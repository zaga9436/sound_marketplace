package service

import (
	"errors"

	"github.com/soundmarket/backend/internal/apierr"
	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/notifications"
	"github.com/soundmarket/backend/internal/repository"
)

type OrderService struct {
	store    repository.Store
	notifier notifications.Service
}

func NewOrderService(store repository.Store, notifier notifications.Service) *OrderService {
	return &OrderService{store: store, notifier: notifier}
}

func (s *OrderService) CreateFromOffer(customer domain.User, cardID string) (domain.Order, error) {
	if err := ensureActiveUser(s.store, customer); err != nil {
		return domain.Order{}, err
	}
	if customer.Role != domain.RoleCustomer {
		return domain.Order{}, apierr.Forbidden("only customer can create order")
	}

	var created domain.Order
	err := s.store.WithTx(func(tx repository.Store) error {
		card, err := tx.GetCard(cardID)
		if err != nil {
			return apierr.NotFound("offer not found")
		}
		if card.CardType != domain.CardTypeOffer {
			return apierr.BadRequest("card is not an offer")
		}
		if _, err := tx.GetOrderByCardAndCustomer(cardID, customer.ID); err == nil {
			return apierr.BadRequest("order for this offer already exists")
		} else if !errors.Is(err, repository.ErrNotFound) {
			return err
		}
		balance, err := tx.GetBalance(customer.ID)
		if err != nil {
			return err
		}
		if balance < card.Price {
			return apierr.BadRequest("insufficient balance")
		}

		created, err = tx.CreateOrder(domain.Order{
			CardID:     card.ID,
			CustomerID: customer.ID,
			EngineerID: card.AuthorID,
			Amount:     card.Price,
			Status:     domain.OrderStatusOnHold,
		})
		if err != nil {
			return err
		}

		_, err = tx.CreateTransaction(domain.Transaction{
			UserID:  customer.ID,
			OrderID: created.ID,
			Type:    domain.TransactionTypeHold,
			Amount:  card.Price,
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return domain.Order{}, err
	}

	s.notifier.Publish(created.EngineerID, "order_created", "New order created from offer")
	return created, nil
}

func (s *OrderService) CreateFromBid(customer domain.User, bidID string) (domain.Order, error) {
	if err := ensureActiveUser(s.store, customer); err != nil {
		return domain.Order{}, err
	}
	if customer.Role != domain.RoleCustomer {
		return domain.Order{}, apierr.Forbidden("only customer can accept bid")
	}

	var created domain.Order
	err := s.store.WithTx(func(tx repository.Store) error {
		bid, err := tx.GetBid(bidID)
		if err != nil {
			return apierr.NotFound("bid not found")
		}
		requestCard, err := tx.GetCard(bid.RequestID)
		if err != nil {
			return apierr.NotFound("request not found")
		}
		if requestCard.CardType != domain.CardTypeRequest {
			return apierr.BadRequest("bid does not belong to request")
		}
		if requestCard.AuthorID != customer.ID {
			return apierr.Forbidden("forbidden")
		}
		if _, err := tx.GetOrderByBidID(bid.ID); err == nil {
			return apierr.BadRequest("order for this bid already exists")
		} else if !errors.Is(err, repository.ErrNotFound) {
			return err
		}
		balance, err := tx.GetBalance(customer.ID)
		if err != nil {
			return err
		}
		if balance < bid.Price {
			return apierr.BadRequest("insufficient balance")
		}

		created, err = tx.CreateOrder(domain.Order{
			RequestID:  requestCard.ID,
			BidID:      bid.ID,
			CustomerID: customer.ID,
			EngineerID: bid.EngineerID,
			Amount:     bid.Price,
			Status:     domain.OrderStatusOnHold,
		})
		if err != nil {
			return err
		}

		_, err = tx.CreateTransaction(domain.Transaction{
			UserID:  customer.ID,
			OrderID: created.ID,
			Type:    domain.TransactionTypeHold,
			Amount:  bid.Price,
		})
		return err
	})
	if err != nil {
		return domain.Order{}, err
	}

	s.notifier.Publish(created.EngineerID, "order_created", "Your bid was selected and order created")
	return created, nil
}

func (s *OrderService) Get(orderID string, actor domain.User) (domain.Order, error) {
	order, err := s.store.GetOrder(orderID)
	if err != nil {
		return domain.Order{}, apierr.NotFound("order not found")
	}
	if actor.Role != domain.RoleAdmin && actor.ID != order.CustomerID && actor.ID != order.EngineerID {
		return domain.Order{}, apierr.Forbidden("forbidden")
	}
	return order, nil
}

func (s *OrderService) List(actor domain.User) ([]domain.Order, error) {
	var (
		orders []domain.Order
		err    error
	)

	switch actor.Role {
	case domain.RoleCustomer:
		orders, err = s.store.ListOrdersByCustomer(actor.ID)
	case domain.RoleEngineer:
		orders, err = s.store.ListOrdersByEngineer(actor.ID)
	case domain.RoleAdmin:
		orders, err = s.store.ListOrders()
	default:
		return nil, apierr.Forbidden("forbidden")
	}
	if err != nil {
		return nil, err
	}
	if orders == nil {
		return []domain.Order{}, nil
	}
	return orders, nil
}

func (s *OrderService) UpdateStatus(actor domain.User, orderID string, next domain.OrderStatus) (domain.Order, error) {
	var updated domain.Order
	err := s.store.WithTx(func(tx repository.Store) error {
		order, err := tx.GetOrder(orderID)
		if err != nil {
			return apierr.NotFound("order not found")
		}
		if actor.Role != domain.RoleAdmin && actor.ID != order.CustomerID && actor.ID != order.EngineerID {
			return apierr.Forbidden("forbidden")
		}
		if next == domain.OrderStatusDispute {
			return apierr.BadRequest("use dispute endpoint to open dispute")
		}
		if order.Status == domain.OrderStatusDispute {
			return apierr.BadRequest("use dispute close endpoint to resolve disputed order")
		}
		if !isStatusTransitionAllowed(order.Status, next) {
			return apierr.BadRequest("invalid status transition")
		}
		if err := validateStatusActor(actor, order, next); err != nil {
			return err
		}

		switch next {
		case domain.OrderStatusCompleted:
			if _, err := tx.CreateTransaction(domain.Transaction{
				UserID:  order.EngineerID,
				OrderID: order.ID,
				Type:    domain.TransactionTypeRelease,
				Amount:  order.Amount,
			}); err != nil {
				return err
			}
		case domain.OrderStatusCancelled:
			if _, err := tx.CreateTransaction(domain.Transaction{
				UserID:  order.CustomerID,
				OrderID: order.ID,
				Type:    domain.TransactionTypeRefund,
				Amount:  order.Amount,
			}); err != nil {
				return err
			}
		}

		order.Status = next
		updated, err = tx.UpdateOrder(order)
		return err
	})
	if err != nil {
		return domain.Order{}, err
	}

	s.notifier.Publish(updated.CustomerID, "order_status_changed", string(next))
	s.notifier.Publish(updated.EngineerID, "order_status_changed", string(next))
	return updated, nil
}

func validateStatusActor(actor domain.User, order domain.Order, next domain.OrderStatus) error {
	if actor.Role == domain.RoleAdmin {
		return nil
	}

	switch next {
	case domain.OrderStatusInProgress, domain.OrderStatusReview:
		if actor.ID != order.EngineerID {
			return apierr.Forbidden("only engineer can perform this status transition")
		}
	case domain.OrderStatusCompleted, domain.OrderStatusCancelled:
		if actor.ID != order.CustomerID {
			return apierr.Forbidden("only customer can perform this status transition")
		}
	case domain.OrderStatusDispute:
		if actor.ID != order.CustomerID && actor.ID != order.EngineerID {
			return apierr.Forbidden("only order participants can open dispute")
		}
	}

	return nil
}

func isStatusTransitionAllowed(current, next domain.OrderStatus) bool {
	allowed := map[domain.OrderStatus][]domain.OrderStatus{
		domain.OrderStatusCreated:    {domain.OrderStatusOnHold, domain.OrderStatusCancelled},
		domain.OrderStatusOnHold:     {domain.OrderStatusInProgress, domain.OrderStatusCancelled, domain.OrderStatusDispute},
		domain.OrderStatusInProgress: {domain.OrderStatusReview, domain.OrderStatusDispute, domain.OrderStatusCancelled},
		domain.OrderStatusReview:     {domain.OrderStatusCompleted, domain.OrderStatusInProgress, domain.OrderStatusDispute, domain.OrderStatusCancelled},
		domain.OrderStatusDispute:    {domain.OrderStatusCancelled, domain.OrderStatusCompleted},
	}
	for _, candidate := range allowed[current] {
		if candidate == next {
			return true
		}
	}
	return false
}
