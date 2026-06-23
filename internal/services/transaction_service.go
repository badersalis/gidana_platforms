package services

import (
	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/repositories"
	"gorm.io/gorm"
)

var servicePlans = map[string]map[string]float64{
	"starlink": {"Basic": 50, "Standard": 150, "Premium": 500},
	"canal+":   {"Standard": 30},
}

type PayServiceInput struct {
	Service         string
	ServiceProvider string
	Plan            string
	WalletID        uint
	Amount          float64
}

type PayServiceResult struct {
	Amount        float64
	Currency      string
	NewBalance    float64
	TransactionID uint
}

type TransferInput struct {
	WalletID  uint
	Recipient string
	Amount    float64
	Provider  string
}

type TransferResult struct {
	Amount        float64
	Currency      string
	NewBalance    float64
	TransactionID uint
}

type TransactionService interface {
	GetTransactions(userID uint, page int) ([]models.Transaction, int64, error)
	PayService(userID uint, input PayServiceInput) (PayServiceResult, error)
	TransferMoney(userID uint, input TransferInput) (TransferResult, error)
}

type transactionService struct {
	walletRepo repositories.WalletRepository
	txRepo     repositories.TransactionRepository
	db         *gorm.DB
}

func NewTransactionService(walletRepo repositories.WalletRepository, txRepo repositories.TransactionRepository, db *gorm.DB) TransactionService {
	return &transactionService{walletRepo: walletRepo, txRepo: txRepo, db: db}
}

func (s *transactionService) GetTransactions(userID uint, page int) ([]models.Transaction, int64, error) {
	if page < 1 {
		page = 1
	}
	pageSize := 20
	offset := (page - 1) * pageSize
	return s.txRepo.GetByUserID(userID, offset, pageSize)
}

func (s *transactionService) PayService(userID uint, input PayServiceInput) (PayServiceResult, error) {
	wallet, err := s.walletRepo.GetByUserAndID(userID, input.WalletID)
	if err != nil {
		return PayServiceResult{}, ErrNotFound("Wallet not found")
	}

	amount := input.Amount
	if amount == 0 {
		plans, ok := servicePlans[input.ServiceProvider]
		if !ok {
			return PayServiceResult{}, ErrBadRequest("Unknown service provider")
		}
		base, ok := plans[input.Plan]
		if !ok {
			return PayServiceResult{}, ErrBadRequest("Unknown plan")
		}
		amount = base * 1.1
	}

	if wallet.Balance < amount {
		return PayServiceResult{}, ErrBadRequest("Insufficient balance")
	}

	s.walletRepo.UpdateBalance(wallet, wallet.Balance-amount)

	tx := &models.Transaction{
		UserID:          userID,
		WalletID:        wallet.ID,
		Amount:          amount,
		Nature:          "expense",
		Service:         input.Service,
		ServiceProvider: input.ServiceProvider,
		Currency:        wallet.Currency,
		Status:          models.StatusDone,
	}
	s.txRepo.Create(tx)

	return PayServiceResult{
		Amount:        amount,
		Currency:      wallet.Currency,
		NewBalance:    wallet.Balance - amount,
		TransactionID: tx.ID,
	}, nil
}

func (s *transactionService) TransferMoney(userID uint, input TransferInput) (TransferResult, error) {
	senderWallet, err := s.walletRepo.GetByUserAndID(userID, input.WalletID)
	if err != nil {
		return TransferResult{}, ErrNotFound("Wallet not found")
	}

	var recipientWallet *models.Wallet
	if input.Provider == "Nita" || input.Provider == "MPesa" {
		if senderWallet.PhoneNumber == input.Recipient {
			return TransferResult{}, ErrBadRequest("Cannot transfer to yourself")
		}
		recipientWallet, err = s.walletRepo.GetByPhoneAndProvider(input.Recipient, input.Provider)
		if err != nil {
			return TransferResult{}, ErrNotFound("Recipient wallet not found")
		}
	} else if input.Provider == "PayPal" {
		if senderWallet.Email == input.Recipient {
			return TransferResult{}, ErrBadRequest("Cannot transfer to yourself")
		}
		recipientWallet, err = s.walletRepo.GetByEmailAndProvider(input.Recipient, "PayPal")
		if err != nil {
			return TransferResult{}, ErrNotFound("Recipient wallet not found")
		}
	} else {
		return TransferResult{}, ErrBadRequest("Unsupported provider for transfer")
	}

	if senderWallet.Currency != recipientWallet.Currency {
		return TransferResult{}, ErrBadRequest("Currency mismatch between wallets")
	}
	if senderWallet.Balance < input.Amount {
		return TransferResult{}, ErrBadRequest("Insufficient balance")
	}

	s.walletRepo.UpdateBalance(senderWallet, senderWallet.Balance-input.Amount)

	outTx := &models.Transaction{
		UserID:          userID,
		WalletID:        senderWallet.ID,
		Amount:          input.Amount,
		Nature:          "expense",
		Service:         "transfer",
		ServiceProvider: input.Provider,
		Currency:        senderWallet.Currency,
		Status:          models.StatusDone,
	}
	s.txRepo.Create(outTx)

	s.walletRepo.UpdateBalance(recipientWallet, recipientWallet.Balance+input.Amount)

	inTx := &models.Transaction{
		UserID:               recipientWallet.UserID,
		WalletID:             recipientWallet.ID,
		Amount:               input.Amount,
		Nature:               "income",
		Service:              "transfer",
		ServiceProvider:      input.Provider,
		Currency:             recipientWallet.Currency,
		Status:               models.StatusDone,
		RelatedTransactionID: &outTx.ID,
	}
	s.txRepo.Create(inTx)
	s.txRepo.Update(outTx, map[string]interface{}{"related_transaction_id": inTx.ID})

	return TransferResult{
		Amount:        input.Amount,
		Currency:      senderWallet.Currency,
		NewBalance:    senderWallet.Balance - input.Amount,
		TransactionID: outTx.ID,
	}, nil
}
