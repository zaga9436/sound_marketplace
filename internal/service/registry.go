package service

import (
	"github.com/soundmarket/backend/internal/auth"
	"github.com/soundmarket/backend/internal/config"
	"github.com/soundmarket/backend/internal/notifications"
	"github.com/soundmarket/backend/internal/payments"
	"github.com/soundmarket/backend/internal/realtime"
	"github.com/soundmarket/backend/internal/repository"
	"github.com/soundmarket/backend/internal/storage"
	"github.com/soundmarket/backend/internal/worker"
)

type Dependencies struct {
	Config         *config.Config
	Store          repository.Store
	AuthManager    *auth.JWTManager
	Broker         *realtime.Broker
	StorageAdapter storage.Adapter
	PaymentAdapter payments.Provider
	Notifier       notifications.Service
	WorkerQueue    worker.Queue
}

type Registry struct {
	Auth        *AuthService
	Profile     *ProfileService
	Card        *CardService
	Bid         *BidService
	Order       *OrderService
	Dispute     *DisputeService
	Review      *ReviewService
	Media       *MediaService
	Chat        *ChatService
	Notifications *NotificationService
	Payment     *PaymentService
	Health      *HealthService
	Realtime    *RealtimeService
	AuthManager *auth.JWTManager
}

func NewRegistry(deps Dependencies) *Registry {
	chatService := NewChatService(deps.Store, deps.Broker, deps.Notifier)
	notificationService := NewNotificationService(deps.Store, deps.Broker)

	return &Registry{
		Auth:        NewAuthService(deps.Store, deps.AuthManager),
		Profile:     NewProfileService(deps.Store, deps.StorageAdapter),
		Card:        NewCardService(deps.Store, deps.Notifier, deps.StorageAdapter),
		Bid:         NewBidService(deps.Store, deps.Notifier),
		Order:       NewOrderService(deps.Store, deps.Notifier),
		Dispute:     NewDisputeService(deps.Store, deps.Notifier),
		Review:      NewReviewService(deps.Store, deps.Notifier),
		Media:       NewMediaService(deps.Config, deps.Store, deps.StorageAdapter),
		Chat:        chatService,
		Notifications: notificationService,
		Payment:     NewPaymentService(deps.Config, deps.Store, deps.PaymentAdapter, deps.Notifier),
		Health:      NewHealthService(deps.Config),
		Realtime:    NewRealtimeService(deps.Broker, chatService, notificationService),
		AuthManager: deps.AuthManager,
	}
}
