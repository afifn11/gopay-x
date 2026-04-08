package domain

import (
	"context"
	"github.com/google/uuid"
)

type UserProfileRepository interface {
	Create(ctx context.Context, profile *UserProfile) error
	FindByUserID(ctx context.Context, userID uuid.UUID) (*UserProfile, error)
	FindByEmail(ctx context.Context, email string) (*UserProfile, error)
	Update(ctx context.Context, profile *UserProfile) error
	UpdateKYCStatus(ctx context.Context, userID uuid.UUID, status KYCStatus, note string) error
	Delete(ctx context.Context, userID uuid.UUID) error
}

type KYCDocumentRepository interface {
	Create(ctx context.Context, doc *KYCDocument) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]KYCDocument, error)
	FindLatestByUserID(ctx context.Context, userID uuid.UUID) (*KYCDocument, error)
}