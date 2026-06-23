package repositories

import (
	"github.com/badersalis/gidana_backend/internal/models"
	"gorm.io/gorm"
)

type FavoriteRepository interface {
	GetByUserID(userID uint, offset, limit int) ([]models.Favorite, int64, error)
	GetByUserAndProperty(userID, propID uint) (*models.Favorite, error)
	ExistsForUserAndProperty(userID, propID uint) (bool, error)
	Create(fav *models.Favorite) error
	Delete(fav *models.Favorite) error
}

type favoriteRepository struct{ db *gorm.DB }

func NewFavoriteRepository(db *gorm.DB) FavoriteRepository {
	return &favoriteRepository{db: db}
}

func (r *favoriteRepository) GetByUserID(userID uint, offset, limit int) ([]models.Favorite, int64, error) {
	var total int64
	r.db.Model(&models.Favorite{}).Where("user_id = ?", userID).Count(&total)

	var favs []models.Favorite
	err := r.db.Where("user_id = ?", userID).
		Preload("Property.Images").
		Preload("Property.Reviews").
		Offset(offset).Limit(limit).
		Find(&favs).Error
	return favs, total, err
}

func (r *favoriteRepository) GetByUserAndProperty(userID, propID uint) (*models.Favorite, error) {
	var fav models.Favorite
	err := r.db.Where("user_id = ? AND property_id = ?", userID, propID).First(&fav).Error
	if err != nil {
		return nil, err
	}
	return &fav, nil
}

func (r *favoriteRepository) ExistsForUserAndProperty(userID, propID uint) (bool, error) {
	var fav models.Favorite
	err := r.db.Where("user_id = ? AND property_id = ?", userID, propID).First(&fav).Error
	if err == nil {
		return true, nil
	}
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	return false, err
}

func (r *favoriteRepository) Create(fav *models.Favorite) error {
	return r.db.Create(fav).Error
}

func (r *favoriteRepository) Delete(fav *models.Favorite) error {
	return r.db.Delete(fav).Error
}
