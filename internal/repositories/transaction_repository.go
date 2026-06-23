package repositories

import (
	"github.com/badersalis/gidana_backend/internal/models"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	GetByUserID(userID uint, offset, limit int) ([]models.Transaction, int64, error)
	Create(tx *models.Transaction) error
	Update(tx *models.Transaction, updates map[string]interface{}) error
}

type transactionRepository struct{ db *gorm.DB }

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) GetByUserID(userID uint, offset, limit int) ([]models.Transaction, int64, error) {
	var total int64
	r.db.Model(&models.Transaction{}).Where("user_id = ?", userID).Count(&total)

	var txs []models.Transaction
	err := r.db.Where("user_id = ?", userID).
		Order("created_at desc").
		Offset(offset).Limit(limit).
		Preload("Wallet").
		Find(&txs).Error
	return txs, total, err
}

func (r *transactionRepository) Create(tx *models.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *transactionRepository) Update(tx *models.Transaction, updates map[string]interface{}) error {
	return r.db.Model(tx).Updates(updates).Error
}
