package handlers

import (
	"github.com/badersalis/gidana_backend/internal/middleware"
	"github.com/badersalis/gidana_backend/internal/services"
	"github.com/badersalis/gidana_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	service services.MessageService
}

func NewMessageHandler(svc services.MessageService) *MessageHandler {
	return &MessageHandler{service: svc}
}

func (h *MessageHandler) StartConversation(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var input struct {
		RecipientID uint   `json:"recipient_id" binding:"required"`
		PropertyID  *uint  `json:"property_id"`
		Message     string `json:"message" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	conv, err := h.service.StartConversation(userID, services.StartConversationInput{
		RecipientID: input.RecipientID,
		PropertyID:  input.PropertyID,
		Message:     input.Message,
	})
	if handleErr(c, err) {
		return
	}
	utils.Created(c, conv)
}

func (h *MessageHandler) GetConversations(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	convs, err := h.service.GetConversations(userID)
	if handleErr(c, err) {
		return
	}
	utils.OK(c, convs)
}

func (h *MessageHandler) GetConversation(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	convID := paramUint(c, "id")

	conv, err := h.service.GetConversation(convID, userID)
	if handleErr(c, err) {
		return
	}
	utils.OK(c, conv)
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	convID := paramUint(c, "id")

	var input struct {
		Content string `json:"content" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	msg, err := h.service.SendMessage(convID, userID, input.Content)
	if handleErr(c, err) {
		return
	}
	utils.Created(c, msg)
}

func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	msgID := paramUint(c, "msgID")

	if handleErr(c, h.service.DeleteMessage(msgID, userID)) {
		return
	}
	utils.OK(c, gin.H{"message": "Message deleted"})
}
