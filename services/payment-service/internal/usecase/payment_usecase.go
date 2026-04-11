package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/afifn11/gopay-x/services/payment-service/internal/domain"
	"github.com/afifn11/gopay-x/services/payment-service/internal/gateway"
)

var (
	ErrPaymentNotFound      = errors.New("payment not found")
	ErrDuplicatePayment     = errors.New("duplicate payment request")
	ErrLockFailed           = errors.New("failed to acquire lock, please retry")
	ErrInvalidAmount        = errors.New("amount must be greater than zero")
	ErrSameUser             = errors.New("sender and receiver cannot be the same")
)

type PaymentUsecase interface {
	CreateTransfer(ctx context.Context, req *TransferRequest) (*domain.Payment, error)
	CreateTopUpViaGateway(ctx context.Context, req *GatewayTopUpRequest) (*domain.Payment, error)
	GetPaymentByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error)
	GetPaymentHistory(ctx context.Context, userID uuid.UUID, page, limit int) (*PaymentListResponse, error)
	HandleCallback(ctx context.Context, req *CallbackRequest) error
}

type TransferRequest struct {
	SenderUserID   string `json:"sender_user_id"`
	ReceiverUserID string `json:"receiver_user_id" validate:"required,uuid"`
	Amount         int64  `json:"amount" validate:"required,gt=0"`
	ReferenceID    string `json:"reference_id" validate:"required"`
	Description    string `json:"description"`
}

type GatewayTopUpRequest struct {
	UserID      string `json:"user_id"`
	Amount      int64  `json:"amount" validate:"required,gt=0"`
	ReferenceID string `json:"reference_id" validate:"required"`
	Description string `json:"description"`
}

type CallbackRequest struct {
	ExternalID string `json:"external_id"`
	Status     string `json:"status"`
	RawPayload string `json:"raw_payload"`
}

type PaymentListResponse struct {
	Data       []domain.Payment `json:"data"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
}

type paymentUsecase struct {
	paymentRepo  domain.PaymentRepository
	callbackRepo domain.PaymentCallbackRepository
	lockRepo     domain.LockRepository
	gateway      gateway.PaymentGateway
}

func NewPaymentUsecase(
	paymentRepo domain.PaymentRepository,
	callbackRepo domain.PaymentCallbackRepository,
	lockRepo domain.LockRepository,
	gw gateway.PaymentGateway,
) PaymentUsecase {
	return &paymentUsecase{
		paymentRepo:  paymentRepo,
		callbackRepo: callbackRepo,
		lockRepo:     lockRepo,
		gateway:      gw,
	}
}

func (uc *paymentUsecase) CreateTransfer(ctx context.Context, req *TransferRequest) (*domain.Payment, error) {
	if req.Amount <= 0 {
		return nil, ErrInvalidAmount
	}
	if req.SenderUserID == req.ReceiverUserID {
		return nil, ErrSameUser
	}

	// Idempotency check
	existing, err := uc.paymentRepo.FindByReferenceID(ctx, req.ReferenceID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, ErrDuplicatePayment
	}

	senderID, _ := uuid.Parse(req.SenderUserID)
	receiverID, _ := uuid.Parse(req.ReceiverUserID)

	// Distributed lock
	lockKey := fmt.Sprintf("payment:transfer:%s", senderID.String())
	acquired, err := uc.lockRepo.AcquireLock(ctx, lockKey, 15)
	if err != nil || !acquired {
		return nil, ErrLockFailed
	}
	defer uc.lockRepo.ReleaseLock(ctx, lockKey)

	now := time.Now()
	payment := &domain.Payment{
		SenderUserID:   senderID,
		ReceiverUserID: &receiverID,
		Type:           domain.PaymentTypeTransfer,
		Method:         domain.MethodWallet,
		Amount:         req.Amount,
		Fee:            0,
		TotalAmount:    req.Amount,
		Status:         domain.PaymentSuccess,
		ReferenceID:    req.ReferenceID,
		Description:    req.Description,
		CompletedAt:    &now,
	}

	if err := uc.paymentRepo.Create(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (uc *paymentUsecase) CreateTopUpViaGateway(ctx context.Context, req *GatewayTopUpRequest) (*domain.Payment, error) {
	if req.Amount <= 0 {
		return nil, ErrInvalidAmount
	}

	// Idempotency check
	existing, err := uc.paymentRepo.FindByReferenceID(ctx, req.ReferenceID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, ErrDuplicatePayment
	}

	userID, _ := uuid.Parse(req.UserID)
	expiredAt := time.Now().Add(24 * time.Hour)

	// Call mock gateway
	gwResp, err := uc.gateway.CreateTransaction(ctx, &gateway.GatewayRequest{
		OrderID:     req.ReferenceID,
		Amount:      req.Amount,
		Description: req.Description,
		CustomerID:  req.UserID,
	})
	if err != nil {
		return nil, err
	}

	payment := &domain.Payment{
		SenderUserID: userID,
		Type:         domain.PaymentTypeTopUp,
		Method:       domain.MethodVirtualAccount,
		Amount:       req.Amount,
		Fee:          0,
		TotalAmount:  req.Amount,
		Status:       domain.PaymentPending,
		ReferenceID:  req.ReferenceID,
		ExternalID:   gwResp.ExternalID,
		Description:  req.Description,
		ExpiredAt:    &expiredAt,
	}

	if err := uc.paymentRepo.Create(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (uc *paymentUsecase) GetPaymentByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error) {
	payment, err := uc.paymentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}
	return payment, nil
}

func (uc *paymentUsecase) GetPaymentHistory(ctx context.Context, userID uuid.UUID, page, limit int) (*PaymentListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	payments, total, err := uc.paymentRepo.FindBySenderUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &PaymentListResponse{
		Data:       payments,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (uc *paymentUsecase) HandleCallback(ctx context.Context, req *CallbackRequest) error {
	payment, err := uc.paymentRepo.FindByReferenceID(ctx, req.ExternalID)
	if err != nil || payment == nil {
		return ErrPaymentNotFound
	}

	// Store callback log
	callback := &domain.PaymentCallback{
		PaymentID:  payment.ID,
		ExternalID: req.ExternalID,
		Status:     req.Status,
		RawPayload: req.RawPayload,
		ReceivedAt: time.Now(),
	}
	_ = uc.callbackRepo.Create(ctx, callback)

	// Update payment status
	var status domain.PaymentStatus
	var completedAt *time.Time
	now := time.Now()

	switch req.Status {
	case "success", "settlement":
		status = domain.PaymentSuccess
		completedAt = &now
	case "failure", "cancel", "expire":
		status = domain.PaymentFailed
	default:
		return nil
	}

	return uc.paymentRepo.UpdateStatus(ctx, payment.ID, status, completedAt, "")
}