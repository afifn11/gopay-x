package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionType string
type TransactionStatus string

const (
	TxTypeTopUp    TransactionType = "top_up"
	TxTypeTransfer TransactionType = "transfer"
	TxTypeWithdraw TransactionType = "withdraw"
	TxTypeRefund   TransactionType = "refund"
	TxTypeFee      TransactionType = "fee"

	TxStatusPending  TransactionStatus = "pending"
	TxStatusSuccess  TransactionStatus = "success"
	TxStatusFailed   TransactionStatus = "failed"
	TxStatusReversed TransactionStatus = "reversed"
)

// Transaction adalah immutable ledger entry
type Transaction struct {
	ID            uuid.UUID         `gorm:"type:uuid;primary_key" json:"id"`
	UserID        uuid.UUID         `gorm:"type:uuid;not null;index" json:"user_id"`
	CounterpartID *uuid.UUID        `gorm:"type:uuid;index" json:"counterpart_id"` // receiver/sender
	Type          TransactionType   `gorm:"not null" json:"type"`
	Amount        int64             `gorm:"not null" json:"amount"`
	Fee           int64             `gorm:"default:0" json:"fee"`
	BalanceBefore int64             `gorm:"not null" json:"balance_before"`
	BalanceAfter  int64             `gorm:"not null" json:"balance_after"`
	Status        TransactionStatus `gorm:"default:success" json:"status"`
	ReferenceID   string            `gorm:"uniqueIndex;not null" json:"reference_id"`
	ServiceSource string            `gorm:"not null" json:"service_source"` // payment-service, wallet-service
	Description   string            `json:"description"`
	Metadata      string            `gorm:"type:jsonb" json:"metadata"`
	CreatedAt     time.Time         `gorm:"not null" json:"created_at"`
}

// Transaction bersifat immutable — tidak ada UpdatedAt
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}