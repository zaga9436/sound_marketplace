package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/soundmarket/backend/internal/apierr"
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
	query, err := cardQueryFromRequest(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	cards, err := h.service.List(query)
	if err != nil {
		response.FromError(w, err)
		return
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

func cardQueryFromRequest(r *http.Request) (domain.CardQuery, error) {
	values := r.URL.Query()
	query := domain.CardQuery{
		CardType:  domain.CardType(values.Get("card_type")),
		Kind:      domain.CardKind(values.Get("kind")),
		AuthorID:  values.Get("author_id"),
		Query:     values.Get("q"),
		Tag:       values.Get("tag"),
		SortBy:    values.Get("sort_by"),
		SortOrder: values.Get("sort_order"),
	}
	if limit := values.Get("limit"); limit != "" {
		parsed, err := strconv.Atoi(limit)
		if err != nil {
			return domain.CardQuery{}, apierr.BadRequest("invalid limit")
		}
		query.Limit = parsed
	}
	if offset := values.Get("offset"); offset != "" {
		parsed, err := strconv.Atoi(offset)
		if err != nil {
			return domain.CardQuery{}, apierr.BadRequest("invalid offset")
		}
		query.Offset = parsed
	}
	if minPrice := values.Get("min_price"); minPrice != "" {
		parsed, err := strconv.ParseInt(minPrice, 10, 64)
		if err != nil {
			return domain.CardQuery{}, apierr.BadRequest("invalid min_price")
		}
		query.MinPrice = &parsed
	}
	if maxPrice := values.Get("max_price"); maxPrice != "" {
		parsed, err := strconv.ParseInt(maxPrice, 10, 64)
		if err != nil {
			return domain.CardQuery{}, apierr.BadRequest("invalid max_price")
		}
		query.MaxPrice = &parsed
	}
	return query, nil
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
