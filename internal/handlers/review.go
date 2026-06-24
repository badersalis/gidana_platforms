package handlers

import (
	"github.com/badersalis/gidana_backend/internal/middleware"
	"github.com/badersalis/gidana_backend/internal/services"
	"github.com/badersalis/gidana_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	service services.ReviewService
}

func NewReviewHandler(svc services.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: svc}
}

func (h *ReviewHandler) GetPropertyReviews(c *gin.Context) {
	propID := paramUint(c, "id")
	reviews, err := h.service.GetPropertyReviews(propID)
	if handleErr(c, err) {
		return
	}
	utils.OK(c, reviews)
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	propID := paramUint(c, "id")

	var input struct {
		Rating  int    `json:"rating" binding:"required,min=1,max=5"`
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	review, err := h.service.CreateReview(userID, propID, services.ReviewInput{
		Rating:  input.Rating,
		Comment: input.Comment,
	})
	if handleErr(c, err) {
		return
	}
	utils.Created(c, review)
}

func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	reviewID := paramUint(c, "id")

	if handleErr(c, h.service.DeleteReview(reviewID, userID)) {
		return
	}
	utils.OK(c, gin.H{"message": "Review deleted"})
}
