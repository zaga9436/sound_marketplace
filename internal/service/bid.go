package service

import (
	"fmt"

	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/notifications"
	"github.com/soundmarket/backend/internal/repository"
)

type BidService struct {
	store    *repository.MemoryStore
	notifier notifications.Service
}

func NewBidService(store *repository.MemoryStore, notifier notifications.Service) *BidService {
	return &BidService{store: store, notifier: notifier}
}

func (s *BidService) Create(actor domain.User, requestID string, price int64, message string) (domain.Bid, error) {
	if actor.Role != domain.RoleEngineer {
		return domain.Bid{}, fmt.Errorf("only engineer can submit bids")
	}
	card, err := s.store.GetCard(requestID)
	if err != nil {
		return domain.Bid{}, err
	}
	if card.CardType != domain.CardTypeRequest {
		return domain.Bid{}, fmt.Errorf("bids are allowed only for requests")
	}
	bid := s.store.CreateBid(domain.Bid{
		RequestID:  requestID,
		EngineerID: actor.ID,
		Price:      price,
		Message:    message,
	})
	s.notifier.Publish(card.AuthorID, "new_bid", "New bid received")
	return bid, nil
}

func (s *BidService) List(requestID string) []domain.Bid {
	return s.store.ListBidsByRequest(requestID)
}
