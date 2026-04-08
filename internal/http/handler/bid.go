package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/soundmarket/backend/internal/domain"
	httprequest "github.com/soundmarket/backend/internal/http/request"
	"github.com/soundmarket/backend/internal/http/middleware"
	"github.com/soundmarket/backend/internal/http/response"
	"github.com/soundmarket/backend/internal/service"
)

type BidHandler struct {
	service *service.BidService
}

type bidRequest struct {
	Price   int64  `json:"price"`
	Message string `json:"message"`
}

func NewBidHandler(service *service.BidService) *BidHandler {
	return &BidHandler{service: service}
}

func (h *BidHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	var req bidRequest
	if err := httprequest.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	bid, err := h.service.Create(user, chi.URLParam(r, "id"), req.Price, req.Message)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, bid)
}

func (h *BidHandler) List(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	bids, err := h.service.List(user, chi.URLParam(r, "id"))
	if err != nil {
		response.FromError(w, err)
		return
	}
	if bids == nil {
		bids = []domain.Bid{}
	}
	response.JSON(w, http.StatusOK, bids)
}
