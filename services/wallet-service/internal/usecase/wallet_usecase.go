package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/afifn11/gopay-x/services/wallet-service/internal/domain"
)

var (
	ErrWalletNotFound      = errors.New("wallet not found")
	ErrWalletAlreadyExists = errors.New("wallet already exists")
	ErrWalletSuspended     = errors.New("wallet is suspended")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrDuplicateTransaction = errors.New("duplicate transaction")
	ErrLockFailed          = errors.New("failed to acquire lock, please retry")
	ErrInvalidAmount       = errors.New("amount must be greater than zero")
)

type WalletUsecase interface {
	CreateWallet(ctx context.Context, userID uuid.UUID) (*domain.Wallet, error)
	GetWallet(ctx context.Context, userID uuid.UUID) (*domain.Wallet, error)
	TopUp(ctx context.Context, req *TopUpRequest) (*domain.WalletTransaction, error)
	GetTransactionHistory(ctx context.Context, userID uuid.UUID, page, limit int) (*TransactionListResponse, error)
}

type TopUpRequest struct {
	UserID      string `json:"user_id" validate:"required,uuid"`
	Amount      int64  `json:"amount" validate:"required,gt=0"`
	ReferenceID string `json:"reference_id" validate:"required"`
	Description string `json:"description"`
}

type TransactionListResponse struct {
	Data       []domain.WalletTransaction `json:"data"`
	Total      int64                      `json:"total"`
	Page       int                        `json:"page"`
	Limit      int                        `json:"limit"`
	TotalPages int                        `json:"total_pages"`
}

type walletUsecase struct {
	walletRepo domain.WalletRepository
	txRepo     domain.WalletTransactionRepository
	lockRepo   domain.LockRepository
}

func NewWalletUsecase(
	walletRepo domain.WalletRepository,
	txRepo domain.WalletTransactionRepository,
	lockRepo domain.LockRepository,
) WalletUsecase {
	return &walletUsecase{
		walletRepo: walletRepo,
		txRepo:     txRepo,
		lockRepo:   lockRepo,
	}
}

func (uc *walletUsecase) CreateWallet(ctx context.Context, userID uuid.UUID) (*domain.Wallet, error) {
	existing, err := uc.walletRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrWalletAlreadyExists
	}

	wallet := &domain.Wallet{
		UserID:  userID,
		Balance: 0,
		Status:  domain.WalletActive,
	}

	if err := uc.walletRepo.Create(ctx, wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (uc *walletUsecase) GetWallet(ctx context.Context, userID uuid.UUID) (*domain.Wallet, error) {
	wallet, err := uc.walletRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, ErrWalletNotFound
	}
	return wallet, nil
}

func (uc *walletUsecase) TopUp(ctx context.Context, req *TopUpRequest) (*domain.WalletTransaction, error) {
	if req.Amount <= 0 {
		return nil, ErrInvalidAmount
	}

	// Idempotency check — cegah double top-up
	existing, err := uc.txRepo.FindByReferenceID(ctx, req.ReferenceID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, ErrDuplicateTransaction
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, err
	}

	// Distributed lock — cegah race condition
	lockKey := fmt.Sprintf("wallet:topup:%s", userID.String())
	acquired, err := uc.lockRepo.AcquireLock(ctx, lockKey, 10)
	if err != nil || !acquired {
		return nil, ErrLockFailed
	}
	defer uc.lockRepo.ReleaseLock(ctx, lockKey)

	wallet, err := uc.walletRepo.FindByUserID(ctx, userID)
	if err != nil || wallet == nil {
		return nil, ErrWalletNotFound
	}

	if wallet.Status != domain.WalletActive {
		return nil, ErrWalletSuspended
	}

	balanceBefore := wallet.Balance
	balanceAfter := balanceBefore + req.Amount

	// Update balance
	if err := uc.walletRepo.UpdateBalance(ctx, wallet.ID, balanceAfter); err != nil {
		return nil, err
	}

	// Record transaction
	tx := &domain.WalletTransaction{
		WalletID:      wallet.ID,
		UserID:        userID,
		Type:          domain.TypeTopUp,
		Amount:        req.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Status:        domain.StatusSuccess,
		ReferenceID:   req.ReferenceID,
		Description:   req.Description,
		CreatedAt:     time.Now(),
	}

	if err := uc.txRepo.Create(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (uc *walletUsecase) GetTransactionHistory(ctx context.Context, userID uuid.UUID, page, limit int) (*TransactionListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	txs, total, err := uc.txRepo.FindByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &TransactionListResponse{
		Data:       txs,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}