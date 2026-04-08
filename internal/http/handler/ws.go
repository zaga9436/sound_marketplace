package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/soundmarket/backend/internal/http/response"
	"github.com/soundmarket/backend/internal/service"
)

type WSHandler struct {
	service *service.RealtimeService
}

func NewWSHandler(service *service.RealtimeService) *WSHandler {
	return &WSHandler{service: service}
}

func (h *WSHandler) Connect(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, h.service.ChatInfo(chi.URLParam(r, "id")))
}
