package services

import (
	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/repositories"
)

type FavoriteService interface {
	GetFavorites(userID uint, page int) ([]models.Property, int64, error)
	ToggleFavorite(userID, propID uint) (bool, error)
}

type favoriteService struct {
	repo     repositories.FavoriteRepository
	propRepo repositories.PropertyRepository
}

func NewFavoriteService(repo repositories.FavoriteRepository, propRepo repositories.PropertyRepository) FavoriteService {
	return &favoriteService{repo: repo, propRepo: propRepo}
}

func (s *favoriteService) GetFavorites(userID uint, page int) ([]models.Property, int64, error) {
	if page < 1 {
		page = 1
	}
	pageSize := 10
	offset := (page - 1) * pageSize

	favs, total, err := s.repo.GetByUserID(userID, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	props := make([]models.Property, 0, len(favs))
	for _, f := range favs {
		f.Property.ComputeRating()
		f.Property.IsFavorited = true
		props = append(props, f.Property)
	}
	return props, total, nil
}

func (s *favoriteService) ToggleFavorite(userID, propID uint) (bool, error) {
	if _, err := s.propRepo.GetByID(propID); err != nil {
		return false, ErrNotFound("Property not found")
	}

	fav, err := s.repo.GetByUserAndProperty(userID, propID)
	if err == nil {
		if delErr := s.repo.Delete(fav); delErr != nil {
			return false, delErr
		}
		return false, nil
	}

	newFav := &models.Favorite{UserID: userID, PropertyID: propID}
	if err := s.repo.Create(newFav); err != nil {
		return false, err
	}
	return true, nil
}
