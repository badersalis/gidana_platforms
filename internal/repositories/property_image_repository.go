package repositories

import (
	"github.com/badersalis/gidana_backend/internal/models"
	"gorm.io/gorm"
)

type PropertyImageRepository interface {
	Create(img *models.PropertyImage) error
	GetByID(id uint) (*models.PropertyImage, error)
	CountByPropertyID(propID uint) (int64, error)
	Delete(img *models.PropertyImage) error
	UnsetMainForProperty(propID uint) error
	SetMain(img *models.PropertyImage) error
}

type propertyImageRepository struct{ db *gorm.DB }

func NewPropertyImageRepository(db *gorm.DB) PropertyImageRepository {
	return &propertyImageRepository{db: db}
}

func (r *propertyImageRepository) Create(img *models.PropertyImage) error {
	return r.db.Create(img).Error
}

func (r *propertyImageRepository) GetByID(id uint) (*models.PropertyImage, error) {
	var img models.PropertyImage
	err := r.db.First(&img, id).Error
	if err != nil {
		return nil, err
	}
	return &img, nil
}

func (r *propertyImageRepository) CountByPropertyID(propID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.PropertyImage{}).Where("property_id = ?", propID).Count(&count).Error
	return count, err
}

func (r *propertyImageRepository) Delete(img *models.PropertyImage) error {
	return r.db.Delete(img).Error
}

func (r *propertyImageRepository) UnsetMainForProperty(propID uint) error {
	return r.db.Model(&models.PropertyImage{}).Where("property_id = ?", propID).Update("is_main", false).Error
}

func (r *propertyImageRepository) SetMain(img *models.PropertyImage) error {
	return r.db.Model(img).Update("is_main", true).Error
}
