package payments

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/soundmarket/backend/internal/config"
)

const yooKassaBaseURL = "https://api.yookassa.ru/v3"

type YooKassaProvider struct {
	cfg    *config.Config
	client *http.Client
}

type yooKassaCreatePaymentRequest struct {
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Capture      bool `json:"capture"`
	Confirmation struct {
		Type      string `json:"type"`
		ReturnURL string `json:"return_url"`
	} `json:"confirmation"`
	Description string            `json:"description,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

type yooKassaPaymentResponse struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	Paid         bool   `json:"paid"`
	Amount       struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Confirmation struct {
		Type            string `json:"type"`
		ConfirmationURL string `json:"confirmation_url"`
	} `json:"confirmation"`
	Metadata map[string]string `json:"metadata"`
}

type yooKassaWebhookPayload struct {
	Type   string                   `json:"type"`
	Event  string                   `json:"event"`
	Object yooKassaPaymentResponse  `json:"object"`
}

func NewYooKassaProvider(cfg *config.Config) *YooKassaProvider {
	return &YooKassaProvider{
		cfg: cfg,
		client: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

func (p *YooKassaProvider) CreatePayment(ctx context.Context, input CreatePaymentInput) (*PaymentSession, error) {
	var reqBody yooKassaCreatePaymentRequest
	reqBody.Amount.Value = fmt.Sprintf("%d.00", input.Amount)
	reqBody.Amount.Currency = "RUB"
	reqBody.Capture = true
	reqBody.Confirmation.Type = "redirect"
	reqBody.Confirmation.ReturnURL = p.cfg.YooKassaReturnURL
	reqBody.Description = input.Description
	reqBody.Metadata = map[string]string{
		"user_id": input.UserID,
		"kind":    "deposit",
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, yooKassaBaseURL+"/payments", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotence-Key", uuid.NewString())
	req.Header.Set("Authorization", "Basic "+p.basicAuthToken())

	respBody, statusCode, err := p.do(req)
	if err != nil {
		return nil, err
	}
	if statusCode < 200 || statusCode >= 300 {
		return nil, fmt.Errorf("yookassa create payment failed: status %d body %s", statusCode, respBody)
	}

	var paymentResp yooKassaPaymentResponse
	if err := json.Unmarshal(respBody, &paymentResp); err != nil {
		return nil, err
	}
	return &PaymentSession{
		ExternalID:      paymentResp.ID,
		Status:          paymentResp.Status,
		ConfirmationURL: paymentResp.Confirmation.ConfirmationURL,
		Provider:        "yookassa",
		RequestRaw:      string(payload),
		Raw:             string(respBody),
	}, nil
}

func (p *YooKassaProvider) GetPayment(ctx context.Context, externalID string) (*PaymentInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, yooKassaBaseURL+"/payments/"+externalID, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Basic "+p.basicAuthToken())

	respBody, statusCode, err := p.do(req)
	if err != nil {
		return nil, err
	}
	if statusCode < 200 || statusCode >= 300 {
		return nil, fmt.Errorf("yookassa get payment failed: status %d body %s", statusCode, respBody)
	}

	var paymentResp yooKassaPaymentResponse
	if err := json.Unmarshal(respBody, &paymentResp); err != nil {
		return nil, err
	}
	return &PaymentInfo{
		ExternalID: paymentResp.ID,
		Status:     paymentResp.Status,
		Paid:       paymentResp.Paid || paymentResp.Status == "succeeded",
		Amount:     parseRubAmount(paymentResp.Amount.Value),
		Provider:   "yookassa",
		Raw:        string(respBody),
	}, nil
}

func (p *YooKassaProvider) HandleWebhook(ctx context.Context, payload []byte, _ http.Header) (*WebhookResult, error) {
	var webhook yooKassaWebhookPayload
	if err := json.Unmarshal(payload, &webhook); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidWebhook, err)
	}
	if webhook.Object.ID == "" {
		return nil, fmt.Errorf("%w: yookassa webhook missing object.id", ErrInvalidWebhook)
	}

	paymentInfo, err := p.GetPayment(ctx, webhook.Object.ID)
	if err != nil {
		return nil, err
	}
	return &WebhookResult{
		ExternalID: paymentInfo.ExternalID,
		Status:     paymentInfo.Status,
		Paid:       paymentInfo.Paid,
		Amount:     paymentInfo.Amount,
		Provider:   "yookassa",
		Raw:        string(payload),
	}, nil
}

func (p *YooKassaProvider) do(req *http.Request) ([]byte, int, error) {
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return body, resp.StatusCode, nil
}

func (p *YooKassaProvider) basicAuthToken() string {
	raw := p.cfg.YooKassaShopID + ":" + p.cfg.YooKassaSecretKey
	return base64.StdEncoding.EncodeToString([]byte(raw))
}

func formatRubAmount(value int64) string {
	return fmt.Sprintf("%d.00", value)
}

func parseRubAmount(value string) int64 {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}
	parts := strings.SplitN(value, ".", 3)
	rubles := int64(0)
	for _, ch := range parts[0] {
		if ch < '0' || ch > '9' {
			return 0
		}
		rubles = rubles*10 + int64(ch-'0')
	}
	return rubles
}
