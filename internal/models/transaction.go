package models

import "time"

type TransactionStatus string

const (
	TransactionDone    TransactionStatus = "done"
	TransactionFailed  TransactionStatus = "failed"
	TransactionOngoing TransactionStatus = "ongoing"
)

type Transaction struct {
	ID                    uint              `gorm:"primarykey" json:"id"`
	CreatedAt             time.Time         `json:"created_at"`
	UserID                uint              `gorm:"not null;index" json:"user_id"`
	Amount                float64           `gorm:"not null" json:"amount"`
	Currency              string            `gorm:"size:5;default:'XOF'" json:"currency"`
	Nature                string            `gorm:"size:20;default:'debit'" json:"nature"`
	Service               string            `gorm:"size:50;not null" json:"service"`
	Plan                  string            `gorm:"size:50" json:"plan,omitempty"`
	Status                TransactionStatus `gorm:"size:20;default:'ongoing'" json:"status"`
	CinetpayTransactionID string            `gorm:"size:100;index" json:"cinetpay_transaction_id,omitempty"`
}
