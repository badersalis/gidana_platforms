package services

import (
	"mime/multipart"
	"time"

	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/repositories"
	"github.com/badersalis/gidana_backend/internal/storage"
	"github.com/badersalis/gidana_backend/internal/utils"
)

type UpdateProfileInput struct {
	FirstName   string
	LastName    string
	Gender      string
	DateOfBirth string
	Locale      string
	Timezone    string
}

type UserService interface {
	UpdateProfile(userID uint, input UpdateProfileInput) (*models.User, error)
	UploadProfilePicture(userID uint, file *multipart.FileHeader) (*models.User, error)
	UpdatePushToken(userID uint, token string) error
	ChangePassword(userID uint, currentPassword, newPassword string) error
	RequestDeleteAccount(userID uint) error
}

type userService struct {
	repo      repositories.UserRepository
	fileStore storage.FileStorage
}

func NewUserService(repo repositories.UserRepository, fileStore storage.FileStorage) UserService {
	return &userService{repo: repo, fileStore: fileStore}
}

func (s *userService) UpdateProfile(userID uint, input UpdateProfileInput) (*models.User, error) {
	updates := map[string]interface{}{}
	if input.FirstName != "" {
		updates["first_name"] = input.FirstName
	}
	if input.LastName != "" {
		updates["last_name"] = input.LastName
	}
	if input.Gender != "" {
		updates["gender"] = input.Gender
	}
	if input.Locale != "" {
		updates["locale"] = input.Locale
	}
	if input.Timezone != "" {
		updates["timezone"] = input.Timezone
	}
	if input.DateOfBirth != "" {
		t, err := time.Parse("2006-01-02", input.DateOfBirth)
		if err == nil {
			updates["date_of_birth"] = t
		}
	}

	if err := s.repo.Update(userID, updates); err != nil {
		return nil, ErrInternal("Failed to update profile")
	}
	return s.repo.GetByID(userID)
}

func (s *userService) UploadProfilePicture(userID uint, file *multipart.FileHeader) (*models.User, error) {
	url, err := s.fileStore.Save(file)
	if err != nil {
		return nil, ErrBadRequest(err.Error())
	}
	if err := s.repo.Update(userID, map[string]interface{}{"profile_picture": url}); err != nil {
		return nil, ErrInternal("Failed to update profile picture")
	}
	return s.repo.GetByID(userID)
}

func (s *userService) UpdatePushToken(userID uint, token string) error {
	return s.repo.Update(userID, map[string]interface{}{"expo_push_token": token})
}

func (s *userService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return ErrNotFound("User not found")
	}
	if !utils.CheckPassword(currentPassword, user.PasswordHash) {
		return ErrBadRequest("Current password is incorrect")
	}
	hash, _ := utils.HashPassword(newPassword)
	return s.repo.Update(userID, map[string]interface{}{"password_hash": hash})
}

func (s *userService) RequestDeleteAccount(userID uint) error {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return ErrNotFound("User not found")
	}

	has, err := s.repo.HasPendingDeletion(userID)
	if err != nil {
		return ErrInternal("Failed to process deletion request")
	}
	if has {
		return ErrBadRequest("Account deletion already requested")
	}

	snapshot := &models.DeletedAccount{
		UserID:         user.ID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Email:          user.Email,
		PhoneNumber:    user.PhoneNumber,
		ProfilePicture: user.ProfilePicture,
		Gender:         user.Gender,
		DateOfBirth:    user.DateOfBirth,
		MemberSince:    user.MemberSince,
		Locale:         user.Locale,
		Timezone:       user.Timezone,
		RequestedAt:    time.Now(),
		Status:         "pending",
	}

	if err := s.repo.CreateDeletedAccountSnapshot(snapshot); err != nil {
		return ErrInternal("Failed to process deletion request")
	}

	s.repo.Deactivate(user)
	s.repo.SoftDelete(user)

	utils.SendExpoPush(
		user.ExpoPushToken,
		"Deletion Request Received",
		"Your account deletion request has been received. Our compliance team will review it and notify you once the process is complete.",
		map[string]any{"type": "account_deletion_requested"},
	)
	return nil
}
