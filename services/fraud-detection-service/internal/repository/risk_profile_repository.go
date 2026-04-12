package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/afifn11/gopay-x/services/fraud-detection-service/internal/domain"
)

type userRiskProfileRepository struct {
	db *gorm.DB
}

func NewUserRiskProfileRepository(db *gorm.DB) domain.UserRiskProfileRepository {
	return &userRiskProfileRepository{db: db}
}

func (r *userRiskProfileRepository) Upsert(ctx context.Context, profile *domain.UserRiskProfile) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"total_flagged_count", "total_blocked_count", "last_flagged_at", "is_blocked", "updated_at"}),
		}).Create(profile).Error
}

func (r *userRiskProfileRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.UserRiskProfile, error) {
	var profile domain.UserRiskProfile
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &profile, err
}

func (r *userRiskProfileRepository) IncrementFlagged(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&domain.UserRiskProfile{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"total_flagged_count": gorm.Expr("total_flagged_count + 1"),
			"last_flagged_at":     now,
			"updated_at":          now,
		}).Error
}

func (r *userRiskProfileRepository) BlockUser(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.UserRiskProfile{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"is_blocked":          true,
			"total_blocked_count": gorm.Expr("total_blocked_count + 1"),
			"updated_at":          time.Now(),
		}).Error
}