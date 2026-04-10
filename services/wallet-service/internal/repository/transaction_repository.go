package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/afifn11/gopay-x/services/wallet-service/internal/domain"
)

type walletTransactionRepository struct {
	db *gorm.DB
}

func NewWalletTransactionRepository(db *gorm.DB) domain.WalletTransactionRepository {
	return &walletTransactionRepository{db: db}
}

func (r *walletTransactionRepository) Create(ctx context.Context, tx *domain.WalletTransaction) error {
	return r.db.WithContext(ctx).Create(tx).Error
}

func (r *walletTransactionRepository) FindByReferenceID(ctx context.Context, referenceID string) (*domain.WalletTransaction, error) {
	var tx domain.WalletTransaction
	err := r.db.WithContext(ctx).Where("reference_id = ?", referenceID).First(&tx).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &tx, err
}

func (r *walletTransactionRepository) FindByWalletID(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]domain.WalletTransaction, int64, error) {
	var txs []domain.WalletTransaction
	var total int64

	r.db.WithContext(ctx).Model(&domain.WalletTransaction{}).
		Where("wallet_id = ?", walletID).Count(&total)

	err := r.db.WithContext(ctx).
		Where("wallet_id = ?", walletID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&txs).Error

	return txs, total, err
}

func (r *walletTransactionRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.WalletTransaction, int64, error) {
	var txs []domain.WalletTransaction
	var total int64

	r.db.WithContext(ctx).Model(&domain.WalletTransaction{}).
		Where("user_id = ?", userID).Count(&total)

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&txs).Error

	return txs, total, err
}

func (r *walletTransactionRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TransactionStatus) error {
	return r.db.WithContext(ctx).Model(&domain.WalletTransaction{}).
		Where("id = ?", id).
		Update("status", status).Error
}