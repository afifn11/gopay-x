package usecase

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/afifn11/gopay-x/services/fraud-detection-service/internal/domain"
	"github.com/afifn11/gopay-x/services/fraud-detection-service/internal/rules"
)

var (
	ErrAlreadyChecked = errors.New("transaction already fraud-checked")
	ErrNotFound       = errors.New("fraud check not found")
)

type FraudUsecase interface {
	CheckTransaction(ctx context.Context, req *CheckRequest) (*domain.FraudCheck, error)
	GetChecksByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]domain.FraudCheck, int64, error)
	GetRiskProfile(ctx context.Context, userID uuid.UUID) (*domain.UserRiskProfile, error)
}

type CheckRequest struct {
	UserID      string
	ReferenceID string
	EventType   string
	Amount      int64
}

type fraudUsecase struct {
	fraudRepo   domain.FraudCheckRepository
	profileRepo domain.UserRiskProfileRepository
	engine      *rules.RuleEngine
}

func NewFraudUsecase(
	fraudRepo domain.FraudCheckRepository,
	profileRepo domain.UserRiskProfileRepository,
	engine *rules.RuleEngine,
) FraudUsecase {
	return &fraudUsecase{
		fraudRepo:   fraudRepo,
		profileRepo: profileRepo,
		engine:      engine,
	}
}

func (uc *fraudUsecase) CheckTransaction(ctx context.Context, req *CheckRequest) (*domain.FraudCheck, error) {
	// Idempotency
	existing, err := uc.fraudRepo.FindByReferenceID(ctx, req.ReferenceID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, ErrAlreadyChecked
	}

	userID, _ := uuid.Parse(req.UserID)

	// Run rule engine
	result, err := uc.engine.Evaluate(ctx, &rules.EvaluationInput{
		UserID:      userID,
		ReferenceID: req.ReferenceID,
		EventType:   req.EventType,
		Amount:      req.Amount,
	})
	if err != nil {
		return nil, err
	}

	check := &domain.FraudCheck{
		UserID:      userID,
		ReferenceID: req.ReferenceID,
		EventType:   req.EventType,
		Amount:      req.Amount,
		RiskScore:   result.RiskScore,
		RiskLevel:   result.RiskLevel,
		Status:      result.Status,
		Reasons:     rules.ReasonsToJSON(result.Reasons),
		CreatedAt:   time.Now(),
	}

	if err := uc.fraudRepo.Create(ctx, check); err != nil {
		return nil, err
	}

	// Update risk profile
	now := time.Now()
	profile := &domain.UserRiskProfile{
		UserID:    userID,
		UpdatedAt: now,
	}

	if result.Status == domain.FraudFlagged {
		profile.TotalFlaggedCount = 1
		profile.LastFlaggedAt = &now
		_ = uc.profileRepo.Upsert(ctx, profile)
		_ = uc.profileRepo.IncrementFlagged(ctx, userID)
		log.Printf("⚠️  Fraud flagged: user=%s ref=%s score=%d", req.UserID, req.ReferenceID, result.RiskScore)
	}

	if result.Status == domain.FraudBlocked {
		profile.IsBlocked = true
		profile.TotalBlockedCount = 1
		_ = uc.profileRepo.Upsert(ctx, profile)
		_ = uc.profileRepo.BlockUser(ctx, userID)
		log.Printf("🚫 Fraud BLOCKED: user=%s ref=%s score=%d", req.UserID, req.ReferenceID, result.RiskScore)
	}

	return check, nil
}

func (uc *fraudUsecase) GetChecksByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]domain.FraudCheck, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return uc.fraudRepo.FindByUserID(ctx, userID, limit, (page-1)*limit)
}

func (uc *fraudUsecase) GetRiskProfile(ctx context.Context, userID uuid.UUID) (*domain.UserRiskProfile, error) {
	profile, err := uc.profileRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return &domain.UserRiskProfile{UserID: userID}, nil
	}
	return profile, nil
}