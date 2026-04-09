package service

import (
	"errors"
	"strings"

	"github.com/soundmarket/backend/internal/apierr"
	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/notifications"
	"github.com/soundmarket/backend/internal/repository"
)

type BidService struct {
	store    repository.Store
	notifier notifications.Service
}

func NewBidService(store repository.Store, notifier notifications.Service) *BidService {
	return &BidService{store: store, notifier: notifier}
}

func (s *BidService) Create(actor domain.User, requestID string, price int64, message string) (domain.Bid, error) {
	if err := ensureActiveUser(s.store, actor); err != nil {
		return domain.Bid{}, err
	}
	if actor.Role != domain.RoleEngineer {
		return domain.Bid{}, apierr.Forbidden("only engineer can submit bids")
	}
	if price <= 0 {
		return domain.Bid{}, apierr.BadRequest("price must be positive")
	}
	if strings.TrimSpace(message) == "" {
		return domain.Bid{}, apierr.BadRequest("message is required")
	}
	card, err := s.store.GetCard(requestID)
	if err != nil {
		return domain.Bid{}, apierr.NotFound("request not found")
	}
	if card.CardType != domain.CardTypeRequest {
		return domain.Bid{}, apierr.BadRequest("bids are allowed only for requests")
	}
	if card.AuthorID == actor.ID {
		return domain.Bid{}, apierr.Forbidden("request author cannot bid on own request")
	}
	if _, err := s.store.GetBidByRequestAndEngineer(requestID, actor.ID); err == nil {
		return domain.Bid{}, apierr.BadRequest("bid for this request already exists")
	} else if !errors.Is(err, repository.ErrNotFound) {
		return domain.Bid{}, err
	}
	bid, err := s.store.CreateBid(domain.Bid{
		RequestID:  requestID,
		EngineerID: actor.ID,
		Price:      price,
		Message:    strings.TrimSpace(message),
	})
	if err != nil {
		return domain.Bid{}, err
	}
	s.notifier.Publish(card.AuthorID, "new_bid", "New bid received")
	return bid, nil
}

func (s *BidService) List(actor domain.User, requestID string) ([]domain.Bid, error) {
	card, err := s.store.GetCard(requestID)
	if err != nil {
		return nil, apierr.NotFound("request not found")
	}
	if card.CardType != domain.CardTypeRequest {
		return nil, apierr.BadRequest("bids are allowed only for requests")
	}
	if actor.Role != domain.RoleAdmin && actor.ID != card.AuthorID {
		return nil, apierr.Forbidden("forbidden")
	}

	if actor.Role == domain.RoleAdmin {
		bids, err := s.store.ListBidsByRequest(requestID)
		if err != nil {
			return nil, err
		}
		if bids == nil {
			return []domain.Bid{}, nil
		}
		return bids, nil
	}

	bids, err := s.store.ListBidsByRequestForAuthor(requestID, actor.ID)
	if err != nil {
		return nil, err
	}
	if bids == nil {
		return []domain.Bid{}, nil
	}
	return bids, nil
}
