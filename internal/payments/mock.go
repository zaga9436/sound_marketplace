package payments

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/soundmarket/backend/internal/config"
)

type MockProvider struct {
	cfg *config.Config
}

type mockWebhookPayload struct {
	ExternalID string `json:"external_id"`
	Status     string `json:"status"`
	Amount     int64  `json:"amount"`
}

func NewMockProvider(cfg *config.Config) *MockProvider {
	return &MockProvider{cfg: cfg}
}

func (p *MockProvider) CreatePayment(_ context.Context, input CreatePaymentInput) (*PaymentSession, error) {
	externalID := uuid.NewString()
	requestRaw := fmt.Sprintf(`{"amount":{"value":"%d.00","currency":"RUB"},"description":"Balance deposit"}`, input.Amount)
	return &PaymentSession{
		ExternalID:      externalID,
		Status:          "pending",
		ConfirmationURL: fmt.Sprintf("%s?payment_id=%s&user_id=%s&amount=%d", p.cfg.YooKassaReturnURL, externalID, input.UserID, input.Amount),
		Provider:        "mock",
		RequestRaw:      requestRaw,
		Raw:             fmt.Sprintf(`{"id":"%s","status":"pending","amount":{"value":"%s","currency":"RUB"}}`, externalID, formatRubAmount(input.Amount)),
	}, nil
}

func (p *MockProvider) GetPayment(_ context.Context, externalID string) (*PaymentInfo, error) {
	return &PaymentInfo{
		ExternalID: externalID,
		Status:     "succeeded",
		Paid:       true,
		Provider:   "mock",
	}, nil
}

func (p *MockProvider) HandleWebhook(_ context.Context, payload []byte, _ http.Header) (*WebhookResult, error) {
	var req mockWebhookPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidWebhook, err)
	}
	status := req.Status
	if status == "" {
		status = "succeeded"
	}
	return &WebhookResult{
		ExternalID: req.ExternalID,
		Status:     status,
		Paid:       status == "succeeded",
		Amount:     req.Amount,
		Provider:   "mock",
		Raw:        string(payload),
	}, nil
}
