package service

import (
	"fmt"

	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/notifications"
	"github.com/soundmarket/backend/internal/repository"
)

type CardService struct {
	store    *repository.MemoryStore
	notifier notifications.Service
}

func NewCardService(store *repository.MemoryStore, notifier notifications.Service) *CardService {
	return &CardService{store: store, notifier: notifier}
}

func (s *CardService) Create(actor domain.User, payload domain.Card) (domain.Card, error) {
	if payload.CardType == domain.CardTypeOffer && actor.Role != domain.RoleEngineer {
		return domain.Card{}, fmt.Errorf("only engineer can create offers")
	}
	if payload.CardType == domain.CardTypeRequest && actor.Role != domain.RoleCustomer {
		return domain.Card{}, fmt.Errorf("only customer can create requests")
	}
	payload.AuthorID = actor.ID
	payload.IsPublished = true
	card := s.store.CreateCard(payload)
	s.notifier.Publish(actor.ID, "card_published", "Card published")
	return card, nil
}

func (s *CardService) Update(actor domain.User, cardID string, payload domain.Card) (domain.Card, error) {
	card, err := s.store.GetCard(cardID)
	if err != nil {
		return domain.Card{}, err
	}
	if actor.Role != domain.RoleAdmin && card.AuthorID != actor.ID {
		return domain.Card{}, fmt.Errorf("forbidden")
	}
	return s.store.UpdateCard(cardID, payload)
}

func (s *CardService) List(cardType, query string) []domain.Card {
	return s.store.ListCards(cardType, query)
}

func (s *CardService) Get(cardID string) (domain.Card, error) {
	return s.store.GetCard(cardID)
}
