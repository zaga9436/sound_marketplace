package notifications

import (
	"context"
	"time"

	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/realtime"
	"github.com/soundmarket/backend/internal/repository"
)

type Service interface {
	Publish(userID, eventType, message string)
}

type RepositoryBackedService struct {
	store  repository.Store
	broker *realtime.Broker
}

func NewRepositoryBackedService(store repository.Store, broker *realtime.Broker) *RepositoryBackedService {
	return &RepositoryBackedService{store: store, broker: broker}
}

func (s *RepositoryBackedService) Publish(userID, eventType, message string) {
	notification, err := s.store.CreateNotification(userID, eventType, message)
	if err != nil {
		return
	}
	ctx := context.Background()
	_, _ = s.broker.IncrementCounter(ctx, realtime.NotificationsUnreadKey(userID), 1)
	_ = s.broker.PublishJSON(ctx, realtime.NotificationsChannel(userID), map[string]interface{}{
		"type":         "notification",
		"notification": notification,
	})
}

func AsMap(notification domain.Notification) map[string]interface{} {
	return map[string]interface{}{
		"id":         notification.ID,
		"user_id":    notification.UserID,
		"type":       notification.Type,
		"message":    notification.Message,
		"is_read":    notification.IsRead,
		"created_at": notification.CreatedAt.Format(time.RFC3339),
	}
}
