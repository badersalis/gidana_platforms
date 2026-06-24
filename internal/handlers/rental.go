package handlers

import (
	"github.com/badersalis/gidana_backend/internal/middleware"
	"github.com/badersalis/gidana_backend/internal/services"
	"github.com/badersalis/gidana_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type RentalHandler struct {
	service services.RentalService
}

func NewRentalHandler(svc services.RentalService) *RentalHandler {
	return &RentalHandler{service: svc}
}

func (h *RentalHandler) GetMyRentals(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	rentals, err := h.service.GetMyRentals(userID)
	if handleErr(c, err) {
		return
	}
	utils.OK(c, rentals)
}

func (h *RentalHandler) CreateRental(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var input struct {
		PropertyID   uint    `json:"property_id" binding:"required"`
		StartDate    string  `json:"start_date" binding:"required"`
		EndDate      string  `json:"end_date"`
		MonthlyPrice float64 `json:"monthly_price" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	rental, err := h.service.CreateRental(services.RentalInput{
		PropertyID:   input.PropertyID,
		StartDate:    input.StartDate,
		EndDate:      input.EndDate,
		MonthlyPrice: input.MonthlyPrice,
	}, userID)
	if handleErr(c, err) {
		return
	}
	utils.Created(c, rental)
}

func (h *RentalHandler) UpdateRentalStatus(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	rentalID := paramUint(c, "id")

	var input struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if handleErr(c, h.service.UpdateRentalStatus(rentalID, userID, input.Status)) {
		return
	}
	utils.OK(c, gin.H{"message": "Status updated", "status": input.Status})
}
