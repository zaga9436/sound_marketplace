package service

import (
	"testing"

	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/notifications"
	"github.com/soundmarket/backend/internal/repository"
)

func TestOfferRequiresEngineerRole(t *testing.T) {
	svc := NewCardService(repository.NewMemoryStore(), notifications.NewInMemoryService())
	_, err := svc.Create(domain.User{ID: "u1", Role: domain.RoleCustomer}, domain.Card{
		CardType: domain.CardTypeOffer,
		Kind:     domain.CardKindProduct,
		Title:    "Beat",
		Price:    1000,
	})
	if err == nil {
		t.Fatal("expected customer offer creation to fail")
	}
}
