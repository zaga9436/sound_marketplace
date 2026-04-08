package service

import (
	"fmt"

	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/notifications"
	"github.com/soundmarket/backend/internal/repository"
)

type OrderService struct {
	store    *repository.MemoryStore
	notifier notifications.Service
}

func NewOrderService(store *repository.MemoryStore, notifier notifications.Service) *OrderService {
	return &OrderService{store: store, notifier: notifier}
}

func (s *OrderService) CreateFromOffer(customer domain.User, cardID string) (domain.Order, error) {
	if customer.Role != domain.RoleCustomer {
		return domain.Order{}, fmt.Errorf("only customer can create order")
	}
	card, err := s.store.GetCard(cardID)
	if err != nil {
		return domain.Order{}, err
	}
	if card.CardType != domain.CardTypeOffer {
		return domain.Order{}, fmt.Errorf("card is not an offer")
	}
	if s.store.GetBalance(customer.ID) < card.Price {
		return domain.Order{}, fmt.Errorf("insufficient balance")
	}
	order := s.store.CreateOrder(domain.Order{
		CardID:      card.ID,
		CustomerID:  customer.ID,
		EngineerID:  card.AuthorID,
		Amount:      card.Price,
		Status:      domain.OrderStatusOnHold,
	})
	s.store.CreateTransaction(domain.Transaction{
		UserID:  customer.ID,
		OrderID: order.ID,
		Type:    domain.TransactionTypeHold,
		Amount:  card.Price,
	})
	s.notifier.Publish(card.AuthorID, "order_created", "New order created")
	return order, nil
}

func (s *OrderService) CreateFromBid(customer domain.User, bidID string) (domain.Order, error) {
	if customer.Role != domain.RoleCustomer {
		return domain.Order{}, fmt.Errorf("only customer can accept bid")
	}
	bid, err := s.store.GetBid(bidID)
	if err != nil {
		return domain.Order{}, err
	}
	requestCard, err := s.store.GetCard(bid.RequestID)
	if err != nil {
		return domain.Order{}, err
	}
	if requestCard.AuthorID != customer.ID {
		return domain.Order{}, fmt.Errorf("forbidden")
	}
	if s.store.GetBalance(customer.ID) < bid.Price {
		return domain.Order{}, fmt.Errorf("insufficient balance")
	}
	order := s.store.CreateOrder(domain.Order{
		RequestID:   requestCard.ID,
		BidID:       bid.ID,
		CustomerID:  customer.ID,
		EngineerID:  bid.EngineerID,
		Amount:      bid.Price,
		Status:      domain.OrderStatusOnHold,
	})
	s.store.CreateTransaction(domain.Transaction{
		UserID:  customer.ID,
		OrderID: order.ID,
		Type:    domain.TransactionTypeHold,
		Amount:  bid.Price,
	})
	s.notifier.Publish(bid.EngineerID, "bid_selected", "Your bid was selected")
	return order, nil
}

func (s *OrderService) Get(orderID string, actor domain.User) (domain.Order, error) {
	order, err := s.store.GetOrder(orderID)
	if err != nil {
		return domain.Order{}, err
	}
	if actor.Role != domain.RoleAdmin && actor.ID != order.CustomerID && actor.ID != order.EngineerID {
		return domain.Order{}, fmt.Errorf("forbidden")
	}
	return order, nil
}

func (s *OrderService) UpdateStatus(actor domain.User, orderID string, next domain.OrderStatus) (domain.Order, error) {
	order, err := s.Get(orderID, actor)
	if err != nil {
		return domain.Order{}, err
	}
	if !isStatusTransitionAllowed(order.Status, next) {
		return domain.Order{}, fmt.Errorf("invalid status transition")
	}
	if next == domain.OrderStatusCompleted {
		s.store.CreateTransaction(domain.Transaction{
			UserID:  order.EngineerID,
			OrderID: order.ID,
			Type:    domain.TransactionTypeRelease,
			Amount:  order.Amount,
		})
	}
	order.Status = next
	updated := s.store.UpdateOrder(order)
	s.notifier.Publish(order.CustomerID, "order_status_changed", string(next))
	s.notifier.Publish(order.EngineerID, "order_status_changed", string(next))
	return updated, nil
}

func isStatusTransitionAllowed(current, next domain.OrderStatus) bool {
	allowed := map[domain.OrderStatus][]domain.OrderStatus{
		domain.OrderStatusCreated:    {domain.OrderStatusOnHold, domain.OrderStatusCancelled},
		domain.OrderStatusOnHold:     {domain.OrderStatusInProgress, domain.OrderStatusCancelled, domain.OrderStatusDispute},
		domain.OrderStatusInProgress: {domain.OrderStatusReview, domain.OrderStatusDispute},
		domain.OrderStatusReview:     {domain.OrderStatusCompleted, domain.OrderStatusInProgress, domain.OrderStatusDispute},
		domain.OrderStatusDispute:    {domain.OrderStatusCancelled, domain.OrderStatusCompleted},
	}
	for _, candidate := range allowed[current] {
		if candidate == next {
			return true
		}
	}
	return false
}
