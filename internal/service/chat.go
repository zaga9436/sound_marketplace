package service

import (
	"context"
	"strings"
	"time"

	"github.com/soundmarket/backend/internal/apierr"
	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/notifications"
	"github.com/soundmarket/backend/internal/realtime"
	"github.com/soundmarket/backend/internal/repository"
)

type ChatService struct {
	store    repository.Store
	broker   *realtime.Broker
	notifier notifications.Service
}

func NewChatService(store repository.Store, broker *realtime.Broker, notifier notifications.Service) *ChatService {
	return &ChatService{store: store, broker: broker, notifier: notifier}
}

func (s *ChatService) ListConversations(ctx context.Context, actor domain.User, limit int) ([]domain.Conversation, error) {
	var (
		conversations []domain.Conversation
		err           error
	)

	switch actor.Role {
	case domain.RoleCustomer:
		conversations, err = s.store.ListConversationsByCustomer(actor.ID, limit)
	case domain.RoleEngineer:
		conversations, err = s.store.ListConversationsByEngineer(actor.ID, limit)
	case domain.RoleAdmin:
		conversations, err = s.store.ListConversations(limit)
	default:
		return nil, apierr.Forbidden("forbidden")
	}
	if err != nil {
		return nil, err
	}
	for i := range conversations {
		count, countErr := s.broker.GetCounter(ctx, realtime.ChatUnreadKey(actor.ID, conversations[i].OrderID))
		if countErr != nil || count == 0 {
			count, _ = s.store.CountUnreadMessages(conversations[i].OrderID, actor.ID)
			_ = s.broker.SetCounter(ctx, realtime.ChatUnreadKey(actor.ID, conversations[i].OrderID), count, 24*time.Hour)
		}
		conversations[i].UnreadCount = count
	}
	return conversations, nil
}

func (s *ChatService) ListMessages(actor domain.User, orderID, beforeID string, limit int) ([]domain.ChatMessage, error) {
	if _, err := s.authorizeOrderAccess(actor, orderID); err != nil {
		return nil, err
	}
	messages, err := s.store.ListMessages(orderID, actor.ID, limit, beforeID)
	if err != nil {
		return nil, err
	}
	if messages == nil {
		return []domain.ChatMessage{}, nil
	}
	return messages, nil
}

func (s *ChatService) SendMessage(ctx context.Context, actor domain.User, orderID, body string) (domain.ChatMessage, error) {
	if err := ensureActiveUser(s.store, actor); err != nil {
		return domain.ChatMessage{}, err
	}
	order, err := s.authorizeOrderAccess(actor, orderID)
	if err != nil {
		return domain.ChatMessage{}, err
	}
	body = strings.TrimSpace(body)
	if body == "" {
		return domain.ChatMessage{}, apierr.BadRequest("message body is required")
	}

	message, err := s.store.CreateMessage(orderID, actor.ID, body)
	if err != nil {
		return domain.ChatMessage{}, err
	}

	for _, recipientID := range participantIDs(order, actor.ID) {
		_, _ = s.broker.IncrementCounter(ctx, realtime.ChatUnreadKey(recipientID, orderID), 1)
		_ = s.broker.PublishJSON(ctx, realtime.NotificationsChannel(recipientID), map[string]interface{}{
			"type": "chat_unread",
			"order_id": orderID,
		})
		s.notifier.Publish(recipientID, "new_message", "New message in order chat")
	}
	_ = s.broker.ResetCounter(ctx, realtime.ChatUnreadKey(actor.ID, orderID))
	_ = s.broker.PublishJSON(ctx, realtime.ChatChannel(orderID), map[string]interface{}{
		"type":    "message",
		"message": message,
	})
	return message, nil
}

func (s *ChatService) MarkRead(ctx context.Context, actor domain.User, orderID string) error {
	if _, err := s.authorizeOrderAccess(actor, orderID); err != nil {
		return err
	}
	readAt := time.Now().UTC()
	if err := s.store.MarkChatRead(orderID, actor.ID, readAt); err != nil {
		return err
	}
	_ = s.broker.ResetCounter(ctx, realtime.ChatUnreadKey(actor.ID, orderID))
	_ = s.broker.PublishJSON(ctx, realtime.ChatChannel(orderID), map[string]interface{}{
		"type":      "read",
		"order_id":  orderID,
		"user_id":   actor.ID,
		"read_at":   readAt,
		"unread_count": 0,
	})
	return nil
}

func (s *ChatService) authorizeOrderAccess(actor domain.User, orderID string) (domain.Order, error) {
	order, err := s.store.GetOrder(orderID)
	if err != nil {
		return domain.Order{}, apierr.NotFound("order not found")
	}
	if actor.Role != domain.RoleAdmin && actor.ID != order.CustomerID && actor.ID != order.EngineerID {
		return domain.Order{}, apierr.Forbidden("forbidden")
	}
	return order, nil
}

func participantIDs(order domain.Order, senderID string) []string {
	participants := make([]string, 0, 2)
	if order.CustomerID != senderID {
		participants = append(participants, order.CustomerID)
	}
	if order.EngineerID != senderID && order.EngineerID != order.CustomerID {
		participants = append(participants, order.EngineerID)
	}
	return participants
}
