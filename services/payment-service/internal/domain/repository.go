package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *Payment) error
	FindByID(ctx context.Context, id uuid.UUID) (*Payment, error)
	FindByReferenceID(ctx context.Context, referenceID string) (*Payment, error)
	FindBySenderUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]Payment, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status PaymentStatus, completedAt *time.Time, failureReason string) error
}

type PaymentCallbackRepository interface {
	Create(ctx context.Context, callback *PaymentCallback) error
}

type LockRepository interface {
	AcquireLock(ctx context.Context, key string, ttlSeconds int) (bool, error)
	ReleaseLock(ctx context.Context, key string) error
}