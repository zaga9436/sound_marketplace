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
	orderHandler := handler.NewOrderHandler(services.Order)
	paymentHandler := handler.NewPaymentHandler(services.Payment)
	wsHandler := handler.NewWSHandler(services.Realtime)

	r.Get("/health", healthHandler.Get)

	r.Route("/api/v1", func(api chi.Router) {
		api.Post("/auth/register", authHandler.Register)
		api.Post("/auth/login", authHandler.Login)
		api.Get("/cards", cardHandler.List)
		api.Get("/cards/{id}", cardHandler.Get)
		api.Get("/profiles/{id}", profileHandler.Public)
		api.Post("/payments/webhook", paymentHandler.Webhook)

		api.Group(func(secure chi.Router) {
			secure.Use(middleware.RequireAuth(services.AuthManager))
			secure.Get("/auth/me", authHandler.Me)
			secure.Get("/users/me", authHandler.Me)
			secure.Get("/profiles/me", profileHandler.Me)
			secure.Put("/profiles/me", profileHandler.Update)
			secure.Post("/cards", cardHandler.Create)
			secure.Put("/cards/{id}", cardHandler.Update)
			secure.Get("/requests/{id}/bids", bidHandler.List)
			secure.Post("/requests/{id}/bids", bidHandler.Create)
			secure.Post("/orders/from-offer", orderHandler.CreateFromOffer)
			secure.Post("/orders/from-bid", orderHandler.CreateFromBid)
			secure.Get("/orders/{id}", orderHandler.Get)
			secure.Patch("/orders/{id}/status", orderHandler.UpdateStatus)
			secure.Post("/payments/deposits", paymentHandler.CreateDeposit)
			secure.Get("/payments/balance", paymentHandler.Balance)
			secure.Get("/ws/orders/{id}", wsHandler.Connect)
		})
	})

	return r
}
