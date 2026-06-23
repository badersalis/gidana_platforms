package services

import (
	"strings"
	"time"

	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/repositories"
)

type SearchService interface {
	GetSuggestions(q string) ([]string, error)
	SaveSearchHistory(userID uint, term string) error
	GetSearchHistory(userID uint) ([]models.SearchHistory, error)
	DeleteSearchHistory(userID uint) error
}

type searchService struct {
	repo repositories.SearchRepository
}

func NewSearchService(repo repositories.SearchRepository) SearchService {
	return &searchService{repo: repo}
}

func (s *searchService) GetSuggestions(q string) ([]string, error) {
	if len(q) < 2 {
		return []string{}, nil
	}
	return s.repo.GetSuggestions(strings.ToLower(q), 5)
}

func (s *searchService) SaveSearchHistory(userID uint, term string) error {
	oneHourAgo := time.Now().Add(-time.Hour)
	_, err := s.repo.FindRecentByUserAndTerm(userID, term, oneHourAgo)
	if err != nil {
		sh := &models.SearchHistory{UserID: userID, SearchTerm: term}
		return s.repo.CreateHistory(sh)
	}
	return nil
}

func (s *searchService) GetSearchHistory(userID uint) ([]models.SearchHistory, error) {
	return s.repo.GetHistoryByUser(userID, 10)
}

func (s *searchService) DeleteSearchHistory(userID uint) error {
	return s.repo.DeleteHistoryByUser(userID)
}
