package service

import (
	"context"
	"time"

	"github.com/soundmarket/backend/internal/apierr"
	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/realtime"
	"github.com/soundmarket/backend/internal/repository"
)

type NotificationListResult struct {
	Items       []domain.Notification `json:"items"`
	UnreadCount int64                 `json:"unread_count"`
}

type NotificationService struct {
	store  repository.Store
	broker *realtime.Broker
}

func NewNotificationService(store repository.Store, broker *realtime.Broker) *NotificationService {
	return &NotificationService{store: store, broker: broker}
}

func (s *NotificationService) List(ctx context.Context, actor domain.User, limit int, beforeID string) (NotificationListResult, error) {
	items, err := s.store.ListNotifications(actor.ID, limit, beforeID)
	if err != nil {
		return NotificationListResult{}, err
	}
	unreadCount, err := s.broker.GetCounter(ctx, realtime.NotificationsUnreadKey(actor.ID))
	if err != nil || unreadCount == 0 {
		unreadCount, _ = s.store.CountUnreadNotifications(actor.ID)
		_ = s.broker.SetCounter(ctx, realtime.NotificationsUnreadKey(actor.ID), unreadCount, 24*time.Hour)
	}
	if items == nil {
		items = []domain.Notification{}
	}
	return NotificationListResult{
		Items:       items,
		UnreadCount: unreadCount,
	}, nil
}

func (s *NotificationService) MarkRead(ctx context.Context, actor domain.User, ids []string) (NotificationListResult, error) {
	if err := s.store.MarkNotificationsRead(actor.ID, ids); err != nil {
		return NotificationListResult{}, err
	}
	unreadCount, err := s.store.CountUnreadNotifications(actor.ID)
	if err != nil {
		return NotificationListResult{}, err
	}
	_ = s.broker.SetCounter(ctx, realtime.NotificationsUnreadKey(actor.ID), unreadCount, 24*time.Hour)
	items, err := s.store.ListNotifications(actor.ID, 20, "")
	if err != nil {
		return NotificationListResult{}, err
	}
	return NotificationListResult{
		Items:       items,
		UnreadCount: unreadCount,
	}, nil
}

func (s *NotificationService) PublishRealtimeSnapshot(ctx context.Context, actor domain.User) error {
	result, err := s.List(ctx, actor, 20, "")
	if err != nil {
		return apierr.BadRequest("unable to load notifications")
	}
	return s.broker.PublishJSON(ctx, realtime.NotificationsChannel(actor.ID), map[string]interface{}{
		"type":          "notifications_snapshot",
		"items":         result.Items,
		"unread_count":  result.UnreadCount,
	})
}
