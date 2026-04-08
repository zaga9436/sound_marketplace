package handler

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/soundmarket/backend/internal/domain"
	httprequest "github.com/soundmarket/backend/internal/http/request"
	"github.com/soundmarket/backend/internal/http/middleware"
	"github.com/soundmarket/backend/internal/http/response"
	"github.com/soundmarket/backend/internal/service"
)

type DisputeHandler struct {
	service *service.DisputeService
}

type openDisputeRequest struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type closeDisputeRequest struct {
	Resolution domain.DisputeResolution `json:"resolution"`
	Message    string                   `json:"message"`
}

func NewDisputeHandler(service *service.DisputeService) *DisputeHandler {
	return &DisputeHandler{service: service}
}

func (h *DisputeHandler) Open(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	var req openDisputeRequest
	if err := httprequest.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		reason = strings.TrimSpace(req.Message)
	}
	dispute, err := h.service.Open(user, chi.URLParam(r, "id"), reason)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, dispute)
}

func (h *DisputeHandler) Get(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	dispute, err := h.service.Get(user, chi.URLParam(r, "id"))
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, dispute)
}

func (h *DisputeHandler) Close(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	var req closeDisputeRequest
	if err := httprequest.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	dispute, err := h.service.Close(user, chi.URLParam(r, "id"), req.Resolution)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, dispute)
}
