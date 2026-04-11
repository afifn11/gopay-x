package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/afifn11/gopay-x/services/payment-service/internal/domain"
)

type paymentCallbackRepository struct {
	db *gorm.DB
}

func NewPaymentCallbackRepository(db *gorm.DB) domain.PaymentCallbackRepository {
	return &paymentCallbackRepository{db: db}
}

func (r *paymentCallbackRepository) Create(ctx context.Context, callback *domain.PaymentCallback) error {
	return r.db.WithContext(ctx).Create(callback).Error
}