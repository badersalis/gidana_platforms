package handlers

import (
	"github.com/badersalis/gidana_backend/internal/middleware"
	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/services"
	"github.com/badersalis/gidana_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	service services.SearchService
}

func NewSearchHandler(svc services.SearchService) *SearchHandler {
	return &SearchHandler{service: svc}
}

func (h *SearchHandler) GetSearchSuggestions(c *gin.Context) {
	suggestions, err := h.service.GetSuggestions(c.Query("q"))
	if err != nil {
		utils.InternalError(c, "Failed to fetch suggestions")
		return
	}
	utils.OK(c, suggestions)
}

func (h *SearchHandler) SaveSearchHistory(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.OK(c, gin.H{"saved": false})
		return
	}

	var input struct {
		SearchTerm string `json:"search_term" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	h.service.SaveSearchHistory(userID, input.SearchTerm)
	utils.OK(c, gin.H{"saved": true})
}

func (h *SearchHandler) GetSearchHistory(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	history, err := h.service.GetSearchHistory(userID)
	if err != nil {
		utils.InternalError(c, "Failed to fetch history")
		return
	}
	if history == nil {
		history = []models.SearchHistory{}
	}
	utils.OK(c, history)
}

func (h *SearchHandler) DeleteSearchHistory(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	h.service.DeleteSearchHistory(userID)
	utils.OK(c, gin.H{"message": "Search history cleared"})
}
