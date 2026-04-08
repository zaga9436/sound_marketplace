package repository

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/soundmarket/backend/internal/domain"
)

var ErrNotFound = errors.New("not found")

type MemoryStore struct {
	mu           sync.RWMutex
	users        map[string]domain.User
	profiles     map[string]domain.Profile
	cards        map[string]domain.Card
	bids         map[string]domain.Bid
	orders       map[string]domain.Order
	transactions map[string]domain.Transaction
	payments     map[string]domain.Payment
	balances     map[string]int64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users:        map[string]domain.User{},
		profiles:     map[string]domain.Profile{},
		cards:        map[string]domain.Card{},
		bids:         map[string]domain.Bid{},
		orders:       map[string]domain.Order{},
		transactions: map[string]domain.Transaction{},
		payments:     map[string]domain.Payment{},
		balances:     map[string]int64{},
	}
}

func (s *MemoryStore) CreateUser(email, passwordHash string, role domain.Role) (domain.User, domain.Profile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, user := range s.users {
		if strings.EqualFold(user.Email, email) {
			return domain.User{}, domain.Profile{}, errors.New("email already exists")
		}
	}
	now := time.Now()
	user := domain.User{
		ID:           uuid.NewString(),
		Email:        strings.ToLower(email),
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    now,
	}
	profile := domain.Profile{
		UserID:      user.ID,
		DisplayName: strings.Split(user.Email, "@")[0],
		CreatedAt:   now,
	}
	s.users[user.ID] = user
	s.profiles[user.ID] = profile
	return user, profile, nil
}

func (s *MemoryStore) FindUserByEmail(email string) (domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, user := range s.users {
		if strings.EqualFold(user.Email, email) {
			return user, nil
		}
	}
	return domain.User{}, ErrNotFound
}

func (s *MemoryStore) GetUser(userID string) (domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[userID]
	if !ok {
		return domain.User{}, ErrNotFound
	}
	return user, nil
}

func (s *MemoryStore) GetProfile(userID string) (domain.Profile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	profile, ok := s.profiles[userID]
	if !ok {
		return domain.Profile{}, ErrNotFound
	}
	return profile, nil
}

func (s *MemoryStore) UpdateProfile(userID, displayName, bio string) (domain.Profile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	profile, ok := s.profiles[userID]
	if !ok {
		return domain.Profile{}, ErrNotFound
	}
	profile.DisplayName = displayName
	profile.Bio = bio
	s.profiles[userID] = profile
	return profile, nil
}

func (s *MemoryStore) CreateCard(card domain.Card) domain.Card {
	s.mu.Lock()
	defer s.mu.Unlock()
	card.ID = uuid.NewString()
	card.CreatedAt = time.Now()
	s.cards[card.ID] = card
	return card
}

func (s *MemoryStore) UpdateCard(cardID string, payload domain.Card) (domain.Card, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	card, ok := s.cards[cardID]
	if !ok {
		return domain.Card{}, ErrNotFound
	}
	card.Title = payload.Title
	card.Description = payload.Description
	card.Price = payload.Price
	card.Tags = payload.Tags
	card.Kind = payload.Kind
	card.IsPublished = payload.IsPublished
	s.cards[cardID] = card
	return card, nil
}

func (s *MemoryStore) ListCards(cardType, query string) []domain.Card {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]domain.Card, 0)
	for _, card := range s.cards {
		if cardType != "" && string(card.CardType) != cardType {
			continue
		}
		if query != "" && !strings.Contains(strings.ToLower(card.Title+" "+card.Description), strings.ToLower(query)) {
			continue
		}
		result = append(result, card)
	}
	return result
}

func (s *MemoryStore) GetCard(cardID string) (domain.Card, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	card, ok := s.cards[cardID]
	if !ok {
		return domain.Card{}, ErrNotFound
	}
	return card, nil
}

func (s *MemoryStore) CreateBid(bid domain.Bid) domain.Bid {
	s.mu.Lock()
	defer s.mu.Unlock()
	bid.ID = uuid.NewString()
	bid.CreatedAt = time.Now()
	s.bids[bid.ID] = bid
	return bid
}

func (s *MemoryStore) ListBidsByRequest(requestID string) []domain.Bid {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]domain.Bid, 0)
	for _, bid := range s.bids {
		if bid.RequestID == requestID {
			result = append(result, bid)
		}
	}
	return result
}

func (s *MemoryStore) GetBid(bidID string) (domain.Bid, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	bid, ok := s.bids[bidID]
	if !ok {
		return domain.Bid{}, ErrNotFound
	}
	return bid, nil
}

func (s *MemoryStore) CreateOrder(order domain.Order) domain.Order {
	s.mu.Lock()
	defer s.mu.Unlock()
	order.ID = uuid.NewString()
	order.CreatedAt = time.Now()
	order.LastStatusTime = order.CreatedAt
	s.orders[order.ID] = order
	return order
}

func (s *MemoryStore) GetOrder(orderID string) (domain.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ok := s.orders[orderID]
	if !ok {
		return domain.Order{}, ErrNotFound
	}
	return order, nil
}

func (s *MemoryStore) UpdateOrder(order domain.Order) domain.Order {
	s.mu.Lock()
	defer s.mu.Unlock()
	order.LastStatusTime = time.Now()
	s.orders[order.ID] = order
	return order
}

func (s *MemoryStore) CreateTransaction(tx domain.Transaction) domain.Transaction {
	s.mu.Lock()
	defer s.mu.Unlock()
	tx.ID = uuid.NewString()
	tx.CreatedAt = time.Now()
	s.transactions[tx.ID] = tx
	switch tx.Type {
	case domain.TransactionTypeDeposit, domain.TransactionTypeRelease, domain.TransactionTypeRefund, domain.TransactionTypePartialRefund:
		s.balances[tx.UserID] += tx.Amount
	case domain.TransactionTypeHold:
		s.balances[tx.UserID] -= tx.Amount
	}
	return tx
}

func (s *MemoryStore) GetBalance(userID string) int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.balances[userID]
}

func (s *MemoryStore) CreatePayment(payment domain.Payment) domain.Payment {
	s.mu.Lock()
	defer s.mu.Unlock()
	payment.ID = uuid.NewString()
	payment.CreatedAt = time.Now()
	s.payments[payment.ID] = payment
	return payment
}
