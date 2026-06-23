package services

import (
	"fmt"

	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/repositories"
	"github.com/badersalis/gidana_backend/internal/utils"
	appws "github.com/badersalis/gidana_backend/internal/ws"
	"gorm.io/gorm"
)

type StartConversationInput struct {
	RecipientID uint
	PropertyID  *uint
	Message     string
}

type MessageService interface {
	StartConversation(userID uint, input StartConversationInput) (*models.Conversation, error)
	GetConversations(userID uint) ([]models.Conversation, error)
	GetConversation(convID, userID uint) (*models.Conversation, error)
	SendMessage(convID, userID uint, content string) (*models.Message, error)
	DeleteMessage(msgID, userID uint) error
}

type messageService struct {
	convRepo repositories.ConversationRepository
	msgRepo  repositories.MessageRepository
	userRepo repositories.UserRepository
	hub      *appws.Hub
}

func NewMessageService(
	convRepo repositories.ConversationRepository,
	msgRepo repositories.MessageRepository,
	userRepo repositories.UserRepository,
	hub *appws.Hub,
) MessageService {
	return &messageService{convRepo: convRepo, msgRepo: msgRepo, userRepo: userRepo, hub: hub}
}

func (s *messageService) notifyRecipient(recipientID uint, senderName string, msg models.Message) {
	s.hub.Emit(recipientID, appws.Event{Type: "new_message", Data: msg})

	if !s.hub.IsOnline(recipientID) {
		if recipient, err := s.userRepo.GetByIDWithToken(recipientID); err == nil && recipient.ExpoPushToken != "" {
			utils.SendExpoPush(
				recipient.ExpoPushToken,
				senderName,
				msg.Content,
				map[string]any{"conversation_id": msg.ConversationID},
			)
		}
	}
}

func (s *messageService) StartConversation(userID uint, input StartConversationInput) (*models.Conversation, error) {
	if input.RecipientID == userID {
		return nil, ErrBadRequest("Cannot start a conversation with yourself")
	}
	if _, err := s.userRepo.GetByID(input.RecipientID); err != nil {
		return nil, ErrNotFound("Recipient not found")
	}

	conv, err := s.convRepo.FindBetweenUsers(userID, input.RecipientID, input.PropertyID)
	if err == gorm.ErrRecordNotFound || conv == nil {
		conv = &models.Conversation{
			InitiatorID: userID,
			RecipientID: input.RecipientID,
			PropertyID:  input.PropertyID,
		}
		if err := s.convRepo.Create(conv); err != nil {
			return nil, ErrInternal("Failed to create conversation")
		}
	} else if err != nil {
		return nil, ErrInternal("Failed to find conversation")
	}

	msg := &models.Message{
		ConversationID: conv.ID,
		SenderID:       userID,
		Content:        input.Message,
	}
	s.msgRepo.Create(msg)
	s.convRepo.TouchUpdatedAt(conv.ID)
	s.msgRepo.ReloadWithSender(msg)

	full, err := s.convRepo.GetWithMessages(conv.ID)
	if err != nil {
		return nil, ErrInternal("Failed to load conversation")
	}

	senderName := fmt.Sprintf("%s %s", msg.Sender.FirstName, msg.Sender.LastName)
	go s.notifyRecipient(input.RecipientID, senderName, *msg)

	return full, nil
}

func (s *messageService) GetConversations(userID uint) ([]models.Conversation, error) {
	convs, err := s.convRepo.GetForUser(userID)
	if err != nil {
		return nil, err
	}
	for i := range convs {
		if last, err := s.msgRepo.GetLatest(convs[i].ID); err == nil {
			convs[i].LastMessage = last
		}
		if count, err := s.msgRepo.CountUnread(convs[i].ID, userID); err == nil {
			convs[i].UnreadCount = int(count)
		}
	}
	return convs, nil
}

func (s *messageService) GetConversation(convID, userID uint) (*models.Conversation, error) {
	conv, err := s.convRepo.GetByID(convID)
	if err != nil {
		return nil, ErrNotFound("Conversation not found")
	}
	if conv.InitiatorID != userID && conv.RecipientID != userID {
		return nil, ErrForbidden("Not authorized")
	}

	msgs, err := s.msgRepo.GetByConversationID(convID)
	if err != nil {
		return nil, ErrInternal("Failed to load messages")
	}
	conv.Messages = msgs
	s.msgRepo.MarkAsRead(convID, userID)
	return conv, nil
}

func (s *messageService) SendMessage(convID, userID uint, content string) (*models.Message, error) {
	conv, err := s.convRepo.GetByID(convID)
	if err != nil {
		return nil, ErrNotFound("Conversation not found")
	}
	if conv.InitiatorID != userID && conv.RecipientID != userID {
		return nil, ErrForbidden("Not authorized")
	}

	msg := &models.Message{
		ConversationID: convID,
		SenderID:       userID,
		Content:        content,
	}
	s.msgRepo.Create(msg)
	s.convRepo.TouchUpdatedAt(convID)
	s.msgRepo.ReloadWithSender(msg)

	recipientID := conv.RecipientID
	if userID == conv.RecipientID {
		recipientID = conv.InitiatorID
	}
	senderName := fmt.Sprintf("%s %s", msg.Sender.FirstName, msg.Sender.LastName)
	go s.notifyRecipient(recipientID, senderName, *msg)

	return msg, nil
}

func (s *messageService) DeleteMessage(msgID, userID uint) error {
	msg, err := s.msgRepo.GetByID(msgID)
	if err != nil {
		return ErrNotFound("Message not found")
	}
	if msg.SenderID != userID {
		return ErrForbidden("Not authorized")
	}
	return s.msgRepo.Delete(msg)
}
