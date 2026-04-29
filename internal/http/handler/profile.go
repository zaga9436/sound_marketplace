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

type ProfileHandler struct {
	service      *service.ProfileService
	mediaService *service.MediaService
}

type updateProfileRequest struct {
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
}

func NewProfileHandler(service *service.ProfileService, mediaService *service.MediaService) *ProfileHandler {
	return &ProfileHandler{service: service, mediaService: mediaService}
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
	query, err := profileCardQueryFromRequest(r)
	if err != nil {
		response.FromError(w, err)
		return
	}
	cards, err := h.service.ListCards(chi.URLParam(r, "id"), query)
	if err != nil {
		response.FromError(w, err)
		return
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

func (h *ProfileHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	if err := r.ParseMultipartForm(16 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid multipart form")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	media, err := h.mediaService.UploadProfileAvatar(r.Context(), user, service.MediaUploadInput{
		Filename:    header.Filename,
		ContentType: header.Header.Get("Content-Type"),
		SizeBytes:   header.Size,
		Reader:      file,
	})
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, media)
}

func profileCardQueryFromRequest(r *http.Request) (domain.CardQuery, error) {
	values := r.URL.Query()
	query := domain.CardQuery{
		CardType:  domain.CardType(values.Get("card_type")),
		Kind:      domain.CardKind(values.Get("kind")),
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
