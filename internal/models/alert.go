package models

import "time"

type Alert struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UserID    uint      `gorm:"not null" json:"user_id"`

	// Location filters (all optional — empty means "any")
	Country      string `gorm:"size:50" json:"country"`
	City         string `gorm:"size:100" json:"city"`
	Neighborhood string `gorm:"size:100" json:"neighborhood"`

	// Property filters
	PropertyType    string  `gorm:"size:50" json:"property_type"`
	TransactionType string  `gorm:"size:20" json:"transaction_type"` // rent, sale
	MinRooms        int     `json:"min_rooms"`
	MaxPrice        float64 `json:"max_price"`
	Currency        string  `gorm:"size:5" json:"currency"` // ISO 4217

	IsActive bool `gorm:"default:true" json:"is_active"`
}
