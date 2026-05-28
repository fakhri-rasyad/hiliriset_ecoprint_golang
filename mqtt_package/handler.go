package mqttpackage

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"

	"hiliriset_ecoprint_golang/models"
	"hiliriset_ecoprint_golang/repositories"
	websocketutils "hiliriset_ecoprint_golang/websocket_utils"
)

const (
	espTimeoutSeconds = 60
	inactivityCheck   = 20 * time.Second
)

var requiredFields = []string{"air_temp", "water_temp", "humidity", "water_sufficient"}

type SessionRecordCreator interface {
	CreateRecord(sessionPubID uuid.UUID, data models.SessionRecordInput) error
	InvalidateSession(sessionPubID uuid.UUID)
}

type Publisher interface {
	Publish(topic string, payload string) error
}

type MQTTHandler struct {
    boSeRepo   repositories.BoSeRepository
    srService  SessionRecordCreator
    mqttPub    Publisher
    wsHub      *websocketutils.Hub
    espRepo    repositories.EspRepository
    krRepo     repositories.KomporRepository
    fabricRepo repositories.FabricTypeRepository
    batcher    *TelemetryBatcher

    mu            sync.Map
    timers        sync.Map
    sessionEspMap sync.Map
    espLastSeen   sync.Map
}

func NewMQTTHandler(
    boSeRepo repositories.BoSeRepository,
    srService SessionRecordCreator,
    wsHub *websocketutils.Hub,
    espRepo repositories.EspRepository,
    krRepo repositories.KomporRepository,
    fabricRepo repositories.FabricTypeRepository,
) *MQTTHandler {
    h := &MQTTHandler{
        boSeRepo:   boSeRepo,
        srService:  srService,
        wsHub:      wsHub,
        espRepo:    espRepo,
        krRepo:     krRepo,
        fabricRepo: fabricRepo,
    }
    h.batcher = NewTelemetryBatcher(srService)
    go h.inactivityChecker()
    return h
}
func (h *MQTTHandler) SetPublisher(pub Publisher) {
	h.mqttPub = pub
}

func (h *MQTTHandler) getSessionLock(id uuid.UUID) *sync.Mutex {
	mu, _ := h.mu.LoadOrStore(id.String(), &sync.Mutex{})
	return mu.(*sync.Mutex)
}

// -------------------------------------------------
// TELEMETRY HANDLER
// -------------------------------------------------

func (h *MQTTHandler) handleTelemetry(client paho.Client, msg paho.Message) {
    if msg.Retained() {
        log.Printf("MQTT skipping retained message on topic: %s", msg.Topic())
        return
    }
    parts := strings.Split(msg.Topic(), "/")
    if len(parts) != 3 {
        log.Printf("MQTT unexpected topic format: %s", msg.Topic())
        return
    }

    espPublicID, err := uuid.Parse(parts[1])
    if err != nil {
        log.Printf("MQTT invalid esp public_id in topic: %s", parts[1])
        return
    }

    h.espLastSeen.Store(espPublicID.String(), time.Now())

    var raw map[string]interface{}
    if err := json.Unmarshal(msg.Payload(), &raw); err != nil {
        log.Printf("MQTT malformed JSON from esp %s", espPublicID)
        h.publishToEsp(espPublicID, models.MQTTErrorPayload{Event: "malformed_json"})
        return
    }

    var missing []string
    for _, field := range requiredFields {
        if _, ok := raw[field]; !ok {
            missing = append(missing, field)
        }
    }
    if len(missing) > 0 {
        log.Printf("MQTT missing fields from esp %s: %v", espPublicID, missing)
        h.publishToEsp(espPublicID, models.MQTTErrorPayload{
            Event:  "missing_fields",
            Fields: missing,
        })
        return
    }

    var telemetry models.MQTTTelemetryPayload
    if err := json.Unmarshal(msg.Payload(), &telemetry); err != nil {
        log.Printf("MQTT bad telemetry values from esp %s: %v", espPublicID, err)
        h.publishToEsp(espPublicID, models.MQTTErrorPayload{Event: "bad_payload"})
        return
    }

    session, err := h.boSeRepo.GetActiveSessionByEspPublicID(espPublicID)
    if err != nil {
        log.Printf("MQTT no active session for esp %s", espPublicID)
        h.publishToEsp(espPublicID, models.MQTTErrorPayload{Event: "no_active_session"})
        return
    }

    lock := h.getSessionLock(session.PublicID)
    lock.Lock()
    defer lock.Unlock()

    if telemetry.Event == "steaming" && session.FinishedAt == nil {
        if err := h.handleSteamingEvent(session); err != nil {
            log.Printf("MQTT failed to handle steaming event for session %s: %v", session.PublicID, err)
        }
    }

    h.batcher.Add(session.PublicID, models.SessionRecordInput{
        SessionID:       session.PublicID,
        AirTemp:         telemetry.AirTemp,
        WaterTemp:       telemetry.WaterTemp,
        Humidity:        telemetry.Humidity,
        WaterSufficient: telemetry.WaterSufficient,
    })

    h.wsHub.Broadcast(session.PublicID, telemetry)
}

