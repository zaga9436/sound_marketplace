package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	httprequest "github.com/soundmarket/backend/internal/http/request"
	"github.com/soundmarket/backend/internal/http/middleware"
	"github.com/soundmarket/backend/internal/http/response"
	"github.com/soundmarket/backend/internal/service"
)

type ChatHandler struct {
	service *service.ChatService
}

type createMessageRequest struct {
	Body string `json:"body"`
}

func NewChatHandler(service *service.ChatService) *ChatHandler {
	return &ChatHandler{service: service}
}

func (h *ChatHandler) ListConversations(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	limit := queryInt(r, "limit", 20)
	conversations, err := h.service.ListConversations(r.Context(), user, limit)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, conversations)
}

func (h *ChatHandler) ListMessages(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	messages, err := h.service.ListMessages(user, chi.URLParam(r, "id"), r.URL.Query().Get("before_id"), queryInt(r, "limit", 20))
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, messages)
}

func (h *ChatHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	var req createMessageRequest
	if err := httprequest.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	message, err := h.service.SendMessage(r.Context(), user, chi.URLParam(r, "id"), req.Body)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, message)
}

func (h *ChatHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	if err := h.service.MarkRead(r.Context(), user, chi.URLParam(r, "id")); err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"ok": true})
}

func queryInt(r *http.Request, key string, fallback int) int {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}
