package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/http/middleware"
	"github.com/soundmarket/backend/internal/http/response"
	"github.com/soundmarket/backend/internal/service"
)

type OrderHandler struct {
	service *service.OrderService
}

type createOfferOrderRequest struct {
	CardID string `json:"card_id"`
}

type createBidOrderRequest struct {
	BidID string `json:"bid_id"`
}

type updateStatusRequest struct {
	Status domain.OrderStatus `json:"status"`
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) CreateFromOffer(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	var req createOfferOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	order, err := h.service.CreateFromOffer(user, req.CardID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, order)
}

func (h *OrderHandler) CreateFromBid(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	var req createBidOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	order, err := h.service.CreateFromBid(user, req.BidID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, order)
}

func (h *OrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	order, err := h.service.Get(chi.URLParam(r, "id"), user)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, order)
}

func (h *OrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	var req updateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	order, err := h.service.UpdateStatus(user, chi.URLParam(r, "id"), req.Status)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, order)
}
