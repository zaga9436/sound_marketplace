package handler

import (
	"encoding/json"
	"net/http"

	"github.com/soundmarket/backend/internal/http/middleware"
	"github.com/soundmarket/backend/internal/http/response"
	"github.com/soundmarket/backend/internal/service"
)

type PaymentHandler struct {
	service *service.PaymentService
}

type depositRequest struct {
	Amount int64 `json:"amount"`
}

type webhookRequest struct {
	UserID     string `json:"user_id"`
	ExternalID string `json:"external_id"`
	Amount     int64  `json:"amount"`
}

func NewPaymentHandler(service *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) CreateDeposit(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	var req depositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	payment, err := h.service.CreateDeposit(user, req.Amount)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, payment)
}

func (h *PaymentHandler) Webhook(w http.ResponseWriter, r *http.Request) {
	var req webhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	response.JSON(w, http.StatusOK, h.service.ProcessWebhook(req.UserID, req.ExternalID, req.Amount))
}

func (h *PaymentHandler) Balance(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	response.JSON(w, http.StatusOK, map[string]int64{"balance": h.service.Balance(user.ID)})
}
