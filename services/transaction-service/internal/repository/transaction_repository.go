package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/afifn11/gopay-x/services/transaction-service/internal/domain"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) domain.TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	return r.db.WithContext(ctx).Create(tx).Error
}

func (r *transactionRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	var tx domain.Transaction
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&tx).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &tx, err
}

func (r *transactionRepository) FindByReferenceID(ctx context.Context, referenceID string) (*domain.Transaction, error) {
	var tx domain.Transaction
	err := r.db.WithContext(ctx).Where("reference_id = ?", referenceID).First(&tx).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &tx, err
}

func (r *transactionRepository) FindByUserID(ctx context.Context, userID uuid.UUID, filter *domain.TransactionFilter) ([]domain.Transaction, int64, error) {
	var txs []domain.Transaction
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Transaction{}).Where("user_id = ?", userID)

	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.From != nil {
		query = query.Where("created_at >= ?", filter.From)
	}
	if filter.To != nil {
		query = query.Where("created_at <= ?", filter.To)
	}

	query.Count(&total)

	err := query.Order("created_at DESC").
		Limit(filter.Limit).
		Offset(filter.Offset).
		Find(&txs).Error

	return txs, total, err
}

func (r *transactionRepository) GetSummary(ctx context.Context, userID uuid.UUID, from, to time.Time) (*domain.TransactionSummary, error) {
	var summary domain.TransactionSummary

	r.db.WithContext(ctx).Model(&domain.Transaction{}).
		Where("user_id = ? AND created_at BETWEEN ? AND ? AND type IN ('top_up','refund')", userID, from, to).
		Select("COALESCE(SUM(amount), 0) as total_in, COUNT(*) as tx_count").
		Scan(&summary)

	r.db.WithContext(ctx).Model(&domain.Transaction{}).
		Where("user_id = ? AND created_at BETWEEN ? AND ? AND type IN ('transfer','withdraw')", userID, from, to).
		Select("COALESCE(SUM(amount), 0) as total_out").
		Scan(&summary)

	r.db.WithContext(ctx).Model(&domain.Transaction{}).
		Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, from, to).
		Select("COALESCE(SUM(fee), 0) as total_fee").
		Scan(&summary)

	return &summary, nil
}