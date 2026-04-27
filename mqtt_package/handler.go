package mqttpackage

import (
	"encoding/json"
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

type SessionRecordCreator interface {
    CreateRecord(sessionPubID uuid.UUID, data models.SessionRecordInput) error
    InvalidateSession(sessionPubID uuid.UUID)
}

type Publisher interface {
    Publish(topic string, payload string) error
}

type MQTTHandler struct {
    boSeRepo  repositories.BoSeRepository
    srService SessionRecordCreator
    mqttPub   Publisher
    wsHub     *websocketutils.Hub
    espRepo   repositories.EspRepository
    krRepo    repositories.KomporRepository
    mu            sync.Map
    timers        sync.Map
    sessionEspMap sync.Map
}

func NewMQTTHandler(
    boSeRepo repositories.BoSeRepository,
    srService SessionRecordCreator,
    wsHub *websocketutils.Hub,
    espRepo repositories.EspRepository,
    krRepo repositories.KomporRepository,
) *MQTTHandler {
    return &MQTTHandler{
        boSeRepo:  boSeRepo,
        srService: srService,
        wsHub:     wsHub,
        espRepo:   espRepo,
        krRepo:    krRepo,
    }
}
func (h *MQTTHandler) SetPublisher(pub Publisher) {
    h.mqttPub = pub
}

func (h *MQTTHandler) getSessionLock(id uuid.UUID) *sync.Mutex {
    mu, _ := h.mu.LoadOrStore(id.String(), &sync.Mutex{})
    return mu.(*sync.Mutex)
}

func (h *MQTTHandler) handleTelemetry(client paho.Client, msg paho.Message) {
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

    var telemetry models.MQTTTelemetryPayload
    if err := json.Unmarshal(msg.Payload(), &telemetry); err != nil {
        log.Printf("MQTT failed to parse telemetry from esp %s: %v", espPublicID, err)
        return
    }

    session, err := h.boSeRepo.GetActiveSessionByEspPublicID(espPublicID)
    if err != nil {
        log.Printf("MQTT no active session for esp %s: %v", espPublicID, err)
        return
    }

    lock := h.getSessionLock(session.PublicID)
    lock.Lock()
    defer lock.Unlock()

    if err := h.srService.CreateRecord(session.PublicID, models.SessionRecordInput{
        AirTemp:   telemetry.AirTemp,
        WaterTemp: telemetry.WaterTemp,
        Humidity:  telemetry.Humidity,
    }); err != nil {
        log.Printf("MQTT failed to save record for session %s: %v", session.PublicID, err)
        return
    }

    h.wsHub.Broadcast(session.PublicID, telemetry)
}

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
        go h.finishSession(session.PublicID, espPublicID)
        return
    }

    log.Printf("Session %s timer started, finishes in %v", session.PublicID, remaining)

    time.AfterFunc(remaining, func() {
        espID, ok := h.sessionEspMap.Load(key)
        if !ok {
            return
        }
        h.finishSession(session.PublicID, espID.(uuid.UUID))
    })
}

func (h *MQTTHandler) resolveEspFromSession(session *models.BoilingSessionBase) uuid.UUID {
    if session.EspPublicID != nil {
        return *session.EspPublicID
    }
    return uuid.Nil
}

func (h *MQTTHandler) finishSession(sessionPublicID uuid.UUID, espPublicID uuid.UUID) {
    session, err := h.boSeRepo.GetSessionByPublicID(sessionPublicID)
    if err != nil || session.BoilingStatus == "finished" {
        return
    }

    if err := h.boSeRepo.FinishSession(sessionPublicID); err != nil {
        log.Printf("MQTT failed to finish session %s: %v", sessionPublicID, err)
        return
    }

    // Release ESP and Kompor
    if err := h.espRepo.SetActive(espPublicID, false); err != nil {
        log.Printf("failed to release esp %s: %v", espPublicID, err)
    }
    if session.KomporID != nil {
        if err := h.krRepo.SetActiveByInternalID(*session.KomporID, false); err != nil {
            log.Printf("failed to release kompor: %v", err)
        }
    }

    h.srService.InvalidateSession(sessionPublicID)

    if h.mqttPub != nil && espPublicID != uuid.Nil {
        topic := "esp/" + espPublicID.String() + "/cmd"
        if err := h.mqttPub.Publish(topic, "stop"); err != nil {
            log.Printf("MQTT failed to publish stop to esp %s: %v", espPublicID, err)
        }
    }

    h.timers.Delete(sessionPublicID.String())
    h.mu.Delete(sessionPublicID.String())
    h.sessionEspMap.Delete(sessionPublicID.String())

    h.wsHub.BroadcastFinished(sessionPublicID)
    log.Printf("Session %s finished, esp and kompor released", sessionPublicID)
}
