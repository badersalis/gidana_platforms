package models

import (
	"time"

	"gorm.io/gorm"
)

type Property struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Listing info
	Title           string `gorm:"size:100;not null" json:"title"`
	Description     string `gorm:"type:text" json:"description"`
	PropertyType    string `gorm:"size:50;not null" json:"property_type"`    // Bedsitter, Studio, Apartment, Maisonette, Bungalow, Townhouse, Villa, Commercial
	TransactionType string `gorm:"size:20;not null" json:"transaction_type"` // rent, sale

	// Location — from broad to precise
	Country      string  `gorm:"size:50;not null;default:''" json:"country"`
	City         string  `gorm:"size:100;not null;default:''" json:"city"`
	State        string  `gorm:"size:100" json:"state"`        // state / county / region / province (optional)
	Neighborhood string  `gorm:"size:100" json:"neighborhood"` // sub-area within the city (optional)
	ExactAddress string  `gorm:"size:255" json:"exact_address"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`

	// Size & layout
	Rooms     int     `gorm:"not null" json:"rooms"`
	Bathrooms int     `gorm:"not null" json:"bathrooms"`
	ShowerType string `gorm:"size:20" json:"shower_type"` // en_suite, shared
	Surface   float64 `json:"surface"`                    // m²

	// Amenities
	Furnished      bool `gorm:"default:false" json:"furnished"`
	HasWifi        bool `gorm:"default:false" json:"has_wifi"`
	HasWater       bool `gorm:"default:false" json:"has_water"`
	HasElectricity bool `gorm:"default:false" json:"has_electricity"`
	HasCourtyard   bool `gorm:"default:false" json:"has_courtyard"`

	// Pricing — ISO 4217 currency codes (KES, USD, EUR, GBP, NGN, GHS, TZS, UGX …)
	Price    float64 `gorm:"not null" json:"price"`
	Currency string  `gorm:"size:5;not null" json:"currency"`

	// Contact & status
	WhatsappContact string `gorm:"size:20" json:"-"`
	PhoneContact    string `gorm:"size:20" json:"-"`
	IsAvailable     bool   `gorm:"default:true" json:"is_available"`
	OwnerID         uint   `gorm:"not null" json:"owner_id"`

	// Relations
	Owner     User            `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Images    []PropertyImage `gorm:"foreignKey:PropertyID;constraint:OnDelete:CASCADE" json:"images,omitempty"`
	Rentals   []Rental        `gorm:"foreignKey:PropertyID;constraint:OnDelete:CASCADE" json:"-"`
	Reviews   []Review        `gorm:"foreignKey:PropertyID;constraint:OnDelete:CASCADE" json:"reviews,omitempty"`
	Favorites []Favorite      `gorm:"foreignKey:PropertyID;constraint:OnDelete:CASCADE" json:"-"`

	// Computed (not stored)
	AverageRating float64 `gorm:"-" json:"average_rating"`
	ReviewCount   int     `gorm:"-" json:"review_count"`
	IsFavorited   bool    `gorm:"-" json:"is_favorited"`
}

func (p *Property) ComputeRating() {
	if len(p.Reviews) == 0 {
		p.AverageRating = 0
		p.ReviewCount = 0
		return
	}
	total := 0
	for _, r := range p.Reviews {
		total += r.Rating
	}
	p.ReviewCount = len(p.Reviews)
	p.AverageRating = float64(total) / float64(p.ReviewCount)
}
