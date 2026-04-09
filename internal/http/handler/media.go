package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/http/middleware"
	"github.com/soundmarket/backend/internal/http/response"
	"github.com/soundmarket/backend/internal/service"
)

type MediaHandler struct {
	service *service.MediaService
}

func NewMediaHandler(service *service.MediaService) *MediaHandler {
	return &MediaHandler{service: service}
}

func (h *MediaHandler) UploadPreview(w http.ResponseWriter, r *http.Request) {
	h.upload(w, r, domain.MediaRolePreview)
}

func (h *MediaHandler) UploadFull(w http.ResponseWriter, r *http.Request) {
	h.upload(w, r, domain.MediaRoleFull)
}

func (h *MediaHandler) DownloadFull(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	url, err := h.service.DownloadCardFullMedia(r.Context(), user, chi.URLParam(r, "id"))
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"url": url})
}

func (h *MediaHandler) upload(w http.ResponseWriter, r *http.Request, role domain.MediaRole) {
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

	media, err := h.service.UploadCardMedia(r.Context(), user, chi.URLParam(r, "id"), role, service.MediaUploadInput{
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
