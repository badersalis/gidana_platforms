package services

import (
	"time"

	"github.com/badersalis/gidana_backend/internal/config"
	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/repositories"
)

type SubscriptionService interface {
	UpgradeSeekerPlan(userID uint, plan, currency string) (*models.Transaction, *models.User, error)
	UpgradeLandlordPlan(userID uint, plan, currency string) (*models.Transaction, *models.User, error)
	ActivateFromWebhook(cinetpayID, status string) error
}

type subscriptionService struct {
	userRepo repositories.UserRepository
	txRepo   repositories.TransactionRepository
}

func NewSubscriptionService(
	userRepo repositories.UserRepository,
	txRepo repositories.TransactionRepository,
) SubscriptionService {
	return &subscriptionService{userRepo: userRepo, txRepo: txRepo}
}

func (s *subscriptionService) UpgradeSeekerPlan(userID uint, plan, currency string) (*models.Transaction, *models.User, error) {
	if plan != "essential" && plan != "pro" {
		return nil, nil, ErrBadRequest("Invalid seeker plan. Must be 'essential' or 'pro'")
	}
	if currency == "" {
		currency = "XOF"
	}
	price, ok := config.GetSeekerPrice(plan, currency)
	if !ok {
		return nil, nil, ErrBadRequest("Unsupported currency for this plan")
	}

	tx := &models.Transaction{
		UserID:   userID,
		Amount:   price,
		Currency: currency,
		Nature:   "debit",
		Service:  "subscription",
		Plan:     plan,
		Status:   models.TransactionDone,
	}
	if err := s.txRepo.Create(tx); err != nil {
		return nil, nil, ErrInternal("Failed to create transaction record")
	}

	expiry := time.Now().AddDate(0, 0, 30)
	if err := s.userRepo.Update(userID, map[string]interface{}{
		"subscription_plan":       plan,
		"subscription_expires_at": expiry,
	}); err != nil {
		return nil, nil, ErrInternal("Failed to update subscription plan")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, nil, ErrInternal("Failed to fetch updated user")
	}
	return tx, user, nil
}

func (s *subscriptionService) UpgradeLandlordPlan(userID uint, plan, currency string) (*models.Transaction, *models.User, error) {
	if plan != "standard" && plan != "agency" {
		return nil, nil, ErrBadRequest("Invalid landlord plan. Must be 'standard' or 'agency'")
	}
	if currency == "" {
		currency = "XOF"
	}
	price, ok := config.GetLandlordPrice(plan, currency)
	if !ok {
		return nil, nil, ErrBadRequest("Unsupported currency for this plan")
	}

	tx := &models.Transaction{
		UserID:   userID,
		Amount:   price,
		Currency: currency,
		Nature:   "debit",
		Service:  "landlord_subscription",
		Plan:     plan,
		Status:   models.TransactionDone,
	}
	if err := s.txRepo.Create(tx); err != nil {
		return nil, nil, ErrInternal("Failed to create transaction record")
	}

	expiry := time.Now().AddDate(0, 0, 30)
	if err := s.userRepo.Update(userID, map[string]interface{}{
		"landlord_plan":           plan,
		"subscription_expires_at": expiry,
	}); err != nil {
		return nil, nil, ErrInternal("Failed to update landlord plan")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, nil, ErrInternal("Failed to fetch updated user")
	}
	return tx, user, nil
}

func (s *subscriptionService) ActivateFromWebhook(cinetpayID, status string) error {
	tx, err := s.txRepo.GetByCinetpayID(cinetpayID)
	if err != nil {
		return ErrNotFound("Transaction not found")
	}

	if status != "ACCEPTED" {
		_ = s.txRepo.Update(tx, map[string]interface{}{"status": models.TransactionFailed})
		return nil
	}

	_ = s.txRepo.Update(tx, map[string]interface{}{"status": models.TransactionDone})

	expiry := time.Now().AddDate(0, 0, 30)
	updates := map[string]interface{}{"subscription_expires_at": expiry}
	if tx.Service == "subscription" {
		updates["subscription_plan"] = tx.Plan
	} else {
		updates["landlord_plan"] = tx.Plan
	}
	return s.userRepo.Update(tx.UserID, updates)
}
