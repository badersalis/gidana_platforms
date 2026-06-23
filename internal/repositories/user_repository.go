package repositories

import (
	"github.com/badersalis/gidana_backend/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByIdentifier(identifier string) (*models.User, error)
	Update(userID uint, updates map[string]interface{}) error
	GetByIDWithToken(id uint) (*models.User, error)
	HasPendingDeletion(userID uint) (bool, error)
	CreateDeletedAccountSnapshot(snap *models.DeletedAccount) error
	Deactivate(user *models.User) error
	SoftDelete(user *models.User) error
}

type userRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByIdentifier(identifier string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ? OR phone_number = ?", identifier, identifier).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(userID uint, updates map[string]interface{}) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
}

func (r *userRepository) GetByIDWithToken(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Select("expo_push_token").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) HasPendingDeletion(userID uint) (bool, error) {
	var existing models.DeletedAccount
	err := r.db.Where("user_id = ?", userID).First(&existing).Error
	if err == nil {
		return true, nil
	}
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	return false, err
}

func (r *userRepository) CreateDeletedAccountSnapshot(snap *models.DeletedAccount) error {
	return r.db.Create(snap).Error
}

func (r *userRepository) Deactivate(user *models.User) error {
	return r.db.Model(user).Update("active", false).Error
}

func (r *userRepository) SoftDelete(user *models.User) error {
	return r.db.Delete(user).Error
}