func (h *MQTTHandler) handleSteamingEvent(session *models.BoilingSessionBase) error {
    fabric, err := h.fabricRepo.GetById(session.FabricType)
    if err != nil {
        return fmt.Errorf("fabric type %d not found: %w", session.FabricType, err)
    }

    finishedAt := time.Now().Add(time.Duration(fabric.BoilingMinutes) * time.Minute)

    if err := h.boSeRepo.SetFinishedAt(session.PublicID, finishedAt); err != nil {
        return fmt.Errorf("failed to set finished_at: %w", err)
    }

    session.FinishedAt = &finishedAt
    h.StartSessionTimer(session)

    log.Printf("Session %s steaming started, finishes at %v (%d min)",
        session.PublicID, finishedAt, fabric.BoilingMinutes)

    h.wsHub.Broadcast(session.PublicID, map[string]interface{}{
        "event":       "steaming_started",
        "finished_at": finishedAt,
    })

    return nil
}

// -------------------------------------------------
// SESSION TIMER
// -------------------------------------------------

func (h *MQTTHandler) StartSessionTimer(session *models.BoilingSessionBase) {
	if session.FinishedAt == nil {
		log.Printf("Session %s has no finished_at, skipping timer", session.PublicID)
		return
	}

	key := session.PublicID.String()
	_, alreadyStarted := h.timers.LoadOrStore(key, true)
	if alreadyStarted {
		return
	}

	espPublicID := h.resolveEspFromSession(session)
	h.sessionEspMap.Store(key, espPublicID)

	remaining := time.Until(*session.FinishedAt)
	if remaining <= 0 {
		log.Printf("Session %s already past finish time, finishing immediately", session.PublicID)
		go h.finishSession(session.PublicID, espPublicID, false)
		return
	}

	log.Printf("Session %s timer started, finishes in %v", session.PublicID, remaining)
	time.AfterFunc(remaining, func() {
		espID, ok := h.sessionEspMap.Load(key)
		if !ok {
			return
		}
		h.finishSession(session.PublicID, espID.(uuid.UUID), false)
	})
}

func (h *MQTTHandler) resolveEspFromSession(session *models.BoilingSessionBase) uuid.UUID {
	if session.EspPublicID != nil {
		return *session.EspPublicID
	}
	return uuid.Nil
}

// -------------------------------------------------
// FINISH SESSION
// -------------------------------------------------

