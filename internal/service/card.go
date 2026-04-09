package service

import (
	"context"
	"strings"

	"github.com/soundmarket/backend/internal/apierr"
	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/notifications"
	"github.com/soundmarket/backend/internal/repository"
	"github.com/soundmarket/backend/internal/storage"
)

type CardService struct {
	store    repository.Store
	notifier notifications.Service
	storage  storage.Adapter
}

func NewCardService(store repository.Store, notifier notifications.Service, storageAdapter storage.Adapter) *CardService {
	return &CardService{store: store, notifier: notifier, storage: storageAdapter}
}

func (s *CardService) Create(actor domain.User, payload domain.Card) (domain.Card, error) {
	if payload.CardType != domain.CardTypeOffer && payload.CardType != domain.CardTypeRequest {
		return domain.Card{}, apierr.BadRequest("card_type must be offer or request")
	}
	if payload.Kind != domain.CardKindProduct && payload.Kind != domain.CardKindService {
		return domain.Card{}, apierr.BadRequest("kind must be product or service")
	}
	if payload.CardType == domain.CardTypeOffer && actor.Role != domain.RoleEngineer {
		return domain.Card{}, apierr.Forbidden("only engineer can create offers")
	}
	if payload.CardType == domain.CardTypeRequest && actor.Role != domain.RoleCustomer {
		return domain.Card{}, apierr.Forbidden("only customer can create requests")
	}
	if strings.TrimSpace(payload.Title) == "" {
		return domain.Card{}, apierr.BadRequest("title is required")
	}
	if strings.TrimSpace(payload.Description) == "" {
		return domain.Card{}, apierr.BadRequest("description is required")
	}
	if payload.Price < 0 {
		return domain.Card{}, apierr.BadRequest("price must be non-negative")
	}

	payload.AuthorID = actor.ID
	payload.IsPublished = true
	card, err := s.store.CreateCard(payload)
	if err != nil {
		return domain.Card{}, err
	}
	card.PreviewURLs = []string{}
	s.notifier.Publish(actor.ID, "card_published", "Card published")
	return card, nil
}

func (s *CardService) Update(actor domain.User, cardID string, payload domain.Card) (domain.Card, error) {
	card, err := s.store.GetCard(cardID)
	if err != nil {
		return domain.Card{}, apierr.NotFound("card not found")
	}
	if actor.Role != domain.RoleAdmin && card.AuthorID != actor.ID {
		return domain.Card{}, apierr.Forbidden("forbidden")
	}
	payload.CardType = card.CardType
	payload.AuthorID = card.AuthorID
	updated, err := s.store.UpdateCard(cardID, payload)
	if err != nil {
		return domain.Card{}, err
	}
	updatedCards := []domain.Card{updated}
	if err := s.attachPreviewURLs(context.Background(), updatedCards); err == nil {
		updated = updatedCards[0]
	}
	if updated.PreviewURLs == nil {
		updated.PreviewURLs = []string{}
	}
	return updated, nil
}

func (s *CardService) List(cardType, query string) ([]domain.Card, error) {
	cards, err := s.store.ListCards(cardType, query)
	if err != nil {
		return nil, err
	}
	if err := s.attachPreviewURLs(context.Background(), cards); err != nil {
		return nil, err
	}
	if cards == nil {
		return []domain.Card{}, nil
	}
	return cards, nil
}

func (s *CardService) Get(cardID string) (domain.Card, error) {
	card, err := s.store.GetCard(cardID)
	if err != nil {
		return domain.Card{}, apierr.NotFound("card not found")
	}
	cards := []domain.Card{card}
	if err := s.attachPreviewURLs(context.Background(), cards); err != nil {
		return domain.Card{}, err
	}
	return cards[0], nil
}

func (s *CardService) attachPreviewURLs(_ context.Context, cards []domain.Card) error {
	for i := range cards {
		mediaFiles, err := s.store.ListMediaByCardAndRole(cards[i].ID, domain.MediaRolePreview)
		if err != nil {
			return err
		}
		cards[i].PreviewURLs = make([]string, 0, len(mediaFiles))
		for _, media := range mediaFiles {
			cards[i].PreviewURLs = append(cards[i].PreviewURLs, s.storage.PublicURL(media.FileKey))
		}
	}
	return nil
}
