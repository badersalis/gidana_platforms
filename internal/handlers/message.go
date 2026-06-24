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

// StartConversation godoc
// @Summary      Start a new conversation with another user
// @Tags         conversations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      StartConversationRequest  true  "Conversation details"
// @Success      201   {object}  ConversationResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      401   {object}  ErrorResponse
// @Router       /conversations [post]
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

// GetConversations godoc
// @Summary      List all conversations for the current user
// @Tags         conversations
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  ConversationListResponse
// @Failure      401  {object}  ErrorResponse
// @Router       /conversations [get]
func (h *MessageHandler) GetConversations(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	convs, err := h.service.GetConversations(userID)
	if handleErr(c, err) {
		return
	}
	utils.OK(c, convs)
}

// GetConversation godoc
// @Summary      Get a single conversation with its messages
// @Tags         conversations
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "Conversation ID"
// @Success      200  {object}  ConversationResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /conversations/{id} [get]
func (h *MessageHandler) GetConversation(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	convID := paramUint(c, "id")

	conv, err := h.service.GetConversation(convID, userID)
	if handleErr(c, err) {
		return
	}
	utils.OK(c, conv)
}

// SendMessage godoc
// @Summary      Send a message in a conversation
// @Tags         conversations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  int                true  "Conversation ID"
// @Param        body  body  SendMessageRequest  true  "Message content"
// @Success      201   {object}  MessageResp
// @Failure      400   {object}  ErrorResponse
// @Failure      401   {object}  ErrorResponse
// @Failure      403   {object}  ErrorResponse
// @Router       /conversations/{id}/messages [post]
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

// DeleteMessage godoc
// @Summary      Delete a message
// @Tags         conversations
// @Produce      json
// @Security     BearerAuth
// @Param        id     path  int  true  "Conversation ID"
// @Param        msgID  path  int  true  "Message ID"
// @Success      200   {object}  MessageResponse
// @Failure      401   {object}  ErrorResponse
// @Failure      403   {object}  ErrorResponse
// @Failure      404   {object}  ErrorResponse
// @Router       /conversations/{id}/messages/{msgID} [delete]
func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	msgID := paramUint(c, "msgID")

	if handleErr(c, h.service.DeleteMessage(msgID, userID)) {
		return
	}
	utils.OK(c, gin.H{"message": "Message deleted"})
}
