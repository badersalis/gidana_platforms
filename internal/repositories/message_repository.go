package repositories

import (
	"time"

	"github.com/badersalis/gidana_backend/internal/models"
	"gorm.io/gorm"
)

type ConversationRepository interface {
	FindBetweenUsers(initiatorID, recipientID uint, propertyID *uint) (*models.Conversation, error)
	Create(conv *models.Conversation) error
	GetByID(id uint) (*models.Conversation, error)
	GetForUser(userID uint) ([]models.Conversation, error)
	TouchUpdatedAt(convID uint) error
	GetWithMessages(id uint) (*models.Conversation, error)
}

type conversationRepository struct{ db *gorm.DB }

func NewConversationRepository(db *gorm.DB) ConversationRepository {
	return &conversationRepository{db: db}
}

func (r *conversationRepository) FindBetweenUsers(initiatorID, recipientID uint, propertyID *uint) (*models.Conversation, error) {
	q := r.db.Where(
		"(initiator_id = ? AND recipient_id = ?) OR (initiator_id = ? AND recipient_id = ?)",
		initiatorID, recipientID, recipientID, initiatorID,
	)
	if propertyID != nil {
		q = q.Where("property_id = ?", *propertyID)
	} else {
		q = q.Where("property_id IS NULL")
	}
	var conv models.Conversation
	err := q.First(&conv).Error
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

func (r *conversationRepository) Create(conv *models.Conversation) error {
	return r.db.Create(conv).Error
}

func (r *conversationRepository) GetByID(id uint) (*models.Conversation, error) {
	var conv models.Conversation
	err := r.db.Preload("Initiator").Preload("Recipient").Preload("Property").First(&conv, id).Error
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

func (r *conversationRepository) GetForUser(userID uint) ([]models.Conversation, error) {
	var convs []models.Conversation
	err := r.db.
		Where("initiator_id = ? OR recipient_id = ?", userID, userID).
		Preload("Initiator").
		Preload("Recipient").
		Preload("Property").
		Order("updated_at DESC").
		Find(&convs).Error
	return convs, err
}

func (r *conversationRepository) TouchUpdatedAt(convID uint) error {
	return r.db.Model(&models.Conversation{}).Where("id = ?", convID).UpdateColumn("updated_at", time.Now()).Error
}

func (r *conversationRepository) GetWithMessages(id uint) (*models.Conversation, error) {
	var conv models.Conversation
	err := r.db.
		Preload("Initiator").
		Preload("Recipient").
		Preload("Property").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Sender").Order("messages.created_at ASC")
		}).
		First(&conv, id).Error
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

type MessageRepository interface {
	Create(msg *models.Message) error
	GetByID(id uint) (*models.Message, error)
	GetByConversationID(convID uint) ([]models.Message, error)
	GetLatest(convID uint) (*models.Message, error)
	CountUnread(convID, excludeSenderID uint) (int64, error)
	MarkAsRead(convID, excludeSenderID uint) error
	Delete(msg *models.Message) error
	ReloadWithSender(msg *models.Message) error
}

type messageRepository struct{ db *gorm.DB }

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(msg *models.Message) error {
	return r.db.Create(msg).Error
}

func (r *messageRepository) GetByID(id uint) (*models.Message, error) {
	var msg models.Message
	err := r.db.First(&msg, id).Error
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (r *messageRepository) GetByConversationID(convID uint) ([]models.Message, error) {
	var msgs []models.Message
	err := r.db.Where("conversation_id = ?", convID).
		Preload("Sender").
		Order("created_at ASC").
		Find(&msgs).Error
	return msgs, err
}

func (r *messageRepository) GetLatest(convID uint) (*models.Message, error) {
	var msg models.Message
	err := r.db.Where("conversation_id = ?", convID).
		Preload("Sender").
		Order("created_at DESC").
		First(&msg).Error
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (r *messageRepository) CountUnread(convID, excludeSenderID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Message{}).
		Where("conversation_id = ? AND sender_id != ? AND is_read = false", convID, excludeSenderID).
		Count(&count).Error
	return count, err
}

func (r *messageRepository) MarkAsRead(convID, excludeSenderID uint) error {
	return r.db.Model(&models.Message{}).
		Where("conversation_id = ? AND sender_id != ? AND is_read = false", convID, excludeSenderID).
		Update("is_read", true).Error
}

func (r *messageRepository) Delete(msg *models.Message) error {
	return r.db.Delete(msg).Error
}

func (r *messageRepository) ReloadWithSender(msg *models.Message) error {
	return r.db.Preload("Sender").First(msg, msg.ID).Error
}
