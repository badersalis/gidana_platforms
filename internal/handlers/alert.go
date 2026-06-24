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

// GetAlerts godoc
// @Summary      List the current user's property alerts
// @Tags         alerts
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  AlertListResponse
// @Failure      401  {object}  ErrorResponse
// @Router       /alerts [get]
func (h *AlertHandler) GetAlerts(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	alerts, err := h.service.GetAlerts(userID)
	if handleErr(c, err) {
		return
	}
	utils.OK(c, alerts)
}

// CreateAlert godoc
// @Summary      Create a property alert
// @Tags         alerts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      alertInput  true  "Alert criteria"
// @Success      201   {object}  AlertResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      401   {object}  ErrorResponse
// @Router       /alerts [post]
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

// UpdateAlert godoc
// @Summary      Update a property alert
// @Tags         alerts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  int         true  "Alert ID"
// @Param        body  body  alertInput  true  "Updated alert criteria"
// @Success      200   {object}  AlertResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      401   {object}  ErrorResponse
// @Failure      403   {object}  ErrorResponse
// @Failure      404   {object}  ErrorResponse
// @Router       /alerts/{id} [put]
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

// DeleteAlert godoc
// @Summary      Delete a property alert
// @Tags         alerts
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "Alert ID"
// @Success      200  {object}  MessageResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /alerts/{id} [delete]
func (h *AlertHandler) DeleteAlert(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	alertID := paramUint(c, "id")

	if handleErr(c, h.service.DeleteAlert(alertID, userID)) {
		return
	}
	utils.OK(c, gin.H{"message": "Alert deleted"})
}
