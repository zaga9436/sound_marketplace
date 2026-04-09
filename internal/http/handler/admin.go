package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/soundmarket/backend/internal/domain"
	httprequest "github.com/soundmarket/backend/internal/http/request"
	"github.com/soundmarket/backend/internal/http/middleware"
	"github.com/soundmarket/backend/internal/http/response"
	"github.com/soundmarket/backend/internal/service"
)

type AdminHandler struct {
	service *service.AdminService
}

type moderationReasonRequest struct {
	Reason string `json:"reason"`
}

type adminCloseDisputeRequest struct {
	Resolution domain.DisputeResolution `json:"resolution"`
	Reason     string                   `json:"reason"`
	Message    string                   `json:"message"`
}

func NewAdminHandler(service *service.AdminService) *AdminHandler {
	return &AdminHandler{service: service}
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	users, err := h.service.ListUsers(user, r.URL.Query().Get("role"), r.URL.Query().Get("status"))
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, users)
}

func (h *AdminHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	result, err := h.service.GetUser(user, chi.URLParam(r, "id"))
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *AdminHandler) SuspendUser(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	var req moderationReasonRequest
	if err := httprequest.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	result, err := h.service.SuspendUser(user, chi.URLParam(r, "id"), req.Reason)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *AdminHandler) UnsuspendUser(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	var req moderationReasonRequest
	if err := httprequest.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	result, err := h.service.UnsuspendUser(user, chi.URLParam(r, "id"), req.Reason)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *AdminHandler) ListCards(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	cards, err := h.service.ListCards(user, r.URL.Query().Get("card_type"), r.URL.Query().Get("q"), r.URL.Query().Get("visibility"))
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, cards)
}

func (h *AdminHandler) GetCard(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	card, err := h.service.GetCard(user, chi.URLParam(r, "id"))
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, card)
}

func (h *AdminHandler) HideCard(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	var req moderationReasonRequest
	if err := httprequest.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	card, err := h.service.HideCard(user, chi.URLParam(r, "id"), req.Reason)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, card)
}

func (h *AdminHandler) UnhideCard(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	var req moderationReasonRequest
	if err := httprequest.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	card, err := h.service.UnhideCard(user, chi.URLParam(r, "id"), req.Reason)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, card)
}

func (h *AdminHandler) ListDisputes(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	disputes, err := h.service.ListDisputes(user, r.URL.Query().Get("status"))
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, disputes)
}

func (h *AdminHandler) GetDispute(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	dispute, err := h.service.GetDispute(user, chi.URLParam(r, "id"))
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, dispute)
}

func (h *AdminHandler) CloseDispute(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	var req adminCloseDisputeRequest
	if err := httprequest.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		reason = strings.TrimSpace(req.Message)
	}
	dispute, err := h.service.CloseDispute(user, chi.URLParam(r, "id"), req.Resolution, reason)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, dispute)
}

func (h *AdminHandler) ListModerationActions(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	actions, err := h.service.ListModerationActions(user, r.URL.Query().Get("target_type"), r.URL.Query().Get("target_id"), limit)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, actions)
}
