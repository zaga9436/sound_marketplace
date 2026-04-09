package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httprequest "github.com/soundmarket/backend/internal/http/request"
	"github.com/soundmarket/backend/internal/http/middleware"
	"github.com/soundmarket/backend/internal/http/response"
	"github.com/soundmarket/backend/internal/service"
)

type ReviewHandler struct {
	service *service.ReviewService
}

type createReviewRequest struct {
	Rating int    `json:"rating"`
	Text   string `json:"text"`
}

func NewReviewHandler(service *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: service}
}

func (h *ReviewHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	var req createReviewRequest
	if err := httprequest.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	review, err := h.service.Create(user, chi.URLParam(r, "id"), req.Rating, req.Text)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, review)
}
