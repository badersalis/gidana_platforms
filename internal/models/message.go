package models

import (
	"time"

	"gorm.io/gorm"
)

type Conversation struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	PropertyID  *uint          `gorm:"index" json:"property_id,omitempty"`
	InitiatorID uint           `gorm:"not null;index" json:"tenant_id"`
	RecipientID uint           `gorm:"not null;index" json:"owner_id"`

	Property  *Property `gorm:"foreignKey:PropertyID" json:"property,omitempty"`
	Initiator User      `gorm:"foreignKey:InitiatorID" json:"tenant,omitempty"`
	Recipient User      `gorm:"foreignKey:RecipientID" json:"owner,omitempty"`
	Messages  []Message `gorm:"foreignKey:ConversationID;constraint:OnDelete:CASCADE" json:"messages,omitempty"`

	UnreadCount int      `gorm:"-" json:"unread_count"`
	LastMessage *Message `gorm:"-" json:"last_message,omitempty"`
}

type Message struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	ConversationID uint           `gorm:"not null;index" json:"conversation_id"`
	SenderID       uint           `gorm:"not null" json:"sender_id"`
	Content        string         `gorm:"type:text;not null" json:"content"`
	IsRead         bool           `gorm:"default:false" json:"is_read"`

	Sender User `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
}