func (h *MQTTHandler) finishSession(sessionPublicID uuid.UUID, espPublicID uuid.UUID, cancelled bool) {
	session, err := h.boSeRepo.GetSessionByPublicID(sessionPublicID)
	if err != nil || session.BoilingStatus == "finished" || session.BoilingStatus == "cancelled" {
		return
	}

	if cancelled {
		if err := h.boSeRepo.CancelSession(sessionPublicID); err != nil {
			log.Printf("MQTT failed to cancel session %s: %v", sessionPublicID, err)
			return
		}
		log.Printf("Session %s cancelled due to ESP inactivity", sessionPublicID)
	} else {
		if err := h.boSeRepo.FinishSession(sessionPublicID); err != nil {
			log.Printf("MQTT failed to finish session %s: %v", sessionPublicID, err)
			return
		}
		log.Printf("Session %s finished normally", sessionPublicID)
	}

	if espPublicID != uuid.Nil {
		if err := h.espRepo.SetActive(espPublicID, false); err != nil {
			log.Printf("failed to release esp %s: %v", espPublicID, err)
		}
	}

	if session.KomporID != nil {
		if err := h.krRepo.SetActiveByInternalID(*session.KomporID, false); err != nil {
			log.Printf("failed to release kompor: %v", err)
		}
	}

	h.batcher.flushAll()
	h.srService.InvalidateSession(sessionPublicID)

	if h.mqttPub != nil && espPublicID != uuid.Nil {
		payload, _ := json.Marshal(models.MQTTCommandPayload{Command: "Stop"})
		topic := "esp/" + espPublicID.String() + "/command"
		if err := h.mqttPub.Publish(topic, string(payload)); err != nil {
			log.Printf("MQTT failed to publish stop to esp %s: %v", espPublicID, err)
		}
	}

	h.timers.Delete(sessionPublicID.String())
	h.mu.Delete(sessionPublicID.String())
	h.sessionEspMap.Delete(sessionPublicID.String())
	h.espLastSeen.Delete(espPublicID.String())

	h.wsHub.BroadcastFinished(sessionPublicID)
}

// -------------------------------------------------
// INACTIVITY CHECKER
// -------------------------------------------------

func (h *MQTTHandler) inactivityChecker() {
    ticker := time.NewTicker(inactivityCheck)
    defer ticker.Stop()

    for range ticker.C {
        now := time.Now()

        h.espLastSeen.Range(func(key, value interface{}) bool {
            espID := key.(string)
            lastSeen := value.(time.Time)

            if now.Sub(lastSeen).Seconds() <= espTimeoutSeconds {
                return true
            }

            espPublicID, err := uuid.Parse(espID)

            if err != nil {
                h.espLastSeen.Delete(espID)
                return true
            }

            session, err := h.boSeRepo.GetActiveSessionByEspPublicID(espPublicID)
            if err != nil {
                h.espLastSeen.Delete(espID)
                return true
            }

            log.Printf("ESP %s inactive for >%ds, cancelling session %s", espID, espTimeoutSeconds, session.PublicID)
            go h.finishSession(session.PublicID, espPublicID, true)
            h.espLastSeen.Delete(espID)

            return true
        })
    }
}

// -------------------------------------------------
// HELPERS
// -------------------------------------------------

func (h *MQTTHandler) publishToEsp(espPublicID uuid.UUID, payload interface{}) {
	if h.mqttPub == nil {
		return
	}
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal ESP error payload: %v", err)
		return
	}
	topic := "esp/" + espPublicID.String() + "/command"
	if err := h.mqttPub.Publish(topic, string(data)); err != nil {
		log.Printf("Failed to publish to esp %s: %v", espPublicID, err)
	}
}

func (h *MQTTHandler) Stop() {
	h.batcher.Stop()
}

func (h *MQTTHandler) RegisterEspSession(espPublicID uuid.UUID) {
    h.espLastSeen.Store(espPublicID.String(), time.Now())
    log.Printf("ESP %s registered, inactivity timer started", espPublicID)
}

func (h *MQTTHandler) RecoverStaleSessions() {
    sessions, err := h.boSeRepo.GetActiveSessions()
    if err != nil {
        log.Printf("RecoverStaleSessions: failed to fetch active sessions: %v", err)
        return
    }

    for _, session := range sessions {
        espPublicID := h.resolveEspFromSession(&session)

        // If finished_at is set and already past — finish it
        if session.FinishedAt != nil && time.Now().After(*session.FinishedAt) {
            log.Printf("RecoverStaleSessions: session %s past end time, finishing", session.PublicID)
            go h.finishSession(session.PublicID, espPublicID, false)
            continue
        }

        // If finished_at is set and still in future — restart timer
        if session.FinishedAt != nil {
            log.Printf("RecoverStaleSessions: session %s resuming timer", session.PublicID)
            h.StartSessionTimer(&session)
            continue
        }

        // If finished_at is nil and session is old (>10 min) — cancel it
        if time.Since(session.CreatedAt) > 10*time.Minute {
            log.Printf("RecoverStaleSessions: session %s stuck in %s, cancelling", session.PublicID, session.BoilingStatus)
            go h.finishSession(session.PublicID, espPublicID, true)
        }
    }
}
