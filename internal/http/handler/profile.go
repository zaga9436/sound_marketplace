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

type ProfileHandler struct {
	service *service.ProfileService
}

type updateProfileRequest struct {
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
}

func NewProfileHandler(service *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{service: service}
}

func (h *ProfileHandler) Me(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	profile, err := h.service.Get(user.ID)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, profile)
}

func (h *ProfileHandler) Public(w http.ResponseWriter, r *http.Request) {
	profile, err := h.service.Get(chi.URLParam(r, "id"))
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, profile)
}

func (h *ProfileHandler) Cards(w http.ResponseWriter, r *http.Request) {
	cards, err := h.service.ListCards(chi.URLParam(r, "id"))
	if err != nil {
		response.FromError(w, err)
		return
	}
	if cards == nil {
		cards = []domain.Card{}
	}
	response.JSON(w, http.StatusOK, cards)
}

func (h *ProfileHandler) Reviews(w http.ResponseWriter, r *http.Request) {
	reviews, err := h.service.ListReviews(chi.URLParam(r, "id"))
	if err != nil {
		response.FromError(w, err)
		return
	}
	if reviews == nil {
		reviews = []domain.Review{}
	}
	response.JSON(w, http.StatusOK, reviews)
}

func (h *ProfileHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	var req updateProfileRequest
	if err := httprequest.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	profile, err := h.service.Update(user.ID, req.DisplayName, req.Bio)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, profile)
}
