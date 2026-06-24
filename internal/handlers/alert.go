package handlers

import (
	"github.com/badersalis/gidana_backend/internal/middleware"
	"github.com/badersalis/gidana_backend/internal/services"
	"github.com/badersalis/gidana_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type AlertHandler struct {
	service services.AlertService
}

func NewAlertHandler(svc services.AlertService) *AlertHandler {
	return &AlertHandler{service: svc}
}

type alertInput struct {
	Country         string  `json:"country"`
	City            string  `json:"city"`
	Neighborhood    string  `json:"neighborhood"`
	PropertyType    string  `json:"property_type"`
	TransactionType string  `json:"transaction_type"`
	MinRooms        int     `json:"min_rooms"`
	MaxPrice        float64 `json:"max_price"`
	Currency        string  `json:"currency"`
}

func (h *AlertHandler) GetAlerts(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	alerts, err := h.service.GetAlerts(userID)
	if handleErr(c, err) {
		return
	}
	utils.OK(c, alerts)
}

func (h *AlertHandler) CreateAlert(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var body alertInput
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	alert, err := h.service.CreateAlert(userID, services.AlertInput{
		Country:         body.Country,
		City:            body.City,
		Neighborhood:    body.Neighborhood,
		PropertyType:    body.PropertyType,
		TransactionType: body.TransactionType,
		MinRooms:        body.MinRooms,
		MaxPrice:        body.MaxPrice,
		Currency:        body.Currency,
	})
	if handleErr(c, err) {
		return
	}
	utils.Created(c, alert)
}

func (h *AlertHandler) UpdateAlert(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	alertID := paramUint(c, "id")

	var body struct {
		alertInput
		IsActive *bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	alert, err := h.service.UpdateAlert(alertID, userID, services.UpdateAlertInput{
		AlertInput: services.AlertInput{
			Country:         body.Country,
			City:            body.City,
			Neighborhood:    body.Neighborhood,
			PropertyType:    body.PropertyType,
			TransactionType: body.TransactionType,
			MinRooms:        body.MinRooms,
			MaxPrice:        body.MaxPrice,
			Currency:        body.Currency,
		},
		IsActive: body.IsActive,
	})
	if handleErr(c, err) {
		return
	}
	utils.OK(c, alert)
}

func (h *AlertHandler) DeleteAlert(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	alertID := paramUint(c, "id")

	if handleErr(c, h.service.DeleteAlert(alertID, userID)) {
		return
	}
	utils.OK(c, gin.H{"message": "Alert deleted"})
}
