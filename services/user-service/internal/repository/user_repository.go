package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/afifn11/gopay-x/services/user-service/internal/domain"
)

type userProfileRepository struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) domain.UserProfileRepository {
	return &userProfileRepository{db: db}
}

func (r *userProfileRepository) Create(ctx context.Context, profile *domain.UserProfile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

func (r *userProfileRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.UserProfile, error) {
	var profile domain.UserProfile
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &profile, err
}

func (r *userProfileRepository) FindByEmail(ctx context.Context, email string) (*domain.UserProfile, error) {
	var profile domain.UserProfile
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&profile).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &profile, err
}

func (r *userProfileRepository) Update(ctx context.Context, profile *domain.UserProfile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}

func (r *userProfileRepository) UpdateKYCStatus(ctx context.Context, userID uuid.UUID, status domain.KYCStatus, note string) error {
	return r.db.WithContext(ctx).Model(&domain.UserProfile{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"kyc_status": status,
			"kyc_note":   note,
		}).Error
}

func (r *userProfileRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&domain.UserProfile{}).Error
}