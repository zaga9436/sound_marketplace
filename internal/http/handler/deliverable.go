package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/http/middleware"
	"github.com/soundmarket/backend/internal/http/response"
	"github.com/soundmarket/backend/internal/service"
)

type DeliverableHandler struct {
	service *service.DeliverableService
}

func NewDeliverableHandler(service *service.DeliverableService) *DeliverableHandler {
	return &DeliverableHandler{service: service}
}

func (h *DeliverableHandler) Upload(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid multipart form")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()
	deliverable, err := h.service.Upload(r.Context(), user, chi.URLParam(r, "id"), service.MediaUploadInput{
		Filename:    header.Filename,
		ContentType: header.Header.Get("Content-Type"),
		SizeBytes:   header.Size,
		Reader:      file,
	})
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, deliverable)
}

func (h *DeliverableHandler) List(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	deliverables, err := h.service.List(user, chi.URLParam(r, "id"))
	if err != nil {
		response.FromError(w, err)
		return
	}
	if deliverables == nil {
		deliverables = []domain.Deliverable{}
	}
	response.JSON(w, http.StatusOK, deliverables)
}

func (h *DeliverableHandler) Download(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	url, err := h.service.Download(r.Context(), user, chi.URLParam(r, "id"), chi.URLParam(r, "deliverable_id"))
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"url": url})
}
