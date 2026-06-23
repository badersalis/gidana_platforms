package repositories

import (
	"github.com/badersalis/gidana_backend/internal/models"
	"gorm.io/gorm"
)

type ReviewRepository interface {
	GetByPropertyID(propID uint) ([]models.Review, error)
	GetByUserAndProperty(userID, propID uint) (*models.Review, error)
	GetByID(id uint) (*models.Review, error)
	Create(review *models.Review) error
	Delete(review *models.Review) error
	ReloadWithUser(review *models.Review) error
}

type reviewRepository struct{ db *gorm.DB }

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) GetByPropertyID(propID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Where("property_id = ?", propID).Preload("User").Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) GetByUserAndProperty(userID, propID uint) (*models.Review, error) {
	var review models.Review
	err := r.db.Where("user_id = ? AND property_id = ?", userID, propID).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) GetByID(id uint) (*models.Review, error) {
	var review models.Review
	err := r.db.First(&review, id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) Create(review *models.Review) error {
	return r.db.Create(review).Error
}

func (r *reviewRepository) Delete(review *models.Review) error {
	return r.db.Delete(review).Error
}

func (r *reviewRepository) ReloadWithUser(review *models.Review) error {
	return r.db.Preload("User").First(review, review.ID).Error
}
