package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/badersalis/gidana_backend/internal/config"
	"github.com/badersalis/gidana_backend/internal/services"
	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	subSvc services.SubscriptionService
}

func NewWebhookHandler(subSvc services.SubscriptionService) *WebhookHandler {
	return &WebhookHandler{subSvc: subSvc}
}

// HandleCinetPay godoc
// @Summary      CinetPay payment webhook
// @Tags         webhooks
// @Accept       json
// @Produce      json
// @Param        X-CINETPAY-SIGNATURE  header  string  true  "HMAC-SHA256 signature of raw request body"
// @Success      200  {object}  MessageResponse
// @Failure      400  {object}  ErrorResponse
// @Router       /webhooks/cinetpay [post]
func (h *WebhookHandler) HandleCinetPay(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Failed to read request body"})
		return
	}

	// Validate HMAC-SHA256 signature
	secret := config.App.CinetPayAPISecret
	if secret != "" {
		sig := c.GetHeader("X-CINETPAY-SIGNATURE")
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(body)
		expected := hex.EncodeToString(mac.Sum(nil))
		if !hmac.Equal([]byte(sig), []byte(expected)) {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid signature"})
			return
		}
	}

	var payload struct {
		CpmTransID     string `json:"cpm_trans_id"`
		CpmTransStatus string `json:"cpm_trans_status"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid payload"})
		return
	}

	if err := h.subSvc.ActivateFromWebhook(payload.CpmTransID, payload.CpmTransStatus); err != nil {
		if se, ok := services.IsServiceError(err); ok && se.Code == 404 {
			// Unknown transaction — still acknowledge to prevent retries
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "Transaction not tracked"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Webhook processing failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Webhook processed"})
}
