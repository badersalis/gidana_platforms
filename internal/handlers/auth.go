package handlers

import (
	"github.com/badersalis/gidana_backend/internal/services"
	"github.com/badersalis/gidana_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler(svc services.AuthService) *AuthHandler {
	return &AuthHandler{service: svc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input struct {
		FirstName   string `json:"first_name" binding:"required"`
		LastName    string `json:"last_name" binding:"required"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
		Password    string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Tous les champs obligatoires doivent être remplis")
		return
	}

	user, token, err := h.service.Register(services.RegisterInput{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
		Password:    input.Password,
	})
	if handleErr(c, err) {
		return
	}
	utils.Created(c, gin.H{"user": user, "token": token})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input struct {
		Identifier string `json:"identifier" binding:"required"`
		Password   string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Email/téléphone et mot de passe requis")
		return
	}

	user, token, err := h.service.Login(input.Identifier, input.Password)
	if handleErr(c, err) {
		return
	}
	utils.OK(c, gin.H{"user": user, "token": token})
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	user, _ := c.Get("user")
	utils.OK(c, user)
}
