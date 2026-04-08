package handler

import (
	"net/http"

	"github.com/soundmarket/backend/internal/http/response"
	"github.com/soundmarket/backend/internal/service"
)

type HealthHandler struct {
	service *service.HealthService
}

func NewHealthHandler(service *service.HealthService) *HealthHandler {
	return &HealthHandler{service: service}
}

func (h *HealthHandler) Get(w http.ResponseWriter, _ *http.Request) {
	response.JSON(w, http.StatusOK, h.service.Status())
}
