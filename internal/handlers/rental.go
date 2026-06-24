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

// GetMyRentals godoc
// @Summary      List the current user's rentals
// @Tags         rentals
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  RentalListResponse
// @Failure      401  {object}  ErrorResponse
// @Router       /rentals [get]
func (h *RentalHandler) GetMyRentals(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	rentals, err := h.service.GetMyRentals(userID)
	if handleErr(c, err) {
		return
	}
	utils.OK(c, rentals)
}

// CreateRental godoc
// @Summary      Create a rental agreement
// @Tags         rentals
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      CreateRentalRequest  true  "Rental details"
// @Success      201   {object}  RentalResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      401   {object}  ErrorResponse
// @Router       /rentals [post]
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

// UpdateRentalStatus godoc
// @Summary      Update a rental's status
// @Tags         rentals
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  int                       true  "Rental ID"
// @Param        body  body  UpdateRentalStatusRequest  true  "New status"
// @Success      200   {object}  MessageResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      401   {object}  ErrorResponse
// @Failure      403   {object}  ErrorResponse
// @Failure      404   {object}  ErrorResponse
// @Router       /rentals/{id}/status [patch]
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
