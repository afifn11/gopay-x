package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/afifn11/gopay-x/services/wallet-service/internal/domain"
)

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) domain.WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) Create(ctx context.Context, wallet *domain.Wallet) error {
	return r.db.WithContext(ctx).Create(wallet).Error
}

func (r *walletRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.Wallet, error) {
	var wallet domain.Wallet
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&wallet).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &wallet, err
}

func (r *walletRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Wallet, error) {
	var wallet domain.Wallet
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&wallet).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &wallet, err
}

func (r *walletRepository) UpdateBalance(ctx context.Context, walletID uuid.UUID, newBalance int64) error {
	return r.db.WithContext(ctx).Model(&domain.Wallet{}).
		Where("id = ?", walletID).
		Update("balance", newBalance).Error
}

func (r *walletRepository) UpdateStatus(ctx context.Context, walletID uuid.UUID, status domain.WalletStatus) error {
	return r.db.WithContext(ctx).Model(&domain.Wallet{}).
		Where("id = ?", walletID).
		Update("status", status).Error
}