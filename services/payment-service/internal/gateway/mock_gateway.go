package gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PaymentGateway interface {
	CreateTransaction(ctx context.Context, req *GatewayRequest) (*GatewayResponse, error)
	CheckStatus(ctx context.Context, externalID string) (*GatewayResponse, error)
	Refund(ctx context.Context, externalID string, amount int64) error
}

type GatewayRequest struct {
	OrderID     string
	Amount      int64
	Description string
	CustomerID  string
}

type GatewayResponse struct {
	ExternalID    string
	Status        string
	PaymentURL    string
	VirtualAccount string
	ExpiredAt     time.Time
}

// MockGateway — simulasi Midtrans/Xendit untuk development
type MockGateway struct{}

func NewMockGateway() PaymentGateway {
	return &MockGateway{}
}

func (g *MockGateway) CreateTransaction(ctx context.Context, req *GatewayRequest) (*GatewayResponse, error) {
	return &GatewayResponse{
		ExternalID:     fmt.Sprintf("mock-ext-%s", uuid.New().String()),
		Status:         "pending",
		PaymentURL:     fmt.Sprintf("https://mock-payment.gopay-x.dev/pay/%s", req.OrderID),
		VirtualAccount: fmt.Sprintf("8808%010d", req.Amount),
		ExpiredAt:      time.Now().Add(24 * time.Hour),
	}, nil
}

func (g *MockGateway) CheckStatus(ctx context.Context, externalID string) (*GatewayResponse, error) {
	// Simulasi: semua transaksi otomatis success di mock mode
	return &GatewayResponse{
		ExternalID: externalID,
		Status:     "success",
	}, nil
}

func (g *MockGateway) Refund(ctx context.Context, externalID string, amount int64) error {
	return nil
}