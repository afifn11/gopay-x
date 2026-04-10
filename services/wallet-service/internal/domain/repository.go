package domain

import (
	"context"

	"github.com/google/uuid"
)

type WalletRepository interface {
	Create(ctx context.Context, wallet *Wallet) error
	FindByUserID(ctx context.Context, userID uuid.UUID) (*Wallet, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Wallet, error)
	UpdateBalance(ctx context.Context, walletID uuid.UUID, newBalance int64) error
	UpdateStatus(ctx context.Context, walletID uuid.UUID, status WalletStatus) error
}

type WalletTransactionRepository interface {
	Create(ctx context.Context, tx *WalletTransaction) error
	FindByReferenceID(ctx context.Context, referenceID string) (*WalletTransaction, error)
	FindByWalletID(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]WalletTransaction, int64, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]WalletTransaction, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status TransactionStatus) error
}

type LockRepository interface {
	AcquireLock(ctx context.Context, key string, ttlSeconds int) (bool, error)
	ReleaseLock(ctx context.Context, key string) error
}