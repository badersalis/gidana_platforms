package handlers

import (
	"strings"
	"time"

	"github.com/badersalis/gidana_backend/internal/database"
	"github.com/badersalis/gidana_backend/internal/middleware"
	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func GetSearchSuggestions(c *gin.Context) {
	q := c.Query("q")
	if len(q) < 2 {
		utils.OK(c, []string{})
		return
	}

	like := "%" + strings.ToLower(q) + "%"

	var cities []struct{ City string }
	database.DB.Model(&models.Property{}).
		Select("DISTINCT city").
		Where("LOWER(city) LIKE ?", like).
		Limit(5).
		Scan(&cities)

	var neighborhoods []struct{ Neighborhood string }
	database.DB.Model(&models.Property{}).
		Select("DISTINCT neighborhood").
		Where("LOWER(neighborhood) LIKE ? AND neighborhood != ''", like).
		Limit(5).
		Scan(&neighborhoods)

	seen := map[string]bool{}
	suggestions := make([]string, 0, len(cities)+len(neighborhoods))
	for _, row := range cities {
		if !seen[row.City] {
			suggestions = append(suggestions, row.City)
			seen[row.City] = true
		}
	}
	for _, row := range neighborhoods {
		if !seen[row.Neighborhood] {
			suggestions = append(suggestions, row.Neighborhood)
			seen[row.Neighborhood] = true
		}
	}

	utils.OK(c, suggestions)
}

func SaveSearchHistory(c *gin.Context) {
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

	oneHourAgo := time.Now().Add(-time.Hour)
	var existing models.SearchHistory
	result := database.DB.Where("user_id = ? AND search_term = ? AND created_at > ?",
		userID, input.SearchTerm, oneHourAgo).First(&existing)

	if result.Error != nil {
		sh := models.SearchHistory{UserID: userID, SearchTerm: input.SearchTerm}
		database.DB.Create(&sh)
	}

	utils.OK(c, gin.H{"saved": true})
}

func GetSearchHistory(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	var history []models.SearchHistory
	database.DB.Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(10).
		Find(&history)
	utils.OK(c, history)
}

func DeleteSearchHistory(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	database.DB.Where("user_id = ?", userID).Delete(&models.SearchHistory{})
	utils.OK(c, gin.H{"message": "Search history cleared"})
}
