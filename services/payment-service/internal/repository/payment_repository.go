package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/afifn11/gopay-x/services/payment-service/internal/domain"
)

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) domain.PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *paymentRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&payment).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &payment, err
}

func (r *paymentRepository) FindByReferenceID(ctx context.Context, referenceID string) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.WithContext(ctx).Where("reference_id = ?", referenceID).First(&payment).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &payment, err
}

func (r *paymentRepository) FindBySenderUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Payment, int64, error) {
	var payments []domain.Payment
	var total int64

	r.db.WithContext(ctx).Model(&domain.Payment{}).
		Where("sender_user_id = ?", userID).Count(&total)

	err := r.db.WithContext(ctx).
		Where("sender_user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&payments).Error

	return payments, total, err
}

func (r *paymentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.PaymentStatus, completedAt *time.Time, failureReason string) error {
	updates := map[string]interface{}{
		"status":         status,
		"failure_reason": failureReason,
	}
	if completedAt != nil {
		updates["completed_at"] = completedAt
	}
	return r.db.WithContext(ctx).Model(&domain.Payment{}).
		Where("id = ?", id).Updates(updates).Error
}