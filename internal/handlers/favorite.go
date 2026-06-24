package handlers

import (
	"strconv"

	"github.com/badersalis/gidana_backend/internal/middleware"
	"github.com/badersalis/gidana_backend/internal/services"
	"github.com/badersalis/gidana_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	service services.FavoriteService
}

func NewFavoriteHandler(svc services.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{service: svc}
}

func (h *FavoriteHandler) GetFavorites(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	props, total, err := h.service.GetFavorites(userID, page)
	if handleErr(c, err) {
		return
	}
	pageSize := 10
	if page < 1 {
		page = 1
	}
	utils.Paginated(c, props, total, page, pageSize)
}

func (h *FavoriteHandler) ToggleFavorite(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	propID := paramUint(c, "id")
	if propID == 0 {
		utils.BadRequest(c, "Invalid property ID")
		return
	}

	favorited, err := h.service.ToggleFavorite(userID, propID)
	if handleErr(c, err) {
		return
	}
	if favorited {
		utils.OK(c, gin.H{"favorited": true, "message": "Added to favorites"})
	} else {
		utils.OK(c, gin.H{"favorited": false, "message": "Removed from favorites"})
	}
}
