package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"

	"github.com/soundmarket/backend/internal/http/middleware"
	"github.com/soundmarket/backend/internal/service"
)

type WSHandler struct {
	realtime *service.RealtimeService
	chat     *service.ChatService
}

type wsInboundMessage struct {
	Type string `json:"type"`
	Body string `json:"body"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func NewWSHandler(realtimeService *service.RealtimeService, chatService *service.ChatService) *WSHandler {
	return &WSHandler{realtime: realtimeService, chat: chatService}
}

func (h *WSHandler) ConnectOrder(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	orderID := chi.URLParam(r, "id")
	bootstrap, err := h.realtime.ChatBootstrap(r.Context(), user, orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	pubsub := h.realtime.Broker().Subscribe(ctx, serviceChatChannel(orderID))
	defer pubsub.Close()

	_ = conn.WriteJSON(bootstrap)

	go func() {
		for {
			var inbound wsInboundMessage
			if err := conn.ReadJSON(&inbound); err != nil {
				cancel()
				return
			}
			switch inbound.Type {
			case "message":
				_, _ = h.chat.SendMessage(ctx, user, orderID, inbound.Body)
			case "mark_read":
				_ = h.chat.MarkRead(ctx, user, orderID)
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-pubsub.Channel():
			if !ok {
				return
			}
			var payload map[string]interface{}
			if err := json.Unmarshal([]byte(msg.Payload), &payload); err == nil {
				_ = conn.WriteJSON(payload)
			}
		}
	}
}

func (h *WSHandler) ConnectNotifications(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	bootstrap, err := h.realtime.NotificationsBootstrap(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	pubsub := h.realtime.Broker().Subscribe(ctx, serviceNotificationsChannel(user.ID))
	defer pubsub.Close()

	_ = conn.WriteJSON(bootstrap)

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-pubsub.Channel():
			if !ok {
				return
			}
			var payload map[string]interface{}
			if err := json.Unmarshal([]byte(msg.Payload), &payload); err == nil {
				_ = conn.WriteJSON(payload)
			}
		}
	}
}

func serviceChatChannel(orderID string) string {
	return "chat:order:" + orderID
}

func serviceNotificationsChannel(userID string) string {
	return "notifications:user:" + userID
}

func closePubSub(pubsub *redis.PubSub) {
	if pubsub != nil {
		_ = pubsub.Close()
	}
}
