package repositories

import (
	"github.com/badersalis/gidana_backend/internal/models"
	"gorm.io/gorm"
)

type WalletRepository interface {
	GetByUserID(userID uint) ([]models.Wallet, error)
	GetByUserAndID(userID, id uint) (*models.Wallet, error)
	GetByPhoneAndProvider(phone, provider string) (*models.Wallet, error)
	GetByEmailAndProvider(email, provider string) (*models.Wallet, error)
	Create(wallet *models.Wallet) error
	Update(wallet *models.Wallet, updates map[string]interface{}) error
	UpdateBalance(wallet *models.Wallet, balance float64) error
	Delete(wallet *models.Wallet) error
	DeselectAllForUser(userID uint) error
	DeselectAllExcept(userID, excludeID uint) error
}

type walletRepository struct{ db *gorm.DB }

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) GetByUserID(userID uint) ([]models.Wallet, error) {
	var wallets []models.Wallet
	err := r.db.Where("user_id = ?", userID).Find(&wallets).Error
	return wallets, err
}

func (r *walletRepository) GetByUserAndID(userID, id uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepository) GetByPhoneAndProvider(phone, provider string) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Where("phone_number = ? AND provider = ?", phone, provider).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepository) GetByEmailAndProvider(email, provider string) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Where("email = ? AND provider = ?", email, provider).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepository) Create(wallet *models.Wallet) error {
	return r.db.Create(wallet).Error
}

func (r *walletRepository) Update(wallet *models.Wallet, updates map[string]interface{}) error {
	return r.db.Model(wallet).Updates(updates).Error
}

func (r *walletRepository) UpdateBalance(wallet *models.Wallet, balance float64) error {
	return r.db.Model(wallet).Update("balance", balance).Error
}

func (r *walletRepository) Delete(wallet *models.Wallet) error {
	return r.db.Delete(wallet).Error
}

func (r *walletRepository) DeselectAllForUser(userID uint) error {
	return r.db.Model(&models.Wallet{}).Where("user_id = ?", userID).Update("selected", false).Error
}

func (r *walletRepository) DeselectAllExcept(userID, excludeID uint) error {
	return r.db.Model(&models.Wallet{}).Where("user_id = ? AND id != ?", userID, excludeID).Update("selected", false).Error
}
