package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/afifn11/gopay-x/services/fraud-detection-service/internal/domain"
)

type fraudCheckRepository struct {
	db *gorm.DB
}

func NewFraudCheckRepository(db *gorm.DB) domain.FraudCheckRepository {
	return &fraudCheckRepository{db: db}
}

func (r *fraudCheckRepository) Create(ctx context.Context, check *domain.FraudCheck) error {
	return r.db.WithContext(ctx).Create(check).Error
}

func (r *fraudCheckRepository) FindByReferenceID(ctx context.Context, referenceID string) (*domain.FraudCheck, error) {
	var check domain.FraudCheck
	err := r.db.WithContext(ctx).Where("reference_id = ?", referenceID).First(&check).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &check, err
}

func (r *fraudCheckRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.FraudCheck, int64, error) {
	var checks []domain.FraudCheck
	var total int64

	r.db.WithContext(ctx).Model(&domain.FraudCheck{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&checks).Error

	return checks, total, err
}

func (r *fraudCheckRepository) CountRecentByUserID(ctx context.Context, userID uuid.UUID, minutes int) (int64, error) {
	var count int64
	since := time.Now().Add(-time.Duration(minutes) * time.Minute)
	err := r.db.WithContext(ctx).Model(&domain.FraudCheck{}).
		Where("user_id = ? AND created_at >= ?", userID, since).
		Count(&count).Error
	return count, err
}

func (r *fraudCheckRepository) SumRecentAmountByUserID(ctx context.Context, userID uuid.UUID, minutes int) (int64, error) {
	var total int64
	since := time.Now().Add(-time.Duration(minutes) * time.Minute)
	err := r.db.WithContext(ctx).Model(&domain.FraudCheck{}).
		Where("user_id = ? AND created_at >= ?", userID, since).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}