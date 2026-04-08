package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KYCStatus string
type Gender string

const (
	KYCPending  KYCStatus = "pending"
	KYCVerified KYCStatus = "verified"
	KYCRejected KYCStatus = "rejected"

	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

type UserProfile struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID  `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	FullName    string     `gorm:"not null" json:"full_name"`
	Email       string     `gorm:"uniqueIndex;not null" json:"email"`
	Phone       string     `gorm:"uniqueIndex" json:"phone"`
	Gender      *Gender    `json:"gender"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Address     string     `json:"address"`
	City        string     `json:"city"`
	Country     string     `gorm:"default:ID" json:"country"`
	AvatarURL   string     `json:"avatar_url"`
	KYCStatus   KYCStatus  `gorm:"default:pending" json:"kyc_status"`
	KYCNote     string     `json:"kyc_note"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (u *UserProfile) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type KYCDocument struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	DocumentType string    `gorm:"not null" json:"document_type"` // ktp, passport, sim
	DocumentURL  string    `gorm:"not null" json:"document_url"`
	SubmittedAt  time.Time `json:"submitted_at"`
	ReviewedAt   *time.Time `json:"reviewed_at"`
	CreatedAt    time.Time `json:"created_at"`
}

func (k *KYCDocument) BeforeCreate(tx *gorm.DB) error {
	if k.ID == uuid.Nil {
		k.ID = uuid.New()
	}
	return nil
}