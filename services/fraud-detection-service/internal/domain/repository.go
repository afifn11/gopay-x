package domain

import (
	"context"

	"github.com/google/uuid"
)

type FraudCheckRepository interface {
	Create(ctx context.Context, check *FraudCheck) error
	FindByReferenceID(ctx context.Context, referenceID string) (*FraudCheck, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]FraudCheck, int64, error)
	CountRecentByUserID(ctx context.Context, userID uuid.UUID, minutes int) (int64, error)
	SumRecentAmountByUserID(ctx context.Context, userID uuid.UUID, minutes int) (int64, error)
}

type UserRiskProfileRepository interface {
	Upsert(ctx context.Context, profile *UserRiskProfile) error
	FindByUserID(ctx context.Context, userID uuid.UUID) (*UserRiskProfile, error)
	IncrementFlagged(ctx context.Context, userID uuid.UUID) error
	BlockUser(ctx context.Context, userID uuid.UUID) error
}