package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentStatus string
type PaymentType string
type PaymentMethod string

const (
	PaymentPending  PaymentStatus = "pending"
	PaymentSuccess  PaymentStatus = "success"
	PaymentFailed   PaymentStatus = "failed"
	PaymentExpired  PaymentStatus = "expired"
	PaymentRefunded PaymentStatus = "refunded"

	PaymentTypeTransfer PaymentType = "transfer"
	PaymentTypeTopUp    PaymentType = "top_up"
	PaymentTypeWithdraw PaymentType = "withdraw"

	MethodWallet      PaymentMethod = "wallet"
	MethodBankTransfer PaymentMethod = "bank_transfer"
	MethodVirtualAccount PaymentMethod = "virtual_account"
)

type Payment struct {
	ID              uuid.UUID     `gorm:"type:uuid;primary_key" json:"id"`
	SenderUserID    uuid.UUID     `gorm:"type:uuid;not null;index" json:"sender_user_id"`
	ReceiverUserID  *uuid.UUID    `gorm:"type:uuid;index" json:"receiver_user_id"`
	Type            PaymentType   `gorm:"not null" json:"type"`
	Method          PaymentMethod `gorm:"not null" json:"method"`
	Amount          int64         `gorm:"not null" json:"amount"`
	Fee             int64         `gorm:"default:0" json:"fee"`
	TotalAmount     int64         `gorm:"not null" json:"total_amount"`
	Status          PaymentStatus `gorm:"default:pending" json:"status"`
	ReferenceID     string        `gorm:"uniqueIndex;not null" json:"reference_id"`
	ExternalID      string        `gorm:"index" json:"external_id"`
	Description     string        `json:"description"`
	FailureReason   string        `json:"failure_reason"`
	ExpiredAt       *time.Time    `json:"expired_at"`
	CompletedAt     *time.Time    `json:"completed_at"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type PaymentCallback struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	PaymentID   uuid.UUID `gorm:"type:uuid;not null;index" json:"payment_id"`
	ExternalID  string    `json:"external_id"`
	Status      string    `json:"status"`
	RawPayload  string    `gorm:"type:jsonb" json:"raw_payload"`
	ReceivedAt  time.Time `json:"received_at"`
}

func (p *PaymentCallback) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}