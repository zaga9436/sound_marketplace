package handler

import (
	"net/http"

	httprequest "github.com/soundmarket/backend/internal/http/request"
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
	ExternalID string `json:"external_id"`
}

func NewPaymentHandler(service *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) CreateDeposit(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	var req depositRequest
	if err := httprequest.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	payment, err := h.service.CreateDeposit(user, req.Amount)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, payment)
}

func (h *PaymentHandler) Webhook(w http.ResponseWriter, r *http.Request) {
	var req webhookRequest
	if err := httprequest.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	tx, err := h.service.ProcessWebhook(req.ExternalID)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, tx)
}

func (h *PaymentHandler) Balance(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	balance, err := h.service.Balance(user.ID)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]int64{"balance": balance})
}
