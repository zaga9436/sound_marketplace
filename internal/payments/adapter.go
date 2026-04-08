package payments

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/soundmarket/backend/internal/config"
)

type PaymentResult struct {
	ExternalID  string
	Status      string
	RedirectURL string
}

type Adapter interface {
	CreatePayment(userID string, amount int64) (*PaymentResult, error)
}

type MockYooKassaAdapter struct {
	cfg *config.Config
}

func NewMockYooKassaAdapter(cfg *config.Config) *MockYooKassaAdapter {
	return &MockYooKassaAdapter{cfg: cfg}
}

func (a *MockYooKassaAdapter) CreatePayment(userID string, amount int64) (*PaymentResult, error) {
	externalID := uuid.NewString()
	return &PaymentResult{
		ExternalID:  externalID,
		Status:      "pending",
		RedirectURL: fmt.Sprintf("%s?payment_id=%s&user_id=%s&amount=%d", a.cfg.YooKassaReturnURL, externalID, userID, amount),
	}, nil
}

