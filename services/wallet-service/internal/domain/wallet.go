package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletStatus string
type TransactionType string
type TransactionStatus string

const (
	WalletActive   WalletStatus = "active"
	WalletSuspended WalletStatus = "suspended"
	WalletClosed   WalletStatus = "closed"

	TypeTopUp    TransactionType = "top_up"
	TypeTransfer TransactionType = "transfer"
	TypeWithdraw TransactionType = "withdraw"
	TypeRefund   TransactionType = "refund"

	StatusPending   TransactionStatus = "pending"
	StatusSuccess   TransactionStatus = "success"
	StatusFailed    TransactionStatus = "failed"
	StatusReversed  TransactionStatus = "reversed"
)

type Wallet struct {
	ID        uuid.UUID    `gorm:"type:uuid;primary_key" json:"id"`
	UserID    uuid.UUID    `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	Balance   int64        `gorm:"not null;default:0" json:"balance"` // stored in cents/rupiah
	Status    WalletStatus `gorm:"default:active" json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func (w *Wallet) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}

type WalletTransaction struct {
	ID              uuid.UUID         `gorm:"type:uuid;primary_key" json:"id"`
	WalletID        uuid.UUID         `gorm:"type:uuid;not null;index" json:"wallet_id"`
	UserID          uuid.UUID         `gorm:"type:uuid;not null;index" json:"user_id"`
	Type            TransactionType   `gorm:"not null" json:"type"`
	Amount          int64             `gorm:"not null" json:"amount"`
	BalanceBefore   int64             `gorm:"not null" json:"balance_before"`
	BalanceAfter    int64             `gorm:"not null" json:"balance_after"`
	Status          TransactionStatus `gorm:"default:pending" json:"status"`
	ReferenceID     string            `gorm:"uniqueIndex;not null" json:"reference_id"` // idempotency key
	Description     string            `json:"description"`
	Metadata        string            `gorm:"type:jsonb" json:"metadata"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

func (t *WalletTransaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}