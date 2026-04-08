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
	Store          repository.Store
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
	Dispute     *DisputeService
	Payment     *PaymentService
	Health      *HealthService
	Realtime    *RealtimeService
	AuthManager *auth.JWTManager
}

func NewRegistry(deps Dependencies) *Registry {
	return &Registry{
		Auth:        NewAuthService(deps.Store, deps.AuthManager),
		Profile:     NewProfileService(deps.Store),
		Card:        NewCardService(deps.Store, deps.Notifier),
		Bid:         NewBidService(deps.Store, deps.Notifier),
		Order:       NewOrderService(deps.Store, deps.Notifier),
		Dispute:     NewDisputeService(deps.Store, deps.Notifier),
		Payment:     NewPaymentService(deps.Store, deps.PaymentAdapter, deps.Notifier),
		Health:      NewHealthService(deps.Config),
		Realtime:    NewRealtimeService(deps.WorkerQueue, deps.StorageAdapter),
		AuthManager: deps.AuthManager,
	}
}
