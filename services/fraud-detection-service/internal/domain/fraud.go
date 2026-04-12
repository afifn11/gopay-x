package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RiskLevel string
type FraudStatus string

const (
	RiskLow      RiskLevel = "low"
	RiskMedium   RiskLevel = "medium"
	RiskHigh     RiskLevel = "high"
	RiskCritical RiskLevel = "critical"

	FraudFlagged  FraudStatus = "flagged"
	FraudCleared  FraudStatus = "cleared"
	FraudBlocked  FraudStatus = "blocked"
)

type FraudCheck struct {
	ID          uuid.UUID   `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID   `gorm:"type:uuid;not null;index" json:"user_id"`
	ReferenceID string      `gorm:"uniqueIndex;not null" json:"reference_id"`
	EventType   string      `gorm:"not null" json:"event_type"`
	Amount      int64       `gorm:"not null" json:"amount"`
	RiskScore   int         `gorm:"not null" json:"risk_score"` // 0-100
	RiskLevel   RiskLevel   `gorm:"not null" json:"risk_level"`
	Status      FraudStatus `gorm:"not null" json:"status"`
	Reasons     string      `gorm:"type:jsonb" json:"reasons"` // JSON array of triggered rules
	CreatedAt   time.Time   `json:"created_at"`
}

func (f *FraudCheck) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

type UserRiskProfile struct {
	ID                  uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID              uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	TotalFlaggedCount   int       `gorm:"default:0" json:"total_flagged_count"`
	TotalBlockedCount   int       `gorm:"default:0" json:"total_blocked_count"`
	LastFlaggedAt       *time.Time `json:"last_flagged_at"`
	IsBlocked           bool      `gorm:"default:false" json:"is_blocked"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (u *UserRiskProfile) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}