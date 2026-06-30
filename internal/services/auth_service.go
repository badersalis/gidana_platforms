package services

import (
	"regexp"
	"strings"
	"time"

	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/repositories"
	"github.com/badersalis/gidana_backend/internal/utils"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^\+[1-9]\d{6,14}$`)
)

type RegisterInput struct {
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
	Password    string
}

type AuthService interface {
	Register(input RegisterInput) (*models.User, string, error)
	Login(identifier, password string) (*models.User, string, error)
	GetMe(userID uint) (*models.User, error)
}

type authService struct {
	repo repositories.UserRepository
}

func NewAuthService(repo repositories.UserRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) Register(input RegisterInput) (*models.User, string, error) {
	input.FirstName = strings.TrimSpace(input.FirstName)
	input.LastName = strings.TrimSpace(input.LastName)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.PhoneNumber = strings.TrimSpace(input.PhoneNumber)

	if len(input.FirstName) < 2 || len(input.FirstName) > 100 {
		return nil, "", ErrBadRequest("Le prénom doit avoir entre 2 et 100 caractères")
	}
	if len(input.LastName) < 2 || len(input.LastName) > 100 {
		return nil, "", ErrBadRequest("Le nom doit avoir entre 2 et 100 caractères")
	}
	if input.Email == "" && input.PhoneNumber == "" {
		return nil, "", ErrBadRequest("Email ou numéro de téléphone requis")
	}
	if input.Email != "" && !emailRegex.MatchString(input.Email) {
		return nil, "", ErrBadRequest("Format d'email invalide")
	}
	if input.PhoneNumber != "" && !phoneRegex.MatchString(input.PhoneNumber) {
		return nil, "", ErrBadRequest("Format de téléphone invalide. Utilisez le format international (+XXXXXXXXXXX)")
	}
	if len(input.Password) < 6 {
		return nil, "", ErrBadRequest("Le mot de passe doit avoir au moins 6 caractères")
	}

	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, "", err
	}

	user := &models.User{
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Email:        input.Email,
		PhoneNumber:  input.PhoneNumber,
		PasswordHash: hash,
		MemberSince:  time.Now(),
		Active:       true,
		Locale:       "fr",
	}

	if err := s.repo.Create(user); err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "UNIQUE") {
			return nil, "", ErrBadRequest("Cet email ou numéro de téléphone est déjà utilisé")
		}
		return nil, "", err
	}

	token, _ := utils.GenerateToken(user.ID, user.Email)
	return user, token, nil
}

func (s *authService) GetMe(userID uint) (*models.User, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, ErrNotFound("User not found")
	}
	if user.SubscriptionExpiresAt != nil && user.SubscriptionExpiresAt.Before(time.Now()) {
		_ = s.repo.Update(userID, map[string]interface{}{
			"subscription_plan":       "basic",
			"landlord_plan":           "free",
			"subscription_expires_at": nil,
		})
		user.SubscriptionPlan = "basic"
		user.LandlordPlan = "free"
		user.SubscriptionExpiresAt = nil
	}
	return user, nil
}

func (s *authService) Login(identifier, password string) (*models.User, string, error) {
	identifier = strings.TrimSpace(identifier)
	if strings.Contains(identifier, "@") {
		identifier = strings.ToLower(identifier)
	}

	user, err := s.repo.GetByIdentifier(identifier)
	if err != nil {
		return nil, "", ErrUnauthorized("Identifiants incorrects")
	}

	if !utils.CheckPassword(password, user.PasswordHash) {
		return nil, "", ErrUnauthorized("Identifiants incorrects")
	}

	token, _ := utils.GenerateToken(user.ID, user.Email)
	return user, token, nil
}
