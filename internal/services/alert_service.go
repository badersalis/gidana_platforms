package services

import (
	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/repositories"
)

type AlertInput struct {
	Country         string
	City            string
	Neighborhood    string
	PropertyType    string
	TransactionType string
	MinRooms        int
	MaxPrice        float64
	Currency        string
}

type UpdateAlertInput struct {
	AlertInput
	IsActive *bool
}

type AlertService interface {
	GetAlerts(userID uint) ([]models.Alert, error)
	CreateAlert(userID uint, input AlertInput) (*models.Alert, error)
	UpdateAlert(alertID, userID uint, input UpdateAlertInput) (*models.Alert, error)
	DeleteAlert(alertID, userID uint) error
}

type alertService struct {
	repo repositories.AlertRepository
}

func NewAlertService(repo repositories.AlertRepository) AlertService {
	return &alertService{repo: repo}
}

func (s *alertService) GetAlerts(userID uint) ([]models.Alert, error) {
	return s.repo.GetByUserID(userID)
}

func (s *alertService) CreateAlert(userID uint, input AlertInput) (*models.Alert, error) {
	alert := &models.Alert{
		UserID:          userID,
		Country:         input.Country,
		City:            input.City,
		Neighborhood:    input.Neighborhood,
		PropertyType:    input.PropertyType,
		TransactionType: input.TransactionType,
		MinRooms:        input.MinRooms,
		MaxPrice:        input.MaxPrice,
		Currency:        input.Currency,
		IsActive:        true,
	}
	if err := s.repo.Create(alert); err != nil {
		return nil, ErrInternal("Failed to create alert")
	}
	return alert, nil
}

func (s *alertService) UpdateAlert(alertID, userID uint, input UpdateAlertInput) (*models.Alert, error) {
	alert, err := s.repo.GetByUserAndID(userID, alertID)
	if err != nil {
		return nil, ErrNotFound("Alert not found")
	}

	updates := map[string]interface{}{
		"country":          input.Country,
		"city":             input.City,
		"neighborhood":     input.Neighborhood,
		"property_type":    input.PropertyType,
		"transaction_type": input.TransactionType,
		"min_rooms":        input.MinRooms,
		"max_price":        input.MaxPrice,
		"currency":         input.Currency,
	}
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}

	if err := s.repo.Update(alert, updates); err != nil {
		return nil, ErrInternal("Failed to update alert")
	}
	return alert, nil
}

func (s *alertService) DeleteAlert(alertID, userID uint) error {
	alert, err := s.repo.GetByUserAndID(userID, alertID)
	if err != nil {
		return ErrNotFound("Alert not found")
	}
	return s.repo.Delete(alert)
}
