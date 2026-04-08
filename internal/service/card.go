package service

import (
	"fmt"
	"strings"

	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/notifications"
	"github.com/soundmarket/backend/internal/repository"
)

type CardService struct {
	store    repository.Store
	notifier notifications.Service
}

func NewCardService(store repository.Store, notifier notifications.Service) *CardService {
	return &CardService{store: store, notifier: notifier}
}

func (s *CardService) Create(actor domain.User, payload domain.Card) (domain.Card, error) {
	if payload.CardType != domain.CardTypeOffer && payload.CardType != domain.CardTypeRequest {
		return domain.Card{}, fmt.Errorf("card_type must be offer or request")
	}
	if payload.Kind != domain.CardKindProduct && payload.Kind != domain.CardKindService {
		return domain.Card{}, fmt.Errorf("kind must be product or service")
	}
	if payload.CardType == domain.CardTypeOffer && actor.Role != domain.RoleEngineer {
		return domain.Card{}, fmt.Errorf("only engineer can create offers")
	}
	if payload.CardType == domain.CardTypeRequest && actor.Role != domain.RoleCustomer {
		return domain.Card{}, fmt.Errorf("only customer can create requests")
	}
	if strings.TrimSpace(payload.Title) == "" {
		return domain.Card{}, fmt.Errorf("title is required")
	}
	if strings.TrimSpace(payload.Description) == "" {
		return domain.Card{}, fmt.Errorf("description is required")
	}
	if payload.Price < 0 {
		return domain.Card{}, fmt.Errorf("price must be non-negative")
	}

	payload.AuthorID = actor.ID
	payload.IsPublished = true
	card, err := s.store.CreateCard(payload)
	if err != nil {
		return domain.Card{}, err
	}
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
	payload.CardType = card.CardType
	payload.AuthorID = card.AuthorID
	return s.store.UpdateCard(cardID, payload)
}

func (s *CardService) List(cardType, query string) ([]domain.Card, error) {
	cards, err := s.store.ListCards(cardType, query)
	if err != nil {
		return nil, err
	}
	if cards == nil {
		return []domain.Card{}, nil
	}
	return cards, nil
}

func (s *CardService) Get(cardID string) (domain.Card, error) {
	return s.store.GetCard(cardID)
}
