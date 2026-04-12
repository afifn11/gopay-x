package domain

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KYCStatus string
type UserRole string

const (
	KYCPending  KYCStatus = "pending"
	KYCVerified KYCStatus = "verified"
	KYCRejected KYCStatus = "rejected"

	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type User struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	FullName    string     `gorm:"not null" json:"full_name"`
	Email       string     `gorm:"uniqueIndex;not null" json:"email"`
	Phone       string     `gorm:"uniqueIndex" json:"phone"`
	Password    string     `gorm:"not null" json:"-"`
	Role        UserRole   `gorm:"default:user" json:"role"`
	KYCStatus   KYCStatus  `gorm:"default:pending" json:"kyc_status"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}