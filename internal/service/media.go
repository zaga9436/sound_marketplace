package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/soundmarket/backend/internal/apierr"
	"github.com/soundmarket/backend/internal/config"
	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/repository"
	"github.com/soundmarket/backend/internal/storage"
)

type MediaUploadInput struct {
	Filename    string
	ContentType string
	SizeBytes   int64
	Reader      io.Reader
}

type MediaService struct {
	cfg     *config.Config
	store   repository.Store
	storage storage.Adapter
}

func NewMediaService(cfg *config.Config, store repository.Store, storageAdapter storage.Adapter) *MediaService {
	return &MediaService{cfg: cfg, store: store, storage: storageAdapter}
}

func (s *MediaService) UploadCardMedia(ctx context.Context, actor domain.User, cardID string, role domain.MediaRole, input MediaUploadInput) (domain.MediaFile, error) {
	if err := ensureActiveUser(s.store, actor); err != nil {
		return domain.MediaFile{}, err
	}
	if role != domain.MediaRolePreview && role != domain.MediaRoleFull {
		return domain.MediaFile{}, apierr.BadRequest("media_role must be preview or full")
	}
	card, err := s.store.GetCard(cardID)
	if err != nil {
		return domain.MediaFile{}, apierr.NotFound("card not found")
	}
	if actor.Role != domain.RoleAdmin && actor.ID != card.AuthorID {
		return domain.MediaFile{}, apierr.Forbidden("forbidden")
	}
	if err := s.validateUpload(input); err != nil {
		return domain.MediaFile{}, err
	}
	if strings.TrimSpace(input.ContentType) == "" {
		input.ContentType = inferAudioContentType(input.Filename)
	}

	ext := strings.ToLower(filepath.Ext(input.Filename))
	key := storage.BuildObjectKey(fmt.Sprintf("cards/%s/%s/%s", cardID, role, uuid.NewString()), "original"+ext)

	object, err := s.storage.Upload(ctx, key, input.ContentType, input.Reader, storage.UploadOptions{
		Public: role == domain.MediaRolePreview,
	})
	if err != nil {
		return domain.MediaFile{}, err
	}

	created, err := s.store.CreateMedia(domain.MediaFile{
		CardID:           card.ID,
		OwnerUserID:      actor.ID,
		FileKey:          object.Key,
		OriginalFilename: input.Filename,
		ContentType:      input.ContentType,
		SizeBytes:        input.SizeBytes,
		MediaRole:        role,
	})
	if err != nil {
		_ = s.storage.Delete(ctx, object.Key)
		return domain.MediaFile{}, err
	}
	if role == domain.MediaRolePreview {
		created.URL = object.PublicURL
	}
	return created, nil
}

func (s *MediaService) DownloadCardFullMedia(ctx context.Context, actor domain.User, cardID string) (string, error) {
	card, err := s.store.GetCard(cardID)
	if err != nil {
		return "", apierr.NotFound("card not found")
	}
	if actor.Role != domain.RoleAdmin && actor.ID != card.AuthorID {
		allowed, err := s.store.UserHasCompletedCardAccess(card.ID, actor.ID)
		if err != nil {
			return "", err
		}
		if !allowed {
			return "", apierr.Forbidden("forbidden")
		}
	}

	media, err := s.store.GetLatestMediaByCardAndRole(card.ID, domain.MediaRoleFull)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", apierr.NotFound("full media not found")
		}
		return "", err
	}
	url, err := s.storage.GenerateSignedURL(ctx, media.FileKey, s.cfg.SignedURLTTL)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (s *MediaService) validateUpload(input MediaUploadInput) error {
	filename := strings.TrimSpace(input.Filename)
	if filename == "" {
		return apierr.BadRequest("file is required")
	}
	if input.SizeBytes <= 0 {
		return apierr.BadRequest("file is empty")
	}
	if input.SizeBytes > s.cfg.MaxUploadSize {
		return apierr.BadRequest("file exceeds max upload size")
	}
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(filename)), ".")
	if ext == "" {
		return apierr.BadRequest("file extension is required")
	}
	if !containsFold(s.cfg.AllowedAudioFormats, ext) {
		return apierr.BadRequest("unsupported audio format")
	}
	contentType := strings.ToLower(strings.TrimSpace(input.ContentType))
	if contentType != "" && !strings.HasPrefix(contentType, "audio/") {
		return apierr.BadRequest("unsupported content type")
	}
	return nil
}

func containsFold(values []string, value string) bool {
	for _, candidate := range values {
		if strings.EqualFold(candidate, value) {
			return true
		}
	}
	return false
}

func inferAudioContentType(filename string) string {
	switch strings.TrimPrefix(strings.ToLower(filepath.Ext(filename)), ".") {
	case "mp3":
		return "audio/mpeg"
	case "wav":
		return "audio/wav"
	case "flac":
		return "audio/flac"
	case "aac":
		return "audio/aac"
	default:
		return "application/octet-stream"
	}
}
