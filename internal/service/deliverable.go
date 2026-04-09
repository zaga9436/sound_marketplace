package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/soundmarket/backend/internal/apierr"
	"github.com/soundmarket/backend/internal/config"
	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/notifications"
	"github.com/soundmarket/backend/internal/repository"
	"github.com/soundmarket/backend/internal/storage"
)

type DeliverableService struct {
	cfg      *config.Config
	store    repository.Store
	storage  storage.Adapter
	notifier notifications.Service
}

func NewDeliverableService(cfg *config.Config, store repository.Store, storageAdapter storage.Adapter, notifier notifications.Service) *DeliverableService {
	return &DeliverableService{cfg: cfg, store: store, storage: storageAdapter, notifier: notifier}
}

func (s *DeliverableService) Upload(ctx context.Context, actor domain.User, orderID string, input MediaUploadInput) (domain.Deliverable, error) {
	if err := ensureActiveUser(s.store, actor); err != nil {
		return domain.Deliverable{}, err
	}
	if err := s.validateUpload(input); err != nil {
		return domain.Deliverable{}, err
	}
	if strings.TrimSpace(input.ContentType) == "" {
		input.ContentType = inferAudioContentType(input.Filename)
	}

	var created domain.Deliverable
	err := s.store.WithTx(func(tx repository.Store) error {
		order, err := tx.GetOrder(orderID)
		if err != nil {
			return apierr.NotFound("order not found")
		}
		if actor.Role != domain.RoleAdmin && actor.ID != order.EngineerID {
			return apierr.Forbidden("only engineer can upload deliverables")
		}
		if order.Status != domain.OrderStatusInProgress && order.Status != domain.OrderStatusReview && order.Status != domain.OrderStatusDispute {
			return apierr.BadRequest("order status does not allow deliverable upload")
		}

		nextVersion := 1
		latest, err := tx.GetLatestDeliverableByOrder(orderID)
		if err == nil {
			nextVersion = latest.Version + 1
		} else if !errors.Is(err, repository.ErrNotFound) {
			return err
		}

		ext := strings.ToLower(filepath.Ext(input.Filename))
		key := storage.BuildObjectKey(fmt.Sprintf("orders/%s/deliverables/v%d/%s", orderID, nextVersion, uuid.NewString()), "deliverable"+ext)
		object, err := s.storage.Upload(ctx, key, input.ContentType, input.Reader, storage.UploadOptions{Public: false})
		if err != nil {
			return err
		}
		if err := tx.DeactivateDeliverablesByOrder(orderID); err != nil {
			_ = s.storage.Delete(ctx, object.Key)
			return err
		}

		created, err = tx.CreateDeliverable(domain.Deliverable{
			OrderID:          orderID,
			UploadedBy:       actor.ID,
			StorageKey:       object.Key,
			OriginalFilename: input.Filename,
			ContentType:      input.ContentType,
			SizeBytes:        input.SizeBytes,
			Version:          nextVersion,
			IsActive:         true,
		})
		if err != nil {
			_ = s.storage.Delete(ctx, object.Key)
			return err
		}
		return nil
	})
	if err != nil {
		return domain.Deliverable{}, err
	}

	eventType := "deliverable_uploaded"
	message := "New deliverable uploaded"
	if created.Version > 1 {
		eventType = "deliverable_updated"
		message = "New deliverable version uploaded"
	}
	order, orderErr := s.store.GetOrder(orderID)
	if orderErr == nil {
		s.notifier.Publish(order.CustomerID, eventType, message)
	}
	return created, nil
}

func (s *DeliverableService) List(actor domain.User, orderID string) ([]domain.Deliverable, error) {
	order, err := s.store.GetOrder(orderID)
	if err != nil {
		return nil, apierr.NotFound("order not found")
	}
	if actor.Role != domain.RoleAdmin && actor.ID != order.CustomerID && actor.ID != order.EngineerID {
		return nil, apierr.Forbidden("forbidden")
	}
	deliverables, err := s.store.ListDeliverablesByOrder(orderID)
	if err != nil {
		return nil, err
	}
	if deliverables == nil {
		return []domain.Deliverable{}, nil
	}
	return deliverables, nil
}

func (s *DeliverableService) Download(ctx context.Context, actor domain.User, orderID, deliverableID string) (string, error) {
	order, err := s.store.GetOrder(orderID)
	if err != nil {
		return "", apierr.NotFound("order not found")
	}
	if actor.Role != domain.RoleAdmin && actor.ID != order.CustomerID && actor.ID != order.EngineerID {
		return "", apierr.Forbidden("forbidden")
	}
	deliverable, err := s.store.GetDeliverable(deliverableID)
	if err != nil {
		return "", apierr.NotFound("deliverable not found")
	}
	if deliverable.OrderID != orderID {
		return "", apierr.NotFound("deliverable not found")
	}
	url, err := s.storage.GenerateSignedURL(ctx, deliverable.StorageKey, s.cfg.SignedURLTTL)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (s *DeliverableService) validateUpload(input MediaUploadInput) error {
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
