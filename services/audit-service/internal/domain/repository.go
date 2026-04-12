package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AuditLogRepository interface {
	Create(ctx context.Context, log *AuditLog) error
	FindByID(ctx context.Context, id uuid.UUID) (*AuditLog, error)
	FindByActorID(ctx context.Context, actorID string, limit, offset int) ([]AuditLog, int64, error)
	FindByResourceID(ctx context.Context, resourceID string, limit, offset int) ([]AuditLog, int64, error)
	FindByServiceAndEvent(ctx context.Context, service, event string, from, to time.Time, limit, offset int) ([]AuditLog, int64, error)
}