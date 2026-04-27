package websocketutils

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/google/uuid"
)

type Hub struct {
    clients map[string][]*websocket.Conn
    mu      sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        clients: make(map[string][]*websocket.Conn),
    }
}

func (h *Hub) Register(sessionID uuid.UUID, conn *websocket.Conn) {
    h.mu.Lock()
    defer h.mu.Unlock()

    key := sessionID.String()
    h.clients[key] = append(h.clients[key], conn)
    log.Printf("WS client registered for session %s, total: %d", key, len(h.clients[key]))
}

func (h *Hub) Unregister(sessionID uuid.UUID, conn *websocket.Conn) {
    h.mu.Lock()
    defer h.mu.Unlock()

    key := sessionID.String()
    conns := h.clients[key]

    for i, c := range conns {
        if c == conn {
            h.clients[key] = append(conns[:i], conns[i+1:]...)
            break
        }
    }

    // Clean up empty session entry
    if len(h.clients[key]) == 0 {
        delete(h.clients, key)
    }

    log.Printf("WS client unregistered from session %s", key)
}

func (h *Hub) Broadcast(sessionID uuid.UUID, payload any) {
    h.mu.RLock()
    defer h.mu.RUnlock()

    key := sessionID.String()
    conns, ok := h.clients[key]
    if !ok || len(conns) == 0 {
        return
    }

    data, err := json.Marshal(payload)
    if err != nil {
        log.Printf("WS failed to marshal broadcast payload: %v", err)
        return
    }

    for _, conn := range conns {
        if err := conn.WriteMessage(1, data); err != nil {
            log.Printf("WS failed to write to client: %v", err)
        }
    }
}

func (h *Hub) BroadcastFinished(sessionID uuid.UUID) {
    type finishedEvent struct {
        Event string `json:"event"`
    }

    data, err := json.Marshal(finishedEvent{Event: "finished"})
    if err != nil {
        log.Printf("WS failed to marshal finished event: %v", err)
        return
    }

    h.mu.Lock()
    defer h.mu.Unlock()

    key := sessionID.String()
    for _, conn := range h.clients[key] {
        if err := conn.WriteMessage(1, data); err != nil {
            log.Printf("WS failed to write finished event: %v", err)
        }
        conn.Close()
    }
    delete(h.clients, key)

    log.Printf("WS session %s finished, all clients disconnected", key)
}

func (h *Hub) Run() {
    log.Println("WS hub running")
    select {}
}
