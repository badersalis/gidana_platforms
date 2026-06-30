package handlers

import (
	"github.com/badersalis/gidana_backend/internal/middleware"
	"github.com/badersalis/gidana_backend/internal/services"
	"github.com/badersalis/gidana_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	service services.SubscriptionService
}

func NewSubscriptionHandler(svc services.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: svc}
}

// UpgradePlan godoc
// @Summary      Upgrade seeker subscription plan
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      UpgradePlanRequest  true  "Plan upgrade details"
// @Success      200   {object}  UpgradePlanResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      401   {object}  ErrorResponse
// @Router       /subscriptions/upgrade [post]
func (h *SubscriptionHandler) UpgradePlan(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var input struct {
		Plan     string `json:"plan" binding:"required"`
		Currency string `json:"currency"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	tx, user, err := h.service.UpgradeSeekerPlan(userID, input.Plan, input.Currency)
	if handleErr(c, err) {
		return
	}
	utils.OK(c, gin.H{"transaction": tx, "user": user})
}

// UpgradeLandlordPlan godoc
// @Summary      Upgrade landlord subscription plan
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      UpgradeLandlordPlanRequest  true  "Landlord plan upgrade details"
// @Success      200   {object}  UpgradePlanResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      401   {object}  ErrorResponse
// @Router       /subscriptions/landlord-upgrade [post]
func (h *SubscriptionHandler) UpgradeLandlordPlan(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var input struct {
		Plan     string `json:"plan" binding:"required"`
		Currency string `json:"currency"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	tx, user, err := h.service.UpgradeLandlordPlan(userID, input.Plan, input.Currency)
	if handleErr(c, err) {
		return
	}
	utils.OK(c, gin.H{"transaction": tx, "user": user})
}
