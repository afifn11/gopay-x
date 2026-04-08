package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/afifn11/gopay-x/services/user-service/internal/domain"
)

type kycDocumentRepository struct {
	db *gorm.DB
}

func NewKYCDocumentRepository(db *gorm.DB) domain.KYCDocumentRepository {
	return &kycDocumentRepository{db: db}
}

func (r *kycDocumentRepository) Create(ctx context.Context, doc *domain.KYCDocument) error {
	return r.db.WithContext(ctx).Create(doc).Error
}

func (r *kycDocumentRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.KYCDocument, error) {
	var docs []domain.KYCDocument
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&docs).Error
	return docs, err
}

func (r *kycDocumentRepository) FindLatestByUserID(ctx context.Context, userID uuid.UUID) (*domain.KYCDocument, error) {
	var doc domain.KYCDocument
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}