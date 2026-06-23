package handlers

import (
	"github.com/badersalis/gidana_backend/internal/middleware"
	"github.com/badersalis/gidana_backend/internal/services"
	"github.com/badersalis/gidana_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(svc services.UserService) *UserHandler {
	return &UserHandler{service: svc}
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var input struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Gender      string `json:"gender"`
		DateOfBirth string `json:"date_of_birth"`
		Locale      string `json:"locale"`
		Timezone    string `json:"timezone"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	user, err := h.service.UpdateProfile(userID, services.UpdateProfileInput{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Gender:      input.Gender,
		DateOfBirth: input.DateOfBirth,
		Locale:      input.Locale,
		Timezone:    input.Timezone,
	})
	if handleErr(c, err) {
		return
	}
	utils.OK(c, user)
}

func (h *UserHandler) UploadProfilePicture(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	file, err := c.FormFile("picture")
	if err != nil {
		utils.BadRequest(c, "No file provided")
		return
	}

	user, err := h.service.UploadProfilePicture(userID, file)
	if handleErr(c, err) {
		return
	}
	utils.OK(c, user)
}

func (h *UserHandler) UpdatePushToken(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var input struct {
		Token string `json:"expo_push_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if handleErr(c, h.service.UpdatePushToken(userID, input.Token)) {
		return
	}
	utils.OK(c, gin.H{"message": "Push token updated"})
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var input struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if handleErr(c, h.service.ChangePassword(userID, input.CurrentPassword, input.NewPassword)) {
		return
	}
	utils.OK(c, gin.H{"message": "Password changed successfully"})
}

func (h *UserHandler) RequestDeleteAccount(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	if handleErr(c, h.service.RequestDeleteAccount(userID)) {
		return
	}
	utils.OK(c, gin.H{"message": "Account deletion request submitted. You will be notified once reviewed."})
}
