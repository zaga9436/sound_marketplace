package service

import (
	"github.com/soundmarket/backend/internal/apierr"
	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/notifications"
	"github.com/soundmarket/backend/internal/payments"
	"github.com/soundmarket/backend/internal/repository"
)

type PaymentService struct {
	store    repository.Store
	adapter  payments.Adapter
	notifier notifications.Service
}

func NewPaymentService(store repository.Store, adapter payments.Adapter, notifier notifications.Service) *PaymentService {
	return &PaymentService{store: store, adapter: adapter, notifier: notifier}
}

func (s *PaymentService) CreateDeposit(user domain.User, amount int64) (domain.Payment, error) {
	if amount <= 0 {
		return domain.Payment{}, apierr.BadRequest("amount must be positive")
	}
	result, err := s.adapter.CreatePayment(user.ID, amount)
	if err != nil {
		return domain.Payment{}, err
	}
	return s.store.CreatePayment(domain.Payment{
		UserID:      user.ID,
		ExternalID:  result.ExternalID,
		Amount:      amount,
		Status:      result.Status,
		Provider:    "yookassa_mock",
		RedirectURL: result.RedirectURL,
	})
}

func (s *PaymentService) ProcessWebhook(externalID string) (domain.Transaction, error) {
	var created domain.Transaction
	err := s.store.WithTx(func(tx repository.Store) error {
		payment, err := tx.GetPaymentByExternalID(externalID)
		if err != nil {
			return apierr.NotFound("payment not found")
		}
		if payment.Status == "succeeded" {
			return apierr.BadRequest("payment already processed")
		}
		if _, err := tx.MarkPaymentSucceeded(externalID); err != nil {
			return err
		}
		created, err = tx.CreateTransaction(domain.Transaction{
			UserID:     payment.UserID,
			Type:       domain.TransactionTypeDeposit,
			Amount:     payment.Amount,
			ExternalID: payment.ExternalID,
		})
		return err
	})
	if err != nil {
		return domain.Transaction{}, err
	}
	s.notifier.Publish(created.UserID, "balance_deposit", "Balance replenished")
	return created, nil
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
