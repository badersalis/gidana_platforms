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

// GetPropertyReviews godoc
// @Summary      List all reviews for a property
// @Tags         reviews
// @Produce      json
// @Param        id  path  int  true  "Property ID"
// @Success      200  {object}  ReviewListResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /properties/{id}/reviews [get]
func (h *ReviewHandler) GetPropertyReviews(c *gin.Context) {
	propID := paramUint(c, "id")
	reviews, err := h.service.GetPropertyReviews(propID)
	if handleErr(c, err) {
		return
	}
	utils.OK(c, reviews)
}

// CreateReview godoc
// @Summary      Submit a review for a property
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  int                 true  "Property ID"
// @Param        body  body  CreateReviewRequest  true  "Review details"
// @Success      201   {object}  ReviewResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      401   {object}  ErrorResponse
// @Router       /properties/{id}/reviews [post]
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

// DeleteReview godoc
// @Summary      Delete a review
// @Tags         reviews
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "Review ID"
// @Success      200  {object}  MessageResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /reviews/{id} [delete]
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	reviewID := paramUint(c, "id")

	if handleErr(c, h.service.DeleteReview(reviewID, userID)) {
		return
	}
	utils.OK(c, gin.H{"message": "Review deleted"})
}
