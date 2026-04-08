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

type CardHandler struct {
	service *service.CardService
}

type cardRequest struct {
	CardType    domain.CardType `json:"card_type"`
	Kind        domain.CardKind `json:"kind"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Price       int64           `json:"price"`
	Tags        []string        `json:"tags"`
	IsPublished bool            `json:"is_published"`
}

func NewCardHandler(service *service.CardService) *CardHandler {
	return &CardHandler{service: service}
}

func (h *CardHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	var req cardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	card, err := h.service.Create(user, domain.Card{
		CardType:    req.CardType,
		Kind:        req.Kind,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Tags:        req.Tags,
		IsPublished: req.IsPublished,
	})
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, card)
}

func (h *CardHandler) List(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, h.service.List(r.URL.Query().Get("card_type"), r.URL.Query().Get("q")))
}

func (h *CardHandler) Get(w http.ResponseWriter, r *http.Request) {
	card, err := h.service.Get(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, card)
}

func (h *CardHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	var req cardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	card, err := h.service.Update(user, chi.URLParam(r, "id"), domain.Card{
		Kind:        req.Kind,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Tags:        req.Tags,
		IsPublished: req.IsPublished,
	})
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, card)
}
