package rules

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/afifn11/gopay-x/services/fraud-detection-service/internal/domain"
)

// RuleEngine mengevaluasi setiap transaksi terhadap rule-rule fraud
type RuleEngine struct {
	fraudRepo domain.FraudCheckRepository
}

func NewRuleEngine(fraudRepo domain.FraudCheckRepository) *RuleEngine {
	return &RuleEngine{fraudRepo: fraudRepo}
}

type EvaluationInput struct {
	UserID      uuid.UUID
	ReferenceID string
	EventType   string
	Amount      int64
}

type EvaluationResult struct {
	RiskScore int
	RiskLevel domain.RiskLevel
	Status    domain.FraudStatus
	Reasons   []string
}

func (e *RuleEngine) Evaluate(ctx context.Context, input *EvaluationInput) (*EvaluationResult, error) {
	score := 0
	var reasons []string

	// Rule 1: Large single transaction (> 10 juta)
	if input.Amount > 10_000_000 {
		score += 30
		reasons = append(reasons, "large_transaction: amount exceeds 10,000,000")
	}

	// Rule 2: Very large transaction (> 50 juta)
	if input.Amount > 50_000_000 {
		score += 40
		reasons = append(reasons, "very_large_transaction: amount exceeds 50,000,000")
	}

	// Rule 3: Velocity check — lebih dari 5 transaksi dalam 10 menit
	recentCount, err := e.fraudRepo.CountRecentByUserID(ctx, input.UserID, 10)
	if err == nil && recentCount >= 5 {
		score += 25
		reasons = append(reasons, "velocity_check: more than 5 transactions in 10 minutes")
	}

	// Rule 4: High volume dalam 10 menit (> 20 juta)
	recentAmount, err := e.fraudRepo.SumRecentAmountByUserID(ctx, input.UserID, 10)
	if err == nil && recentAmount > 20_000_000 {
		score += 20
		reasons = append(reasons, "high_volume: total amount exceeds 20,000,000 in 10 minutes")
	}

	// Rule 5: Round number (sering dipakai di penipuan)
	if input.Amount%1_000_000 == 0 && input.Amount >= 5_000_000 {
		score += 5
		reasons = append(reasons, "round_number: suspiciously round large amount")
	}

	// Cap score at 100
	if score > 100 {
		score = 100
	}

	result := &EvaluationResult{
		RiskScore: score,
		Reasons:   reasons,
	}

	// Determine risk level and status
	switch {
	case score >= 80:
		result.RiskLevel = domain.RiskCritical
		result.Status = domain.FraudBlocked
	case score >= 50:
		result.RiskLevel = domain.RiskHigh
		result.Status = domain.FraudFlagged
	case score >= 25:
		result.RiskLevel = domain.RiskMedium
		result.Status = domain.FraudFlagged
	default:
		result.RiskLevel = domain.RiskLow
		result.Status = domain.FraudCleared
	}

	return result, nil
}

func ReasonsToJSON(reasons []string) string {
	b, _ := json.Marshal(reasons)
	return string(b)
}