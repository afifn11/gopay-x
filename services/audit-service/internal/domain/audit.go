package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditLog struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ServiceName string    `gorm:"not null;index" json:"service_name"`
	EventType   string    `gorm:"not null;index" json:"event_type"`
	ActorID     string    `gorm:"index" json:"actor_id"`
	ActorType   string    `json:"actor_type"` // user, system, admin
	ResourceID  string    `gorm:"index" json:"resource_id"`
	ResourceType string   `json:"resource_type"`
	Action      string    `gorm:"not null" json:"action"`
	Status      string    `gorm:"not null" json:"status"` // success, failed
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	Payload     string    `gorm:"type:jsonb" json:"payload"`
	ErrorMsg    string    `json:"error_msg"`
	CreatedAt   time.Time `gorm:"not null;index" json:"created_at"`
}

// AuditLog bersifat immutable — tidak ada update/delete
func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}