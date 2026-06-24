package repositories

import (
	"strings"

	"github.com/badersalis/gidana_backend/internal/models"
	"gorm.io/gorm"
)

type PropertyFilters struct {
	Q               string
	Country         string
	City            string
	PropertyType    string
	TransactionType string
	MinPrice        string
	MaxPrice        string
}

type PropertyRepository interface {
	List(filters PropertyFilters, offset, limit int) ([]models.Property, int64, error)
	GetByID(id uint) (*models.Property, error)
	GetFeatured() ([]models.Property, error)
	GetByOwnerID(ownerID uint) ([]models.Property, error)
	Create(prop *models.Property) error
	Update(prop *models.Property, updates map[string]interface{}) error
	Delete(prop *models.Property) error
	ToggleAvailability(prop *models.Property, newVal bool) error
}

type propertyRepository struct{ db *gorm.DB }

func NewPropertyRepository(db *gorm.DB) PropertyRepository {
	return &propertyRepository{db: db}
}

func (r *propertyRepository) List(filters PropertyFilters, offset, limit int) ([]models.Property, int64, error) {
	q := r.db.Model(&models.Property{}).Where("is_available = ?", true)

	if filters.Q != "" {
		like := "%" + strings.ToLower(filters.Q) + "%"
		q = q.Where(
			"LOWER(city) LIKE ? OR LOWER(neighborhood) LIKE ? OR LOWER(country) LIKE ? OR LOWER(state) LIKE ? OR LOWER(exact_address) LIKE ? OR LOWER(title) LIKE ?",
			like, like, like, like, like, like,
		)
	}
	if filters.Country != "" {
		q = q.Where("country = ?", filters.Country)
	}
	if filters.City != "" {
		q = q.Where("city = ?", filters.City)
	}
	if filters.PropertyType != "" {
		q = q.Where("property_type = ?", filters.PropertyType)
	}
	if filters.TransactionType != "" {
		q = q.Where("transaction_type = ?", filters.TransactionType)
	}
	if filters.MinPrice != "" {
		q = q.Where("price >= ?", filters.MinPrice)
	}
	if filters.MaxPrice != "" {
		q = q.Where("price <= ?", filters.MaxPrice)
	}

	var total int64
	q.Count(&total)

	var props []models.Property
	err := q.Offset(offset).Limit(limit).Preload("Images").Preload("Owner").Find(&props).Error
	return props, total, err
}

func (r *propertyRepository) GetByID(id uint) (*models.Property, error) {
	var prop models.Property
	err := r.db.Preload("Images").Preload("Owner").Preload("Reviews.User").First(&prop, id).Error
	if err != nil {
		return nil, err
	}
	return &prop, nil
}

func (r *propertyRepository) GetFeatured() ([]models.Property, error) {
	var props []models.Property
	err := r.db.Where("is_available = ?", true).
		Preload("Images").
		Preload("Reviews").
		Find(&props).Error
	return props, err
}

func (r *propertyRepository) GetByOwnerID(ownerID uint) ([]models.Property, error) {
	var props []models.Property
	err := r.db.Where("owner_id = ?", ownerID).Preload("Images").Find(&props).Error
	return props, err
}

func (r *propertyRepository) Create(prop *models.Property) error {
	return r.db.Create(prop).Error
}

func (r *propertyRepository) Update(prop *models.Property, updates map[string]interface{}) error {
	return r.db.Model(prop).Updates(updates).Error
}

func (r *propertyRepository) Delete(prop *models.Property) error {
	return r.db.Select("Images", "Rentals", "Reviews", "Favorites").Delete(prop).Error
}

func (r *propertyRepository) ToggleAvailability(prop *models.Property, newVal bool) error {
	return r.db.Model(prop).Update("is_available", newVal).Error
}
