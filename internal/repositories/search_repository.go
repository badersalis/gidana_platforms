package repositories

import (
	"time"

	"github.com/badersalis/gidana_backend/internal/models"
	"gorm.io/gorm"
)

type SearchRepository interface {
	GetSuggestions(q string, limit int) ([]string, error)
	FindRecentByUserAndTerm(userID uint, term string, since time.Time) (*models.SearchHistory, error)
	CreateHistory(sh *models.SearchHistory) error
	GetHistoryByUser(userID uint, limit int) ([]models.SearchHistory, error)
	DeleteHistoryByUser(userID uint) error
}

type searchRepository struct{ db *gorm.DB }

func NewSearchRepository(db *gorm.DB) SearchRepository {
	return &searchRepository{db: db}
}

func (r *searchRepository) GetSuggestions(q string, limit int) ([]string, error) {
	like := "%" + q + "%"

	var cities []struct{ City string }
	r.db.Model(&models.Property{}).
		Select("DISTINCT city").
		Where("LOWER(city) LIKE ?", like).
		Limit(limit).
		Scan(&cities)

	var neighborhoods []struct{ Neighborhood string }
	r.db.Model(&models.Property{}).
		Select("DISTINCT neighborhood").
		Where("LOWER(neighborhood) LIKE ? AND neighborhood != ''", like).
		Limit(limit).
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
	return suggestions, nil
}

func (r *searchRepository) FindRecentByUserAndTerm(userID uint, term string, since time.Time) (*models.SearchHistory, error) {
	var sh models.SearchHistory
	err := r.db.Where("user_id = ? AND search_term = ? AND created_at > ?", userID, term, since).First(&sh).Error
	if err != nil {
		return nil, err
	}
	return &sh, nil
}

func (r *searchRepository) CreateHistory(sh *models.SearchHistory) error {
	return r.db.Create(sh).Error
}

func (r *searchRepository) GetHistoryByUser(userID uint, limit int) ([]models.SearchHistory, error) {
	var history []models.SearchHistory
	err := r.db.Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(limit).
		Find(&history).Error
	return history, err
}

func (r *searchRepository) DeleteHistoryByUser(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.SearchHistory{}).Error
}
