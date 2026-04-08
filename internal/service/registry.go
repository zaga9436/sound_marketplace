package service

import (
	"github.com/soundmarket/backend/internal/auth"
	"github.com/soundmarket/backend/internal/config"
	"github.com/soundmarket/backend/internal/notifications"
	"github.com/soundmarket/backend/internal/payments"
	"github.com/soundmarket/backend/internal/repository"
	"github.com/soundmarket/backend/internal/storage"
	"github.com/soundmarket/backend/internal/worker"
)

type Dependencies struct {
	Config         *config.Config
	AuthManager    *auth.JWTManager
	StorageAdapter storage.Adapter
	PaymentAdapter payments.Adapter
	Notifier       notifications.Service
	WorkerQueue    worker.Queue
}

type Registry struct {
	Auth        *AuthService
	Profile     *ProfileService
	Card        *CardService
	Bid         *BidService
	Order       *OrderService
	Payment     *PaymentService
	Health      *HealthService
	Realtime    *RealtimeService
	AuthManager *auth.JWTManager
}

func NewRegistry(deps Dependencies) *Registry {
	store := repository.NewMemoryStore()
	return &Registry{
		Auth:        NewAuthService(store, deps.AuthManager),
		Profile:     NewProfileService(store),
		Card:        NewCardService(store, deps.Notifier),
		Bid:         NewBidService(store, deps.Notifier),
		Order:       NewOrderService(store, deps.Notifier),
		Payment:     NewPaymentService(store, deps.PaymentAdapter, deps.Notifier),
		Health:      NewHealthService(deps.Config),
		Realtime:    NewRealtimeService(deps.WorkerQueue, deps.StorageAdapter),
		AuthManager: deps.AuthManager,
	}
}
