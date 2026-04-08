package service

import (
	"testing"

	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/notifications"
	"github.com/soundmarket/backend/internal/repository"
)

func TestBidAllowedOnlyForRequest(t *testing.T) {
	store := repository.NewMemoryStore()
	notifier := notifications.NewInMemoryService()
	cardSvc := NewCardService(store, notifier)
	bidSvc := NewBidService(store, notifier)

	engineer := domain.User{ID: "eng-1", Role: domain.RoleEngineer}
	offer, err := cardSvc.Create(engineer, domain.Card{
		CardType: domain.CardTypeOffer,
		Kind:     domain.CardKindProduct,
		Title:    "Offer",
		Price:    100,
	})
	if err != nil {
		t.Fatalf("create offer: %v", err)
	}
	if _, err := bidSvc.Create(engineer, offer.ID, 90, "my bid"); err == nil {
		t.Fatal("expected bid creation on offer to fail")
	}
}
