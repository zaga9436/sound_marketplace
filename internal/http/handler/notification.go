package handler

import (
	"net/http"

	httprequest "github.com/soundmarket/backend/internal/http/request"
	"github.com/soundmarket/backend/internal/http/middleware"
	"github.com/soundmarket/backend/internal/http/response"
	"github.com/soundmarket/backend/internal/service"
)

type NotificationHandler struct {
	service *service.NotificationService
}

type markNotificationsReadRequest struct {
	IDs []string `json:"ids"`
}

func NewNotificationHandler(service *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	result, err := h.service.List(r.Context(), user, queryInt(r, "limit", 20), r.URL.Query().Get("before_id"))
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *NotificationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	var req markNotificationsReadRequest
	if err := httprequest.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	result, err := h.service.MarkRead(r.Context(), user, req.IDs)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}
