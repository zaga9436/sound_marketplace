package service

import (
	"context"

	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/realtime"
)

type RealtimeService struct {
	broker        *realtime.Broker
	chat          *ChatService
	notifications *NotificationService
}

func NewRealtimeService(broker *realtime.Broker, chat *ChatService, notifications *NotificationService) *RealtimeService {
	return &RealtimeService{
		broker:        broker,
		chat:          chat,
		notifications: notifications,
	}
}

func (s *RealtimeService) SubscribeChat(ctx context.Context, actor domain.User, orderID string) error {
	_, err := s.chat.authorizeOrderAccess(actor, orderID)
	return err
}

func (s *RealtimeService) ChatBootstrap(ctx context.Context, actor domain.User, orderID string) (map[string]interface{}, error) {
	if err := s.SubscribeChat(ctx, actor, orderID); err != nil {
		return nil, err
	}
	messages, err := s.chat.ListMessages(actor, orderID, "", 20)
	if err != nil {
		return nil, err
	}
	conversations, err := s.chat.ListConversations(ctx, actor, 100)
	if err != nil {
		return nil, err
	}
	unreadCount := int64(0)
	for _, conversation := range conversations {
		if conversation.OrderID == orderID {
			unreadCount = conversation.UnreadCount
			break
		}
	}
	return map[string]interface{}{
		"type":         "chat_bootstrap",
		"order_id":     orderID,
		"channel":      realtime.ChatChannel(orderID),
		"messages":     messages,
		"unread_count": unreadCount,
	}, nil
}

func (s *RealtimeService) NotificationsBootstrap(ctx context.Context, actor domain.User) (map[string]interface{}, error) {
	result, err := s.notifications.List(ctx, actor, 20, "")
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"type":         "notifications_bootstrap",
		"channel":      realtime.NotificationsChannel(actor.ID),
		"items":        result.Items,
		"unread_count": result.UnreadCount,
	}, nil
}

func (s *RealtimeService) Broker() *realtime.Broker {
	return s.broker
}
