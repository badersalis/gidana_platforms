package services

import (
	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/repositories"
)

type ReviewInput struct {
	Rating  int
	Comment string
}

type ReviewService interface {
	GetPropertyReviews(propID uint) ([]models.Review, error)
	CreateReview(userID, propID uint, input ReviewInput) (*models.Review, error)
	DeleteReview(reviewID, userID uint) error
}

type reviewService struct {
	repo     repositories.ReviewRepository
	propRepo repositories.PropertyRepository
}

func NewReviewService(repo repositories.ReviewRepository, propRepo repositories.PropertyRepository) ReviewService {
	return &reviewService{repo: repo, propRepo: propRepo}
}

func (s *reviewService) GetPropertyReviews(propID uint) ([]models.Review, error) {
	return s.repo.GetByPropertyID(propID)
}

func (s *reviewService) CreateReview(userID, propID uint, input ReviewInput) (*models.Review, error) {
	if _, err := s.propRepo.GetByID(propID); err != nil {
		return nil, ErrNotFound("Property not found")
	}

	if _, err := s.repo.GetByUserAndProperty(userID, propID); err == nil {
		return nil, ErrBadRequest("You have already reviewed this property")
	}

	review := &models.Review{
		PropertyID: propID,
		UserID:     userID,
		Rating:     input.Rating,
		Comment:    input.Comment,
	}
	if err := s.repo.Create(review); err != nil {
		return nil, ErrInternal("Failed to create review")
	}
	s.repo.ReloadWithUser(review)
	return review, nil
}

func (s *reviewService) DeleteReview(reviewID, userID uint) error {
	review, err := s.repo.GetByID(reviewID)
	if err != nil {
		return ErrNotFound("Review not found")
	}
	if review.UserID != userID {
		return ErrForbidden("Not authorized")
	}
	return s.repo.Delete(review)
}
