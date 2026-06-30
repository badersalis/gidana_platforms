package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	FirstName      string         `gorm:"size:100;not null" json:"first_name"`
	LastName       string         `gorm:"size:100;not null" json:"last_name"`
	PhoneNumber    string         `gorm:"size:20;uniqueIndex" json:"phone_number,omitempty"`
	Email          string         `gorm:"size:255;uniqueIndex" json:"email,omitempty"`
	PasswordHash   string         `gorm:"not null" json:"-"`
	ProfilePicture string         `json:"profile_picture,omitempty"`
	Gender         string         `gorm:"size:20" json:"gender,omitempty"`
	DateOfBirth    *time.Time     `json:"date_of_birth,omitempty"`
	MemberSince    time.Time      `json:"member_since"`
	Active         bool           `gorm:"default:true" json:"active"`
	Locale         string         `gorm:"size:10;default:'en'" json:"locale"`
	Timezone       string         `gorm:"size:50;default:'UTC'" json:"timezone"`
	ExpoPushToken  string         `gorm:"size:200" json:"-"`

	SubscriptionPlan      string     `gorm:"size:20;default:'basic'" json:"subscription_plan"`
	LandlordPlan          string     `gorm:"size:20;default:'free'" json:"landlord_plan"`
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at,omitempty"`

	Properties    []Property      `gorm:"foreignKey:OwnerID" json:"-"`
	Favorites     []Favorite     `gorm:"foreignKey:UserID" json:"-"`
	Alerts        []Alert        `gorm:"foreignKey:UserID" json:"-"`
	Rentals       []Rental       `gorm:"foreignKey:TenantID" json:"-"`
	Reviews       []Review       `gorm:"foreignKey:UserID" json:"-"`
	SearchHistory  []SearchHistory `gorm:"foreignKey:UserID" json:"-"`
}

func (u *User) CanReviewProperty(propertyID uint) bool {
	return true
}
