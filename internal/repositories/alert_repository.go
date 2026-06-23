package repositories

import (
	"github.com/badersalis/gidana_backend/internal/models"
	"gorm.io/gorm"
)

type AlertRepository interface {
	GetByUserID(userID uint) ([]models.Alert, error)
	GetByUserAndID(userID, id uint) (*models.Alert, error)
	Create(alert *models.Alert) error
	Update(alert *models.Alert, updates map[string]interface{}) error
	Delete(alert *models.Alert) error
	FindMatchingAlerts(prop models.Property) ([]models.Alert, error)
}

type alertRepository struct{ db *gorm.DB }

func NewAlertRepository(db *gorm.DB) AlertRepository {
	return &alertRepository{db: db}
}

func (r *alertRepository) GetByUserID(userID uint) ([]models.Alert, error) {
	var alerts []models.Alert
	err := r.db.Where("user_id = ?", userID).Find(&alerts).Error
	return alerts, err
}

func (r *alertRepository) GetByUserAndID(userID, id uint) (*models.Alert, error) {
	var alert models.Alert
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&alert).Error
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

func (r *alertRepository) Create(alert *models.Alert) error {
	return r.db.Create(alert).Error
}

func (r *alertRepository) Update(alert *models.Alert, updates map[string]interface{}) error {
	return r.db.Model(alert).Updates(updates).Error
}

func (r *alertRepository) Delete(alert *models.Alert) error {
	return r.db.Delete(alert).Error
}

func (r *alertRepository) FindMatchingAlerts(prop models.Property) ([]models.Alert, error) {
	var alerts []models.Alert
	err := r.db.Where(
		`is_active = true
		 AND user_id != ?
		 AND (country = '' OR country = ?)
		 AND (city = '' OR city = ?)
		 AND (neighborhood = '' OR neighborhood = ?)
		 AND (property_type = '' OR property_type = ?)
		 AND (transaction_type = '' OR transaction_type = ?)
		 AND (min_rooms = 0 OR min_rooms <= ?)
		 AND (max_price = 0 OR max_price >= ?)`,
		prop.OwnerID,
		prop.Country,
		prop.City,
		prop.Neighborhood,
		prop.PropertyType,
		prop.TransactionType,
		prop.Rooms,
		prop.Price,
	).Find(&alerts).Error
	return alerts, err
}
