package app

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/soundmarket/backend/internal/auth"
	"github.com/soundmarket/backend/internal/config"
	"github.com/soundmarket/backend/internal/http/router"
	"github.com/soundmarket/backend/internal/notifications"
	"github.com/soundmarket/backend/internal/payments"
	"github.com/soundmarket/backend/internal/platform/db"
	"github.com/soundmarket/backend/internal/repository"
	"github.com/soundmarket/backend/internal/service"
	"github.com/soundmarket/backend/internal/storage"
	"github.com/soundmarket/backend/internal/worker"
)

type App struct {
	Config *config.Config
	Router http.Handler
	DB     *sql.DB
	Redis  *redis.Client
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	authManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTTTL)
	postgresDB, err := db.OpenPostgres(cfg)
	if err != nil {
		return nil, err
	}
	if cfg.AutoApplyMigrations {
		log.Printf("app init: auto migrations enabled, dir=%s", cfg.MigrationsDir)
		if err := db.ApplyMigrations(postgresDB, cfg.MigrationsDir); err != nil {
			_ = postgresDB.Close()
			return nil, err
		}
	}
	redisClient, err := db.OpenRedis(cfg)
	if err != nil {
		_ = postgresDB.Close()
		return nil, err
	}
	store := repository.NewPostgresStore(postgresDB)
	if err := ensureDevAdmin(cfg, store); err != nil {
		_ = postgresDB.Close()
		_ = redisClient.Close()
		return nil, err
	}
	storageAdapter, err := storage.NewS3Adapter(cfg)
	if err != nil {
		_ = postgresDB.Close()
		_ = redisClient.Close()
		return nil, err
	}
	paymentAdapter := payments.NewMockYooKassaAdapter(cfg)
	notifier := notifications.NewRepositoryBackedService(store)
	workerQueue := worker.NewInMemoryQueue()

	services := service.NewRegistry(service.Dependencies{
		Config:         cfg,
		Store:          store,
		AuthManager:    authManager,
		StorageAdapter: storageAdapter,
		PaymentAdapter: paymentAdapter,
		Notifier:       notifier,
		WorkerQueue:    workerQueue,
	})

	return &App{
		Config: cfg,
		Router: router.New(services),
		DB:     postgresDB,
		Redis:  redisClient,
	}, nil
}
