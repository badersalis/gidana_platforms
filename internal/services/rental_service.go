package services

import (
	"time"

	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/repositories"
)

type RentalInput struct {
	PropertyID   uint
	StartDate    string
	EndDate      string
	MonthlyPrice float64
}

type RentalService interface {
	GetMyRentals(tenantID uint) ([]models.Rental, error)
	CreateRental(input RentalInput, tenantID uint) (*models.Rental, error)
	UpdateRentalStatus(rentalID, userID uint, status string) error
}

type rentalService struct {
	repo     repositories.RentalRepository
	propRepo repositories.PropertyRepository
}

func NewRentalService(repo repositories.RentalRepository, propRepo repositories.PropertyRepository) RentalService {
	return &rentalService{repo: repo, propRepo: propRepo}
}

func (s *rentalService) GetMyRentals(tenantID uint) ([]models.Rental, error) {
	return s.repo.GetByTenantID(tenantID)
}

func (s *rentalService) CreateRental(input RentalInput, tenantID uint) (*models.Rental, error) {
	prop, err := s.propRepo.GetByID(input.PropertyID)
	if err != nil {
		return nil, ErrNotFound("Property not found")
	}
	if !prop.IsAvailable {
		return nil, ErrBadRequest("Property is not available")
	}

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return nil, ErrBadRequest("Invalid start date format (YYYY-MM-DD)")
	}

	rental := &models.Rental{
		PropertyID:   input.PropertyID,
		TenantID:     tenantID,
		StartDate:    startDate,
		MonthlyPrice: input.MonthlyPrice,
		Status:       "pending",
	}

	if input.EndDate != "" {
		if endDate, err := time.Parse("2006-01-02", input.EndDate); err == nil {
			rental.EndDate = &endDate
		}
	}

	if err := s.repo.Create(rental); err != nil {
		return nil, ErrInternal("Failed to create rental")
	}
	s.propRepo.ToggleAvailability(prop, false)
	s.repo.ReloadWithProperty(rental)
	return rental, nil
}

func (s *rentalService) UpdateRentalStatus(rentalID, userID uint, status string) error {
	validStatuses := map[string]bool{"pending": true, "occupied": true, "available": true, "completed": true}
	if !validStatuses[status] {
		return ErrBadRequest("Invalid status")
	}

	rental, err := s.repo.GetByID(rentalID)
	if err != nil {
		return ErrNotFound("Rental not found")
	}

	prop, err := s.propRepo.GetByID(rental.PropertyID)
	if err != nil {
		return ErrNotFound("Property not found")
	}
	if prop.OwnerID != userID && rental.TenantID != userID {
		return ErrForbidden("Not authorized")
	}

	if err := s.repo.UpdateStatus(rental, status); err != nil {
		return ErrInternal("Failed to update rental status")
	}

	if status == "completed" || status == "available" {
		s.propRepo.ToggleAvailability(prop, true)
	}
	return nil
}
