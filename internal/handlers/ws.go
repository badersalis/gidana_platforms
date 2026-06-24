package handlers

import (
	"github.com/badersalis/gidana_backend/internal/utils"
	appws "github.com/badersalis/gidana_backend/internal/ws"
	"github.com/gin-gonic/gin"
)

type WSHandler struct {
	hub *appws.Hub
}

func NewWSHandler(hub *appws.Hub) *WSHandler {
	return &WSHandler{hub: hub}
}

// ServeWS upgrades an HTTP connection to WebSocket.
// Auth uses ?token= because browsers cannot send Authorization headers in WS handshakes.
func (h *WSHandler) ServeWS(c *gin.Context) {
	tokenStr := c.Query("token")
	claims, err := utils.ParseToken(tokenStr)
	if err != nil {
		utils.Unauthorized(c, "Invalid token")
		return
	}

	conn, err := appws.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	h.hub.Connect(claims.UserID, conn)
	defer h.hub.Disconnect(claims.UserID)

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}
