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
	if err := httprequest.DecodeJSON(r, &req); err != nil {
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
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, card)
}

func (h *CardHandler) List(w http.ResponseWriter, r *http.Request) {
	cards, err := h.service.List(r.URL.Query().Get("card_type"), r.URL.Query().Get("q"))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if cards == nil {
		cards = []domain.Card{}
	}
	response.JSON(w, http.StatusOK, cards)
}

func (h *CardHandler) Get(w http.ResponseWriter, r *http.Request) {
	card, err := h.service.Get(chi.URLParam(r, "id"))
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, card)
}

func (h *CardHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	var req cardRequest
	if err := httprequest.DecodeJSON(r, &req); err != nil {
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
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, card)
}
