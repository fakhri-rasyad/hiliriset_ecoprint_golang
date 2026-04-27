package websocketutils

import (
	"log"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type WSController struct {
    hub *Hub
}

func NewWSController(hub *Hub) *WSController {
    return &WSController{hub: hub}
}

// HandleSession godoc
// @Summary      WebSocket session monitor
// @Description  Connect via WebSocket to monitor live telemetry for a boiling session
// @Tags         WebSocket
// @Param        session_id path string true "Public ID session"
// @Router       /api/v1/sessions/{session_id}/ws [get]
func (wc *WSController) HandleSession(c fiber.Ctx) error {
    if !websocket.IsWebSocketUpgrade(c) {
        return fiber.ErrUpgradeRequired
    }
    return c.Next()
}

func (wc *WSController) HandleSessionWS(c *websocket.Conn) {
    sessionIDStr := c.Params("session_id")

    sessionID, err := uuid.Parse(sessionIDStr)
    if err != nil {
        log.Printf("WS invalid session_id: %s", sessionIDStr)
        c.Close()
        return
    }

    // Register client
    wc.hub.Register(sessionID, c)
    defer wc.hub.Unregister(sessionID, c)

    log.Printf("WS client connected to session %s", sessionID)

    // Keep connection alive — read loop discards incoming messages
    // but detects client disconnect via read error
    for {
        _, _, err := c.ReadMessage()
        if err != nil {
            log.Printf("WS client disconnected from session %s: %v", sessionID, err)
            break
        }
    }
}
