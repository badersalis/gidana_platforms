package handlers

import (
	"strconv"

	"github.com/badersalis/gidana_backend/internal/middleware"
	"github.com/badersalis/gidana_backend/internal/services"
	"github.com/badersalis/gidana_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service services.TransactionService
}

func NewTransactionHandler(svc services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: svc}
}

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	txs, total, err := h.service.GetTransactions(userID, page)
	if handleErr(c, err) {
		return
	}
	pageSize := 20
	if page < 1 {
		page = 1
	}
	utils.Paginated(c, txs, total, page, pageSize)
}

func (h *TransactionHandler) PayService(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var input struct {
		Service         string  `json:"service" binding:"required"`
		ServiceProvider string  `json:"service_provider" binding:"required"`
		Plan            string  `json:"plan"`
		WalletID        uint    `json:"wallet_id" binding:"required"`
		Amount          float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	result, err := h.service.PayService(userID, services.PayServiceInput{
		Service:         input.Service,
		ServiceProvider: input.ServiceProvider,
		Plan:            input.Plan,
		WalletID:        input.WalletID,
		Amount:          input.Amount,
	})
	if handleErr(c, err) {
		return
	}
	utils.OK(c, gin.H{
		"message":        "Payment successful",
		"amount":         result.Amount,
		"currency":       result.Currency,
		"new_balance":    result.NewBalance,
		"transaction_id": result.TransactionID,
	})
}

func (h *TransactionHandler) TransferMoney(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var input struct {
		WalletID  uint    `json:"wallet_id" binding:"required"`
		Recipient string  `json:"recipient" binding:"required"`
		Amount    float64 `json:"amount" binding:"required,gt=0"`
		Provider  string  `json:"provider" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	result, err := h.service.TransferMoney(userID, services.TransferInput{
		WalletID:  input.WalletID,
		Recipient: input.Recipient,
		Amount:    input.Amount,
		Provider:  input.Provider,
	})
	if handleErr(c, err) {
		return
	}
	utils.OK(c, gin.H{
		"message":        "Transfer successful",
		"amount":         result.Amount,
		"currency":       result.Currency,
		"new_balance":    result.NewBalance,
		"transaction_id": result.TransactionID,
	})
}
