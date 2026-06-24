package services

import (
	"strings"

	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/repositories"
	"github.com/badersalis/gidana_backend/internal/utils"
)

type WalletInput struct {
	Provider       string
	Nature         string
	PhoneNumber    string
	Email          string
	CardNumber     string
	CVV            string
	ExpirationDate string
	Password       string
	Currency       string
	Selected       bool
}

type WalletService interface {
	GetWallets(userID uint) ([]models.Wallet, error)
	CreateWallet(userID uint, input WalletInput) (*models.Wallet, error)
	UpdateWallet(userID, walletID uint, input WalletInput) (*models.Wallet, error)
	DeleteWallet(userID, walletID uint) error
	SelectWallet(userID, walletID uint) error
	RefreshWalletBalance(userID, walletID uint) (float64, string, error)
}

var validProviders = map[string]bool{"Nita": true, "MPesa": true, "Visa": true, "Mastercard": true, "PayPal": true}

type walletService struct {
	repo repositories.WalletRepository
}

func NewWalletService(repo repositories.WalletRepository) WalletService {
	return &walletService{repo: repo}
}

func (s *walletService) GetWallets(userID uint) ([]models.Wallet, error) {
	wallets, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	for i := range wallets {
		wallets[i].ApplyMasks()
	}
	return wallets, nil
}

func (s *walletService) CreateWallet(userID uint, input WalletInput) (*models.Wallet, error) {
	if !validProviders[input.Provider] {
		return nil, ErrBadRequest("Invalid provider")
	}

	currency := input.Currency
	if currency == "" {
		currency = "XOF"
	}

	wallet := &models.Wallet{
		UserID:         userID,
		Provider:       input.Provider,
		Nature:         input.Nature,
		PhoneNumber:    input.PhoneNumber,
		Email:          input.Email,
		CardNumber:     input.CardNumber,
		CVV:            input.CVV,
		ExpirationDate: input.ExpirationDate,
		Currency:       currency,
		Balance:        0,
	}

	if input.Password != "" {
		hash, _ := utils.HashPassword(input.Password)
		wallet.Password = hash
	}

	if input.Selected {
		s.repo.DeselectAllForUser(userID)
		wallet.Selected = true
	}

	if err := s.repo.Create(wallet); err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "UNIQUE") {
			return nil, ErrBadRequest("Wallet credentials already in use")
		}
		return nil, ErrInternal("Failed to create wallet")
	}

	wallet.ApplyMasks()
	return wallet, nil
}

func (s *walletService) UpdateWallet(userID, walletID uint, input WalletInput) (*models.Wallet, error) {
	wallet, err := s.repo.GetByUserAndID(userID, walletID)
	if err != nil {
		return nil, ErrNotFound("Wallet not found")
	}

	updates := map[string]interface{}{}
	if input.Nature != "" {
		updates["nature"] = input.Nature
	}
	if input.PhoneNumber != "" {
		updates["phone_number"] = input.PhoneNumber
	}
	if input.Email != "" {
		updates["email"] = input.Email
	}
	if input.Currency != "" {
		updates["currency"] = input.Currency
	}
	if input.Selected {
		s.repo.DeselectAllExcept(userID, wallet.ID)
		updates["selected"] = true
	}

	if err := s.repo.Update(wallet, updates); err != nil {
		return nil, ErrInternal("Failed to update wallet")
	}
	wallet.ApplyMasks()
	return wallet, nil
}

func (s *walletService) DeleteWallet(userID, walletID uint) error {
	wallet, err := s.repo.GetByUserAndID(userID, walletID)
	if err != nil {
		return ErrNotFound("Wallet not found")
	}
	return s.repo.Delete(wallet)
}

func (s *walletService) SelectWallet(userID, walletID uint) error {
	wallet, err := s.repo.GetByUserAndID(userID, walletID)
	if err != nil {
		return ErrNotFound("Wallet not found")
	}
	s.repo.DeselectAllForUser(userID)
	return s.repo.Update(wallet, map[string]interface{}{"selected": true})
}

func (s *walletService) RefreshWalletBalance(userID, walletID uint) (float64, string, error) {
	wallet, err := s.repo.GetByUserAndID(userID, walletID)
	if err != nil {
		return 0, "", ErrNotFound("Wallet not found")
	}

	balances := map[string]float64{"Nita": 1000, "MPesa": 500, "Visa": 750, "Mastercard": 600, "PayPal": 300}
	balance := balances[wallet.Provider]
	s.repo.UpdateBalance(wallet, balance)
	return balance, wallet.Currency, nil
}
