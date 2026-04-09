package service

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/soundmarket/backend/internal/apierr"
	"github.com/soundmarket/backend/internal/config"
	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/notifications"
	"github.com/soundmarket/backend/internal/payments"
	"github.com/soundmarket/backend/internal/repository"
)

type PaymentService struct {
	cfg      *config.Config
	store    repository.Store
	adapter  payments.Provider
	notifier notifications.Service
}

type PaymentSyncResult struct {
	Payment        domain.Payment `json:"payment"`
	DepositCreated bool           `json:"deposit_created"`
}

func NewPaymentService(cfg *config.Config, store repository.Store, adapter payments.Provider, notifier notifications.Service) *PaymentService {
	return &PaymentService{cfg: cfg, store: store, adapter: adapter, notifier: notifier}
}

func (s *PaymentService) CreateDeposit(user domain.User, amount int64) (domain.Payment, error) {
	if amount <= 0 {
		return domain.Payment{}, apierr.BadRequest("amount must be positive")
	}
	result, err := s.adapter.CreatePayment(context.Background(), payments.CreatePaymentInput{
		UserID:      user.ID,
		Amount:      amount,
		Description: "Balance deposit",
	})
	if err != nil {
		return domain.Payment{}, err
	}
	payment, err := s.store.CreatePayment(domain.Payment{
		UserID:      user.ID,
		ExternalID:  result.ExternalID,
		Amount:      amount,
		Status:      result.Status,
		Provider:    result.Provider,
		RedirectURL: result.ConfirmationURL,
		CallbackData: result.RequestRaw,
	})
	if err != nil {
		return domain.Payment{}, err
	}
	payment.ConfirmationURL = payment.RedirectURL
	return payment, nil
}

func (s *PaymentService) ProcessWebhook(ctx context.Context, payload []byte, headers http.Header) (domain.Payment, bool, error) {
	result, err := s.adapter.HandleWebhook(ctx, payload, headers)
	if err != nil {
		if errors.Is(err, payments.ErrInvalidWebhook) {
			return domain.Payment{}, false, apierr.BadRequest("invalid webhook payload")
		}
		return domain.Payment{}, false, err
	}
	if result.ExternalID == "" {
		return domain.Payment{}, false, apierr.BadRequest("payment id missing in webhook")
	}

	updated, created, err := s.applyProviderResult(result.ExternalID, result.Status, result.Paid, result.Raw)
	if err != nil {
		return domain.Payment{}, false, err
	}
	return updated, created, nil
}

func (s *PaymentService) SyncPayment(ctx context.Context, actor domain.User, externalID string) (PaymentSyncResult, error) {
	if strings.EqualFold(s.cfg.AppEnv, "production") {
		return PaymentSyncResult{}, apierr.Forbidden("manual payment sync is disabled in production")
	}
	payment, err := s.store.GetPaymentByExternalID(externalID)
	if err != nil {
		return PaymentSyncResult{}, apierr.NotFound("payment not found")
	}
	if actor.Role != domain.RoleAdmin && actor.ID != payment.UserID {
		return PaymentSyncResult{}, apierr.Forbidden("forbidden")
	}
	if payment.Provider != "yookassa" {
		return PaymentSyncResult{}, apierr.BadRequest("manual sync is supported only for yookassa payments")
	}

	info, err := s.adapter.GetPayment(ctx, externalID)
	if err != nil {
		return PaymentSyncResult{}, err
	}
	updated, created, err := s.applyProviderResult(info.ExternalID, info.Status, info.Paid, info.Raw)
	if err != nil {
		return PaymentSyncResult{}, err
	}
	return PaymentSyncResult{
		Payment:        updated,
		DepositCreated: created,
	}, nil
}

func (s *PaymentService) Refund(order domain.Order, amount int64) (domain.Transaction, error) {
	if amount <= 0 || amount > order.Amount {
		return domain.Transaction{}, apierr.BadRequest("invalid refund amount")
	}
	txType := domain.TransactionTypeRefund
	if amount < order.Amount {
		txType = domain.TransactionTypePartialRefund
	}
	return s.store.CreateTransaction(domain.Transaction{
		UserID:  order.CustomerID,
		OrderID: order.ID,
		Type:    txType,
		Amount:  amount,
	})
}

func (s *PaymentService) Balance(userID string) (int64, error) {
	return s.store.GetBalance(userID)
}

func (s *PaymentService) applyProviderResult(externalID, status string, paid bool, raw string) (domain.Payment, bool, error) {
	var (
		updated domain.Payment
		created bool
		err     error
	)

	err = s.store.WithTx(func(tx repository.Store) error {
		var payment domain.Payment
		payment, err = tx.GetPaymentByExternalID(externalID)
		if err != nil {
			return apierr.NotFound("payment not found")
		}
		if payment.Status == "succeeded" {
			updated = payment
			updated.ConfirmationURL = updated.RedirectURL
			return nil
		}
		if status != "succeeded" || !paid {
			updated = payment
			updated.Status = status
			updated.ConfirmationURL = updated.RedirectURL
			return nil
		}
		updated, err = tx.MarkPaymentSucceeded(externalID)
		if err != nil {
			return err
		}
		updated.CallbackData = raw
		if _, err := tx.CreateTransaction(domain.Transaction{
			UserID:     payment.UserID,
			Type:       domain.TransactionTypeDeposit,
			Amount:     payment.Amount,
			ExternalID: payment.ExternalID,
		}); err != nil {
			return err
		}
		updated.ConfirmationURL = updated.RedirectURL
		created = true
		return nil
	})
	if err != nil {
		return domain.Payment{}, false, err
	}
	if created {
		s.notifier.Publish(updated.UserID, "balance_deposit", "Balance replenished")
	}
	return updated, created, nil
}
