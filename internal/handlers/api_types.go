package handlers

import "github.com/badersalis/gidana_backend/internal/models"

// ─── Request bodies ───────────────────────────────────────────────────────────

type RegisterRequest struct {
	FirstName   string `json:"first_name" example:"John"`
	LastName    string `json:"last_name" example:"Doe"`
	Email       string `json:"email" example:"john@example.com"`
	PhoneNumber string `json:"phone_number" example:"+237600000000"`
	Password    string `json:"password" example:"secret123"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" example:"john@example.com"`
	Password   string `json:"password" example:"secret123"`
}

type UpdateProfileRequest struct {
	FirstName   string `json:"first_name" example:"John"`
	LastName    string `json:"last_name" example:"Doe"`
	Gender      string `json:"gender" example:"male"`
	DateOfBirth string `json:"date_of_birth" example:"1995-04-12"`
	Locale      string `json:"locale" example:"en"`
	Timezone    string `json:"timezone" example:"Africa/Douala"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" example:"oldpass123"`
	NewPassword     string `json:"new_password" example:"newpass123"`
}

type UpdatePushTokenRequest struct {
	ExpoPushToken string `json:"expo_push_token" example:"ExponentPushToken[...]"`
}

type CreateRentalRequest struct {
	PropertyID   uint    `json:"property_id" example:"42"`
	StartDate    string  `json:"start_date" example:"2026-07-01"`
	EndDate      string  `json:"end_date" example:"2027-07-01"`
	MonthlyPrice float64 `json:"monthly_price" example:"150000"`
}

type UpdateRentalStatusRequest struct {
	Status string `json:"status" example:"occupied" enums:"pending,occupied,available,completed"`
}

type StartConversationRequest struct {
	PropertyID uint   `json:"property_id" example:"42"`
	Message    string `json:"message" example:"Is this property still available?"`
}

type UpgradePlanRequest struct {
	Plan     string `json:"plan" example:"pro" enums:"essential,pro"`
	Currency string `json:"currency" example:"XOF" enums:"XOF,USD"`
}

type UpgradeLandlordPlanRequest struct {
	Plan     string `json:"plan" example:"standard" enums:"standard,agency"`
	Currency string `json:"currency" example:"XOF" enums:"XOF,USD"`
}

type SendMessageRequest struct {
	Content string `json:"content" example:"Hello, I am interested in renting this property."`
}

type CreateReviewRequest struct {
	Rating  int    `json:"rating" example:"4" minimum:"1" maximum:"5"`
	Comment string `json:"comment" example:"Great location and very clean."`
}

type SaveSearchHistoryRequest struct {
	SearchTerm string `json:"search_term" example:"apartment Douala"`
}

// ─── Response envelopes ───────────────────────────────────────────────────────

type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error" example:"invalid request"`
}

type MessageResponse struct {
	Success bool              `json:"success" example:"true"`
	Data    map[string]string `json:"data"`
}

// Auth

type AuthData struct {
	User  models.User `json:"user"`
	Token string      `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type AuthResponse struct {
	Success bool     `json:"success" example:"true"`
	Data    AuthData `json:"data"`
}

type UserResponse struct {
	Success bool        `json:"success" example:"true"`
	Data    models.User `json:"data"`
}

// Properties

type PropertyResponse struct {
	Success bool            `json:"success" example:"true"`
	Data    models.Property `json:"data"`
}

type PropertyListResponse struct {
	Success  bool              `json:"success" example:"true"`
	Data     []models.Property `json:"data"`
	Total    int64             `json:"total" example:"100"`
	Page     int               `json:"page" example:"1"`
	PageSize int               `json:"page_size" example:"10"`
	Pages    int               `json:"pages" example:"10"`
}

type PropertyImageResponse struct {
	Success bool                `json:"success" example:"true"`
	Data    models.PropertyImage `json:"data"`
}

// Rentals

type RentalResponse struct {
	Success bool          `json:"success" example:"true"`
	Data    models.Rental `json:"data"`
}

type RentalListResponse struct {
	Success bool            `json:"success" example:"true"`
	Data    []models.Rental `json:"data"`
}

// Reviews

type ReviewResponse struct {
	Success bool          `json:"success" example:"true"`
	Data    models.Review `json:"data"`
}

type ReviewListResponse struct {
	Success bool            `json:"success" example:"true"`
	Data    []models.Review `json:"data"`
}

// Conversations & Messages

type ConversationResponse struct {
	Success bool                `json:"success" example:"true"`
	Data    models.Conversation `json:"data"`
}

type ConversationListResponse struct {
	Success bool                  `json:"success" example:"true"`
	Data    []models.Conversation `json:"data"`
}

type MessageResp struct {
	Success bool          `json:"success" example:"true"`
	Data    models.Message `json:"data"`
}

// Favorites

type FavoriteListResponse struct {
	Success  bool              `json:"success" example:"true"`
	Data     []models.Property `json:"data"`
	Total    int64             `json:"total" example:"5"`
	Page     int               `json:"page" example:"1"`
	PageSize int               `json:"page_size" example:"10"`
	Pages    int               `json:"pages" example:"1"`
}

type FavoriteToggleData struct {
	Favorited bool   `json:"favorited" example:"true"`
	Message   string `json:"message" example:"Added to favorites"`
}

type FavoriteToggleResponse struct {
	Success bool               `json:"success" example:"true"`
	Data    FavoriteToggleData `json:"data"`
}

// Alerts

type AlertResponse struct {
	Success bool         `json:"success" example:"true"`
	Data    models.Alert `json:"data"`
}

type AlertListResponse struct {
	Success bool           `json:"success" example:"true"`
	Data    []models.Alert `json:"data"`
}

// Search

type SearchSuggestionsResponse struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data"`
}

type SearchHistoryResponse struct {
	Success bool                   `json:"success" example:"true"`
	Data    []models.SearchHistory `json:"data"`
}

// Subscriptions

type UpgradePlanResponse struct {
	Success bool                   `json:"success" example:"true"`
	Data    UpgradePlanResponseData `json:"data"`
}

type UpgradePlanResponseData struct {
	Transaction models.Transaction `json:"transaction"`
	User        models.User        `json:"user"`
}

type AvailabilityResponse struct {
	Success bool `json:"success" example:"true"`
	Data    struct {
		IsAvailable bool `json:"is_available" example:"false"`
	} `json:"data"`
}
