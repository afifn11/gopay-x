package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *Transaction) error
	FindByID(ctx context.Context, id uuid.UUID) (*Transaction, error)
	FindByReferenceID(ctx context.Context, referenceID string) (*Transaction, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, filter *TransactionFilter) ([]Transaction, int64, error)
	GetSummary(ctx context.Context, userID uuid.UUID, from, to time.Time) (*TransactionSummary, error)
}

type TransactionFilter struct {
	Type   string
	Status string
	From   *time.Time
	To     *time.Time
	Limit  int
	Offset int
}

type TransactionSummary struct {
	TotalIn    int64 `json:"total_in"`
	TotalOut   int64 `json:"total_out"`
	TotalFee   int64 `json:"total_fee"`
	TxCount    int64 `json:"tx_count"`
}