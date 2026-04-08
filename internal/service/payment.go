package service

import (
	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/notifications"
	"github.com/soundmarket/backend/internal/payments"
	"github.com/soundmarket/backend/internal/repository"
)

type PaymentService struct {
	store    *repository.MemoryStore
	adapter  payments.Adapter
	notifier notifications.Service
}

func NewPaymentService(store *repository.MemoryStore, adapter payments.Adapter, notifier notifications.Service) *PaymentService {
	return &PaymentService{store: store, adapter: adapter, notifier: notifier}
}

func (s *PaymentService) CreateDeposit(user domain.User, amount int64) (domain.Payment, error) {
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
	}), nil
}

func (s *PaymentService) ProcessWebhook(userID, externalID string, amount int64) domain.Transaction {
	tx := s.store.CreateTransaction(domain.Transaction{
		UserID:     userID,
		Type:       domain.TransactionTypeDeposit,
		Amount:     amount,
		ExternalID: externalID,
	})
	s.notifier.Publish(userID, "balance_deposit", "Balance replenished")
	return tx
}

func (s *PaymentService) Balance(userID string) int64 {
	return s.store.GetBalance(userID)
}
