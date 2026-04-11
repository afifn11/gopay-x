package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/afifn11/gopay-x/services/transaction-service/internal/domain"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrDuplicateRecord     = errors.New("transaction already recorded")
)

type TransactionUsecase interface {
	RecordTransaction(ctx context.Context, req *RecordTransactionRequest) (*domain.Transaction, error)
	GetTransaction(ctx context.Context, id uuid.UUID) (*domain.Transaction, error)
	GetHistory(ctx context.Context, userID uuid.UUID, req *HistoryRequest) (*HistoryResponse, error)
	GetSummary(ctx context.Context, userID uuid.UUID, from, to string) (*domain.TransactionSummary, error)
}

type RecordTransactionRequest struct {
	UserID        string
	CounterpartID string
	Type          string
	Amount        int64
	Fee           int64
	BalanceBefore int64
	BalanceAfter  int64
	ReferenceID   string
	ServiceSource string
	Description   string
	Status        string
}

type HistoryRequest struct {
	Type   string
	Status string
	From   string
	To     string
	Page   int
	Limit  int
}

type HistoryResponse struct {
	Data       []domain.Transaction `json:"data"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	Limit      int                  `json:"limit"`
	TotalPages int                  `json:"total_pages"`
}

type transactionUsecase struct {
	txRepo domain.TransactionRepository
}

func NewTransactionUsecase(txRepo domain.TransactionRepository) TransactionUsecase {
	return &transactionUsecase{txRepo: txRepo}
}

func (uc *transactionUsecase) RecordTransaction(ctx context.Context, req *RecordTransactionRequest) (*domain.Transaction, error) {
	// Idempotency check
	existing, err := uc.txRepo.FindByReferenceID(ctx, req.ReferenceID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, ErrDuplicateRecord
	}

	userID, _ := uuid.Parse(req.UserID)

	tx := &domain.Transaction{
		UserID:        userID,
		Type:          domain.TransactionType(req.Type),
		Amount:        req.Amount,
		Fee:           req.Fee,
		BalanceBefore: req.BalanceBefore,
		BalanceAfter:  req.BalanceAfter,
		Status:        domain.TransactionStatus(req.Status),
		ReferenceID:   req.ReferenceID,
		ServiceSource: req.ServiceSource,
		Description:   req.Description,
		CreatedAt:     time.Now(),
	}

	if req.CounterpartID != "" {
		cid, _ := uuid.Parse(req.CounterpartID)
		tx.CounterpartID = &cid
	}

	if err := uc.txRepo.Create(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (uc *transactionUsecase) GetTransaction(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	tx, err := uc.txRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, ErrTransactionNotFound
	}
	return tx, nil
}

func (uc *transactionUsecase) GetHistory(ctx context.Context, userID uuid.UUID, req *HistoryRequest) (*HistoryResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}

	filter := &domain.TransactionFilter{
		Type:   req.Type,
		Status: req.Status,
		Limit:  req.Limit,
		Offset: (req.Page - 1) * req.Limit,
	}

	if req.From != "" {
		t, err := time.Parse("2006-01-02", req.From)
		if err == nil {
			filter.From = &t
		}
	}
	if req.To != "" {
		t, err := time.Parse("2006-01-02", req.To)
		if err == nil {
			filter.To = &t
		}
	}

	txs, total, err := uc.txRepo.FindByUserID(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	return &HistoryResponse{
		Data:       txs,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

func (uc *transactionUsecase) GetSummary(ctx context.Context, userID uuid.UUID, from, to string) (*domain.TransactionSummary, error) {
	fromTime, err := time.Parse("2006-01-02", from)
	if err != nil {
		fromTime = time.Now().AddDate(0, -1, 0)
	}
	toTime, err := time.Parse("2006-01-02", to)
	if err != nil {
		toTime = time.Now()
	}
	return uc.txRepo.GetSummary(ctx, userID, fromTime, toTime)
}