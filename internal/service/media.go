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

func (s *MediaService) UploadProfileAvatar(ctx context.Context, actor domain.User, input MediaUploadInput) (domain.MediaFile, error) {
	if err := ensureActiveUser(s.store, actor); err != nil {
		return domain.MediaFile{}, err
	}
	if err := s.validateUpload(domain.MediaRoleAvatar, input); err != nil {
		return domain.MediaFile{}, err
	}
	if strings.TrimSpace(input.ContentType) == "" {
		input.ContentType = inferContentType(domain.MediaRoleAvatar, input.Filename)
	}

	ext := strings.ToLower(filepath.Ext(input.Filename))
	key := storage.BuildObjectKey(fmt.Sprintf("profiles/%s/avatar/%s", actor.ID, uuid.NewString()), "original"+ext)

	object, err := s.storage.Upload(ctx, key, input.ContentType, input.Reader, storage.UploadOptions{
		Public: true,
	})
	if err != nil {
		return domain.MediaFile{}, err
	}

	created, err := s.store.CreateMedia(domain.MediaFile{
		OwnerUserID:      actor.ID,
		FileKey:          object.Key,
		OriginalFilename: input.Filename,
		ContentType:      input.ContentType,
		SizeBytes:        input.SizeBytes,
		MediaRole:        domain.MediaRoleAvatar,
	})
	if err != nil {
		_ = s.storage.Delete(ctx, object.Key)
		return domain.MediaFile{}, err
	}
	created.URL = object.PublicURL
	return created, nil
}

func (s *MediaService) UploadCardMedia(ctx context.Context, actor domain.User, cardID string, role domain.MediaRole, input MediaUploadInput) (domain.MediaFile, error) {
	if err := ensureActiveUser(s.store, actor); err != nil {
		return domain.MediaFile{}, err
	}
	if role != domain.MediaRoleCover && role != domain.MediaRolePreview && role != domain.MediaRoleFull && role != domain.MediaRoleMaterial {
		return domain.MediaFile{}, apierr.BadRequest("media_role must be cover, preview, full or material")
	}
	card, err := s.store.GetCard(cardID)
	if err != nil {
		return domain.MediaFile{}, apierr.NotFound("card not found")
	}
	if actor.Role != domain.RoleAdmin && actor.ID != card.AuthorID {
		return domain.MediaFile{}, apierr.Forbidden("forbidden")
	}
	if role == domain.MediaRoleMaterial && card.CardType != domain.CardTypeRequest {
		return domain.MediaFile{}, apierr.BadRequest("materials are available only for request cards")
	}
	if (role == domain.MediaRolePreview || role == domain.MediaRoleFull) && (card.CardType != domain.CardTypeOffer || card.Kind != domain.CardKindProduct) {
		return domain.MediaFile{}, apierr.BadRequest("preview and full media are available only for product offers")
	}
	if err := s.validateUpload(role, input); err != nil {
		return domain.MediaFile{}, err
	}
	if strings.TrimSpace(input.ContentType) == "" {
		input.ContentType = inferContentType(role, input.Filename)
	}

	ext := strings.ToLower(filepath.Ext(input.Filename))
	key := storage.BuildObjectKey(fmt.Sprintf("cards/%s/%s/%s", cardID, role, uuid.NewString()), "original"+ext)

	object, err := s.storage.Upload(ctx, key, input.ContentType, input.Reader, storage.UploadOptions{
		Public: role == domain.MediaRolePreview || role == domain.MediaRoleCover,
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
	if role == domain.MediaRolePreview || role == domain.MediaRoleCover {
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
		allowed, err := s.store.UserHasStartedCardAccess(card.ID, actor.ID)
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

func (s *MediaService) ListCardMaterials(ctx context.Context, actor domain.User, cardID string) ([]domain.MediaFile, error) {
	card, err := s.store.GetCard(cardID)
	if err != nil {
		return nil, apierr.NotFound("card not found")
	}
	if actor.Role != domain.RoleAdmin && actor.ID != card.AuthorID {
		allowed, err := s.store.UserHasStartedCardAccess(card.ID, actor.ID)
		if err != nil {
			return nil, err
		}
		if !allowed {
			return nil, apierr.Forbidden("forbidden")
		}
	}
	materials, err := s.store.ListMediaByCardAndRole(card.ID, domain.MediaRoleMaterial)
	if err != nil {
		return nil, err
	}
	return materials, nil
}

func (s *MediaService) DownloadCardMaterial(ctx context.Context, actor domain.User, cardID, mediaID string) (string, error) {
	materials, err := s.ListCardMaterials(ctx, actor, cardID)
	if err != nil {
		return "", err
	}
	var selected domain.MediaFile
	for _, material := range materials {
		if material.ID == mediaID {
			selected = material
			break
		}
	}
	if selected.ID == "" {
		return "", apierr.NotFound("material not found")
	}
	url, err := s.storage.GenerateSignedURL(ctx, selected.FileKey, s.cfg.SignedURLTTL)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (s *MediaService) validateUpload(role domain.MediaRole, input MediaUploadInput) error {
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
	contentType := strings.ToLower(strings.TrimSpace(input.ContentType))

	switch role {
	case domain.MediaRoleAvatar, domain.MediaRoleCover:
		if !containsFold(s.cfg.AllowedImageFormats, ext) {
			return apierr.BadRequest("unsupported image format")
		}
		if contentType != "" && !strings.HasPrefix(contentType, "image/") {
			return apierr.BadRequest("unsupported content type")
		}
	case domain.MediaRolePreview, domain.MediaRoleFull:
		if !containsFold(s.cfg.AllowedAudioFormats, ext) {
			return apierr.BadRequest("unsupported audio format")
		}
		if contentType != "" && !strings.HasPrefix(contentType, "audio/") {
			return apierr.BadRequest("unsupported content type")
		}
	case domain.MediaRoleMaterial:
		if !containsFold(s.cfg.AllowedAudioFormats, ext) && ext != "zip" {
			return apierr.BadRequest("unsupported material format")
		}
		if contentType != "" && !strings.HasPrefix(contentType, "audio/") && contentType != "application/zip" && contentType != "application/x-zip-compressed" && contentType != "application/octet-stream" {
			return apierr.BadRequest("unsupported content type")
		}
	default:
		return apierr.BadRequest("unsupported media role")
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

func inferContentType(role domain.MediaRole, filename string) string {
	if role == domain.MediaRoleAvatar || role == domain.MediaRoleCover {
		return inferImageContentType(filename)
	}
	if role == domain.MediaRoleMaterial && strings.EqualFold(strings.TrimPrefix(filepath.Ext(filename), "."), "zip") {
		return "application/zip"
	}
	return inferAudioContentType(filename)
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

func inferImageContentType(filename string) string {
	switch strings.TrimPrefix(strings.ToLower(filepath.Ext(filename)), ".") {
	case "jpg", "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "webp":
		return "image/webp"
	case "gif":
		return "image/gif"
	default:
		return "application/octet-stream"
	}
}
