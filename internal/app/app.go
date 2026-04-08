package app

import (
	"net/http"

	"github.com/soundmarket/backend/internal/auth"
	"github.com/soundmarket/backend/internal/config"
	"github.com/soundmarket/backend/internal/http/router"
	"github.com/soundmarket/backend/internal/notifications"
	"github.com/soundmarket/backend/internal/payments"
	"github.com/soundmarket/backend/internal/service"
	"github.com/soundmarket/backend/internal/storage"
	"github.com/soundmarket/backend/internal/worker"
)

type App struct {
	Config *config.Config
	Router http.Handler
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	authManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTTTL)
	storageAdapter := storage.NewS3Adapter(cfg)
	paymentAdapter := payments.NewMockYooKassaAdapter(cfg)
	notifier := notifications.NewInMemoryService()
	workerQueue := worker.NewInMemoryQueue()

	services := service.NewRegistry(service.Dependencies{
		Config:         cfg,
		AuthManager:    authManager,
		StorageAdapter: storageAdapter,
		PaymentAdapter: paymentAdapter,
		Notifier:       notifier,
		WorkerQueue:    workerQueue,
	})

	return &App{
		Config: cfg,
		Router: router.New(services),
	}, nil
}
