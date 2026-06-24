package repositories

import (
	"github.com/badersalis/gidana_backend/internal/models"
	"gorm.io/gorm"
)

type RentalRepository interface {
	GetByTenantID(tenantID uint) ([]models.Rental, error)
	GetByID(id uint) (*models.Rental, error)
	Create(rental *models.Rental) error
	UpdateStatus(rental *models.Rental, status string) error
	ReloadWithProperty(rental *models.Rental) error
}

type rentalRepository struct{ db *gorm.DB }

func NewRentalRepository(db *gorm.DB) RentalRepository {
	return &rentalRepository{db: db}
}

func (r *rentalRepository) GetByTenantID(tenantID uint) ([]models.Rental, error) {
	var rentals []models.Rental
	err := r.db.Where("tenant_id = ?", tenantID).
		Preload("Property.Images").
		Find(&rentals).Error
	return rentals, err
}

func (r *rentalRepository) GetByID(id uint) (*models.Rental, error) {
	var rental models.Rental
	err := r.db.First(&rental, id).Error
	if err != nil {
		return nil, err
	}
	return &rental, nil
}

func (r *rentalRepository) Create(rental *models.Rental) error {
	return r.db.Create(rental).Error
}

func (r *rentalRepository) UpdateStatus(rental *models.Rental, status string) error {
	return r.db.Model(rental).Update("status", status).Error
}

func (r *rentalRepository) ReloadWithProperty(rental *models.Rental) error {
	return r.db.Preload("Property.Images").First(rental, rental.ID).Error
}
