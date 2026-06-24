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

// GetSearchSuggestions godoc
// @Summary      Get search auto-complete suggestions
// @Tags         search
// @Produce      json
// @Param        q  query  string  false  "Search query"
// @Success      200  {object}  SearchSuggestionsResponse
// @Router       /search/suggestions [get]
func (h *SearchHandler) GetSearchSuggestions(c *gin.Context) {
	suggestions, err := h.service.GetSuggestions(c.Query("q"))
	if err != nil {
		utils.InternalError(c, "Failed to fetch suggestions")
		return
	}
	utils.OK(c, suggestions)
}

// SaveSearchHistory godoc
// @Summary      Save a search term to history (no-op if not authenticated)
// @Tags         search
// @Accept       json
// @Produce      json
// @Param        body  body      SaveSearchHistoryRequest  false  "Search term"
// @Success      200   {object}  MessageResponse
// @Router       /search/history [post]
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

// GetSearchHistory godoc
// @Summary      Get the current user's search history
// @Tags         search
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  SearchHistoryResponse
// @Failure      401  {object}  ErrorResponse
// @Router       /search/history [get]
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

// DeleteSearchHistory godoc
// @Summary      Clear the current user's search history
// @Tags         search
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  MessageResponse
// @Failure      401  {object}  ErrorResponse
// @Router       /search/history [delete]
func (h *SearchHandler) DeleteSearchHistory(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	h.service.DeleteSearchHistory(userID)
	utils.OK(c, gin.H{"message": "Search history cleared"})
}
