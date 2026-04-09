package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/soundmarket/backend/internal/http/handler"
	"github.com/soundmarket/backend/internal/http/middleware"
	"github.com/soundmarket/backend/internal/service"
)

func New(services *service.Registry) http.Handler {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))

	healthHandler := handler.NewHealthHandler(services.Health)
	authHandler := handler.NewAuthHandler(services.Auth)
	profileHandler := handler.NewProfileHandler(services.Profile)
	cardHandler := handler.NewCardHandler(services.Card)
	bidHandler := handler.NewBidHandler(services.Bid)
	chatHandler := handler.NewChatHandler(services.Chat)
	notificationHandler := handler.NewNotificationHandler(services.Notifications)
	orderHandler := handler.NewOrderHandler(services.Order)
	disputeHandler := handler.NewDisputeHandler(services.Dispute)
	reviewHandler := handler.NewReviewHandler(services.Review)
	mediaHandler := handler.NewMediaHandler(services.Media)
	paymentHandler := handler.NewPaymentHandler(services.Payment)
	wsHandler := handler.NewWSHandler(services.Realtime, services.Chat)

	r.Get("/health", healthHandler.Get)

	r.Route("/api/v1", func(api chi.Router) {
		api.Post("/auth/register", authHandler.Register)
		api.Post("/auth/login", authHandler.Login)
		api.Get("/cards", cardHandler.List)
		api.Get("/cards/{id}", cardHandler.Get)
		api.Get("/profiles/{id}", profileHandler.Public)
		api.Get("/profiles/{id}/cards", profileHandler.Cards)
		api.Get("/profiles/{id}/reviews", profileHandler.Reviews)
		api.Post("/payments/webhook", paymentHandler.Webhook)

		api.Group(func(secure chi.Router) {
			secure.Use(middleware.RequireAuth(services.AuthManager))
			secure.Get("/auth/me", authHandler.Me)
			secure.Get("/users/me", authHandler.Me)
			secure.Get("/profiles/me", profileHandler.Me)
			secure.Put("/profiles/me", profileHandler.Update)
			secure.Post("/cards", cardHandler.Create)
			secure.Put("/cards/{id}", cardHandler.Update)
			secure.Post("/cards/{id}/media/preview", mediaHandler.UploadPreview)
			secure.Post("/cards/{id}/media/full", mediaHandler.UploadFull)
			secure.Get("/cards/{id}/download", mediaHandler.DownloadFull)
			secure.Get("/chats", chatHandler.ListConversations)
			secure.Get("/notifications", notificationHandler.List)
			secure.Post("/notifications/read", notificationHandler.MarkRead)
			secure.Get("/requests/{id}/bids", bidHandler.List)
			secure.Post("/requests/{id}/bids", bidHandler.Create)
			secure.Route("/orders", func(orderRoutes chi.Router) {
				orderRoutes.Post("/from-offer", orderHandler.CreateFromOffer)
				orderRoutes.Post("/from-bid", orderHandler.CreateFromBid)
				orderRoutes.Get("/", orderHandler.List)
				orderRoutes.Route("/{id}", func(order chi.Router) {
					order.Get("/", orderHandler.Get)
					order.Get("/messages", chatHandler.ListMessages)
					order.Post("/messages", chatHandler.CreateMessage)
					order.Post("/messages/read", chatHandler.MarkRead)
					order.Patch("/status", orderHandler.UpdateStatus)
					order.Post("/reviews", reviewHandler.Create)
					order.Post("/dispute", disputeHandler.Open)
					order.Get("/dispute", disputeHandler.Get)
					order.Post("/dispute/close", disputeHandler.Close)
				})
			})
			secure.Post("/payments/deposits", paymentHandler.CreateDeposit)
			secure.Post("/payments/sync", paymentHandler.Sync)
			secure.Get("/payments/balance", paymentHandler.Balance)
			secure.Get("/ws/orders/{id}", wsHandler.ConnectOrder)
			secure.Get("/ws/notifications", wsHandler.ConnectNotifications)
		})
	})

	return r
}
