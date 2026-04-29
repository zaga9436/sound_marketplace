package service

import (
	"errors"
	"strings"

	"github.com/soundmarket/backend/internal/apierr"
	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/notifications"
	"github.com/soundmarket/backend/internal/repository"
)

type DisputeService struct {
	store    repository.Store
	notifier notifications.Service
}

func NewDisputeService(store repository.Store, notifier notifications.Service) *DisputeService {
	return &DisputeService{store: store, notifier: notifier}
}

func (s *DisputeService) Open(actor domain.User, orderID, reason string) (domain.Dispute, error) {
	if err := ensureActiveUser(s.store, actor); err != nil {
		return domain.Dispute{}, err
	}
	if actor.Role == domain.RoleAdmin {
		return domain.Dispute{}, apierr.Forbidden("admin cannot open dispute as transaction side")
	}
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return domain.Dispute{}, apierr.BadRequest("reason is required")
	}

	var created domain.Dispute
	err := s.store.WithTx(func(tx repository.Store) error {
		order, err := tx.GetOrder(orderID)
		if err != nil {
			return apierr.NotFound("order not found")
		}
		if actor.ID != order.CustomerID && actor.ID != order.EngineerID {
			return apierr.Forbidden("forbidden")
		}
		if order.Status == domain.OrderStatusCompleted || order.Status == domain.OrderStatusCancelled {
			return apierr.BadRequest("cannot open dispute for closed order")
		}
		if _, err := tx.GetOpenDisputeByOrderID(orderID); err == nil {
			return apierr.Conflict("open dispute already exists for this order")
		} else if !errors.Is(err, repository.ErrNotFound) {
			return err
		}

		order.Status = domain.OrderStatusDispute
		order.DisputeReason = reason
		if _, err := tx.UpdateOrder(order); err != nil {
			return err
		}

		created, err = tx.CreateDispute(domain.Dispute{
			OrderID:        orderID,
			OpenedByUserID: actor.ID,
			Reason:         reason,
			Status:         domain.DisputeStatusOpen,
		})
		return err
	})
	if err != nil {
		return domain.Dispute{}, err
	}

	order, orderErr := s.store.GetOrder(orderID)
	if orderErr == nil {
		s.notifier.Publish(order.CustomerID, "dispute_opened", "Dispute opened")
		if order.EngineerID != order.CustomerID {
			s.notifier.Publish(order.EngineerID, "dispute_opened", "Dispute opened")
		}
	}
	return created, nil
}

func (s *DisputeService) Get(actor domain.User, orderID string) (domain.Dispute, error) {
	order, err := s.store.GetOrder(orderID)
	if err != nil {
		return domain.Dispute{}, apierr.NotFound("order not found")
	}
	if actor.Role != domain.RoleAdmin && actor.ID != order.CustomerID && actor.ID != order.EngineerID {
		return domain.Dispute{}, apierr.Forbidden("forbidden")
	}
	dispute, err := s.store.GetDisputeByOrderID(orderID)
	if err != nil {
		return domain.Dispute{}, apierr.NotFound("dispute not found")
	}
	return dispute, nil
}

func (s *DisputeService) Close(actor domain.User, orderID string, resolution domain.DisputeResolution) (domain.Dispute, error) {
	if resolution != domain.DisputeResolutionCompleteOrder && resolution != domain.DisputeResolutionCancelOrder {
		return domain.Dispute{}, apierr.BadRequest("resolution must be complete_order or cancel_order")
	}

	var closed domain.Dispute
	err := s.store.WithTx(func(tx repository.Store) error {
		order, err := tx.GetOrder(orderID)
		if err != nil {
			return apierr.NotFound("order not found")
		}
		dispute, err := tx.GetOpenDisputeByOrderID(orderID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				if latest, latestErr := tx.GetDisputeByOrderID(orderID); latestErr == nil && latest.Status == domain.DisputeStatusClosed {
					return apierr.Conflict("dispute already closed")
				}
				return apierr.NotFound("open dispute not found")
			}
			return err
		}
		if actor.Role != domain.RoleAdmin && actor.ID != dispute.OpenedByUserID {
			return apierr.Forbidden("forbidden")
		}

		switch resolution {
		case domain.DisputeResolutionCompleteOrder:
			order.Status = domain.OrderStatusCompleted
			if _, err := tx.CreateTransaction(domain.Transaction{
				UserID:  order.EngineerID,
				OrderID: order.ID,
				Type:    domain.TransactionTypeRelease,
				Amount:  engineerPayoutAmount(order.Amount),
			}); err != nil {
				return err
			}
		case domain.DisputeResolutionCancelOrder:
			order.Status = domain.OrderStatusCancelled
			if _, err := tx.CreateTransaction(domain.Transaction{
				UserID:  order.CustomerID,
				OrderID: order.ID,
				Type:    domain.TransactionTypeRefund,
				Amount:  order.Amount,
			}); err != nil {
				return err
			}
		}
		order.DisputeReason = ""
		if _, err := tx.UpdateOrder(order); err != nil {
			return err
		}

		closed, err = tx.CloseDispute(dispute.ID, resolution)
		return err
	})
	if err != nil {
		return domain.Dispute{}, err
	}

	order, orderErr := s.store.GetOrder(orderID)
	if orderErr == nil {
		s.notifier.Publish(order.CustomerID, "dispute_closed", "Dispute closed")
		if order.EngineerID != order.CustomerID {
			s.notifier.Publish(order.EngineerID, "dispute_closed", "Dispute closed")
		}
	}
	return closed, nil
}
