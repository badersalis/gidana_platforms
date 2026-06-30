package repositories

import (
	"github.com/badersalis/gidana_backend/internal/models"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(tx *models.Transaction) error
	GetByID(id uint) (*models.Transaction, error)
	GetByCinetpayID(cinetpayID string) (*models.Transaction, error)
	Update(tx *models.Transaction, updates map[string]interface{}) error
	GetByUserID(userID uint) ([]models.Transaction, error)
}

type transactionRepository struct{ db *gorm.DB }

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(tx *models.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *transactionRepository) GetByID(id uint) (*models.Transaction, error) {
	var tx models.Transaction
	if err := r.db.First(&tx, id).Error; err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepository) GetByCinetpayID(cinetpayID string) (*models.Transaction, error) {
	var tx models.Transaction
	if err := r.db.Where("cinetpay_transaction_id = ?", cinetpayID).First(&tx).Error; err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepository) Update(tx *models.Transaction, updates map[string]interface{}) error {
	return r.db.Model(tx).Updates(updates).Error
}

func (r *transactionRepository) GetByUserID(userID uint) ([]models.Transaction, error) {
	var txs []models.Transaction
	if err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&txs).Error; err != nil {
		return nil, err
	}
	return txs, nil
}
