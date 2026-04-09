package payments

import (
	"context"
	"errors"
	"net/http"
)

var ErrInvalidWebhook = errors.New("invalid webhook")

type CreatePaymentInput struct {
	UserID      string
	Amount      int64
	Description string
}

type PaymentSession struct {
	ExternalID      string
	Status          string
	ConfirmationURL string
	Provider        string
	Raw             string
	RequestRaw      string
}

type PaymentInfo struct {
	ExternalID string
	Status     string
	Paid       bool
	Amount     int64
	Provider   string
	Raw        string
}

type WebhookResult struct {
	ExternalID string
	Status     string
	Paid       bool
	Amount     int64
	Provider   string
	Raw        string
}

type Provider interface {
	CreatePayment(ctx context.Context, input CreatePaymentInput) (*PaymentSession, error)
	GetPayment(ctx context.Context, externalID string) (*PaymentInfo, error)
	HandleWebhook(ctx context.Context, payload []byte, headers http.Header) (*WebhookResult, error)
}
