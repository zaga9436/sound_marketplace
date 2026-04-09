package service

import (
	"errors"
	"strings"

	"github.com/soundmarket/backend/internal/apierr"
	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/notifications"
	"github.com/soundmarket/backend/internal/repository"
)

type ReviewService struct {
	store    repository.Store
	notifier notifications.Service
}

func NewReviewService(store repository.Store, notifier notifications.Service) *ReviewService {
	return &ReviewService{store: store, notifier: notifier}
}

func (s *ReviewService) Create(actor domain.User, orderID string, rating int, text string) (domain.Review, error) {
	if err := ensureActiveUser(s.store, actor); err != nil {
		return domain.Review{}, err
	}
	text = strings.TrimSpace(text)
	if rating < 1 || rating > 5 {
		return domain.Review{}, apierr.BadRequest("rating must be between 1 and 5")
	}
	if text == "" {
		return domain.Review{}, apierr.BadRequest("text is required")
	}

	var created domain.Review
	err := s.store.WithTx(func(tx repository.Store) error {
		order, err := tx.GetOrder(orderID)
		if err != nil {
			return apierr.NotFound("order not found")
		}
		if actor.ID != order.CustomerID && actor.ID != order.EngineerID && actor.Role != domain.RoleAdmin {
			return apierr.Forbidden("forbidden")
		}
		if order.Status != domain.OrderStatusCompleted {
			return apierr.BadRequest("review is available only for completed order")
		}
		if actor.ID != order.CustomerID {
			return apierr.Forbidden("only customer can leave review for completed order")
		}
		if _, err := tx.GetReviewByOrderAndAuthor(order.ID, actor.ID); err == nil {
			return apierr.Conflict("review for this order already exists")
		} else if !errors.Is(err, repository.ErrNotFound) {
			return err
		}

		created, err = tx.CreateReview(domain.Review{
			OrderID:      order.ID,
			AuthorID:     actor.ID,
			TargetUserID: order.EngineerID,
			Rating:       rating,
			Text:         text,
		})
		if err != nil {
			return err
		}
		if _, err := tx.RefreshProfileRating(order.EngineerID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return domain.Review{}, err
	}

	s.notifier.Publish(created.TargetUserID, "review_received", "New verified review received")
	return created, nil
}
