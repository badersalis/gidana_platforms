package handlers

import (
	"github.com/badersalis/gidana_backend/internal/middleware"
	"github.com/badersalis/gidana_backend/internal/services"
	"github.com/badersalis/gidana_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	service services.WalletService
}

func NewWalletHandler(svc services.WalletService) *WalletHandler {
	return &WalletHandler{service: svc}
}

type walletInput struct {
	Provider       string `json:"provider" binding:"required"`
	Nature         string `json:"nature"`
	PhoneNumber    string `json:"phone_number"`
	Email          string `json:"email"`
	CardNumber     string `json:"card_number"`
	CVV            string `json:"cvv"`
	ExpirationDate string `json:"expiration_date"`
	Password       string `json:"password"`
	Currency       string `json:"currency"`
	Selected       bool   `json:"selected"`
}

func (h *WalletHandler) GetWallets(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	wallets, err := h.service.GetWallets(userID)
	if handleErr(c, err) {
		return
	}
	utils.OK(c, wallets)
}

func (h *WalletHandler) CreateWallet(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var input walletInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	wallet, err := h.service.CreateWallet(userID, services.WalletInput{
		Provider: input.Provider, Nature: input.Nature,
		PhoneNumber: input.PhoneNumber, Email: input.Email,
		CardNumber: input.CardNumber, CVV: input.CVV,
		ExpirationDate: input.ExpirationDate, Password: input.Password,
		Currency: input.Currency, Selected: input.Selected,
	})
	if handleErr(c, err) {
		return
	}
	utils.Created(c, wallet)
}

func (h *WalletHandler) UpdateWallet(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	walletID := paramUint(c, "id")

	var input walletInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	wallet, err := h.service.UpdateWallet(userID, walletID, services.WalletInput{
		Provider: input.Provider, Nature: input.Nature,
		PhoneNumber: input.PhoneNumber, Email: input.Email,
		Currency: input.Currency, Selected: input.Selected,
	})
	if handleErr(c, err) {
		return
	}
	utils.OK(c, wallet)
}

func (h *WalletHandler) DeleteWallet(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	walletID := paramUint(c, "id")

	if handleErr(c, h.service.DeleteWallet(userID, walletID)) {
		return
	}
	utils.OK(c, gin.H{"message": "Wallet deleted"})
}

func (h *WalletHandler) SelectWallet(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	walletID := paramUint(c, "id")

	if handleErr(c, h.service.SelectWallet(userID, walletID)) {
		return
	}
	utils.OK(c, gin.H{"message": "Default wallet updated"})
}

func (h *WalletHandler) RefreshWalletBalance(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	walletID := paramUint(c, "id")

	balance, currency, err := h.service.RefreshWalletBalance(userID, walletID)
	if handleErr(c, err) {
		return
	}
	utils.OK(c, gin.H{"balance": balance, "currency": currency})
}
