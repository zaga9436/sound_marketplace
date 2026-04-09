package service

import (
	"context"
	"strings"

	"github.com/soundmarket/backend/internal/apierr"
	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/repository"
	"github.com/soundmarket/backend/internal/storage"
)

type AdminService struct {
	store    repository.Store
	disputes *DisputeService
	storage  storage.Adapter
}

func NewAdminService(store repository.Store, disputes *DisputeService, storageAdapter storage.Adapter) *AdminService {
	return &AdminService{store: store, disputes: disputes, storage: storageAdapter}
}

func (s *AdminService) ListUsers(actor domain.User, role, status string) ([]domain.User, error) {
	if actor.Role != domain.RoleAdmin {
		return nil, apierr.Forbidden("forbidden")
	}
	users, err := s.store.ListUsers(strings.TrimSpace(role), strings.TrimSpace(status))
	if err != nil {
		return nil, err
	}
	if users == nil {
		return []domain.User{}, nil
	}
	return users, nil
}

func (s *AdminService) GetUser(actor domain.User, userID string) (domain.User, error) {
	if actor.Role != domain.RoleAdmin {
		return domain.User{}, apierr.Forbidden("forbidden")
	}
	user, err := s.store.GetUser(userID)
	if err != nil {
		return domain.User{}, apierr.NotFound("user not found")
	}
	return user, nil
}

func (s *AdminService) SuspendUser(actor domain.User, userID, reason string) (domain.User, error) {
	if actor.Role != domain.RoleAdmin {
		return domain.User{}, apierr.Forbidden("forbidden")
	}
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return domain.User{}, apierr.BadRequest("reason is required")
	}
	target, err := s.store.GetUser(userID)
	if err != nil {
		return domain.User{}, apierr.NotFound("user not found")
	}
	if target.Role == domain.RoleAdmin {
		return domain.User{}, apierr.BadRequest("cannot suspend admin user")
	}
	user, err := s.store.SetUserSuspended(userID, true, reason)
	if err != nil {
		return domain.User{}, err
	}
	_, _ = s.store.CreateModerationAction(domain.ModerationAction{
		AdminUserID: actor.ID,
		TargetType:  "user",
		TargetID:    userID,
		Action:      "suspend",
		Reason:      reason,
	})
	return user, nil
}

func (s *AdminService) UnsuspendUser(actor domain.User, userID, reason string) (domain.User, error) {
	if actor.Role != domain.RoleAdmin {
		return domain.User{}, apierr.Forbidden("forbidden")
	}
	if _, err := s.store.GetUser(userID); err != nil {
		return domain.User{}, apierr.NotFound("user not found")
	}
	user, err := s.store.SetUserSuspended(userID, false, strings.TrimSpace(reason))
	if err != nil {
		return domain.User{}, err
	}
	_, _ = s.store.CreateModerationAction(domain.ModerationAction{
		AdminUserID: actor.ID,
		TargetType:  "user",
		TargetID:    userID,
		Action:      "unsuspend",
		Reason:      strings.TrimSpace(reason),
	})
	return user, nil
}

func (s *AdminService) ListCards(actor domain.User, cardType, query, visibility string) ([]domain.Card, error) {
	if actor.Role != domain.RoleAdmin {
		return nil, apierr.Forbidden("forbidden")
	}
	cards, err := s.store.ListCardsForAdmin(strings.TrimSpace(cardType), strings.TrimSpace(query), strings.TrimSpace(visibility))
	if err != nil {
		return nil, err
	}
	if cards == nil {
		return []domain.Card{}, nil
	}
	if err := s.attachPreviewURLs(context.Background(), cards); err != nil {
		return nil, err
	}
	return cards, nil
}

func (s *AdminService) GetCard(actor domain.User, cardID string) (domain.Card, error) {
	if actor.Role != domain.RoleAdmin {
		return domain.Card{}, apierr.Forbidden("forbidden")
	}
	card, err := s.store.GetCard(cardID)
	if err != nil {
		return domain.Card{}, apierr.NotFound("card not found")
	}
	cards := []domain.Card{card}
	if err := s.attachPreviewURLs(context.Background(), cards); err != nil {
		return domain.Card{}, err
	}
	return cards[0], nil
}

func (s *AdminService) HideCard(actor domain.User, cardID, reason string) (domain.Card, error) {
	if actor.Role != domain.RoleAdmin {
		return domain.Card{}, apierr.Forbidden("forbidden")
	}
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return domain.Card{}, apierr.BadRequest("reason is required")
	}
	if _, err := s.store.GetCard(cardID); err != nil {
		return domain.Card{}, apierr.NotFound("card not found")
	}
	card, err := s.store.SetCardHidden(cardID, true, reason)
	if err != nil {
		return domain.Card{}, err
	}
	_, _ = s.store.CreateModerationAction(domain.ModerationAction{
		AdminUserID: actor.ID,
		TargetType:  "card",
		TargetID:    cardID,
		Action:      "hide",
		Reason:      reason,
	})
	return card, nil
}

func (s *AdminService) UnhideCard(actor domain.User, cardID, reason string) (domain.Card, error) {
	if actor.Role != domain.RoleAdmin {
		return domain.Card{}, apierr.Forbidden("forbidden")
	}
	if _, err := s.store.GetCard(cardID); err != nil {
		return domain.Card{}, apierr.NotFound("card not found")
	}
	card, err := s.store.SetCardHidden(cardID, false, strings.TrimSpace(reason))
	if err != nil {
		return domain.Card{}, err
	}
	_, _ = s.store.CreateModerationAction(domain.ModerationAction{
		AdminUserID: actor.ID,
		TargetType:  "card",
		TargetID:    cardID,
		Action:      "unhide",
		Reason:      strings.TrimSpace(reason),
	})
	return card, nil
}

func (s *AdminService) ListDisputes(actor domain.User, status string) ([]domain.Dispute, error) {
	if actor.Role != domain.RoleAdmin {
		return nil, apierr.Forbidden("forbidden")
	}
	disputes, err := s.store.ListDisputes(strings.TrimSpace(status))
	if err != nil {
		return nil, err
	}
	if disputes == nil {
		return []domain.Dispute{}, nil
	}
	return disputes, nil
}

func (s *AdminService) GetDispute(actor domain.User, disputeID string) (domain.Dispute, error) {
	if actor.Role != domain.RoleAdmin {
		return domain.Dispute{}, apierr.Forbidden("forbidden")
	}
	dispute, err := s.store.GetDispute(disputeID)
	if err != nil {
		return domain.Dispute{}, apierr.NotFound("dispute not found")
	}
	return dispute, nil
}

func (s *AdminService) CloseDispute(actor domain.User, disputeID string, resolution domain.DisputeResolution, reason string) (domain.Dispute, error) {
	if actor.Role != domain.RoleAdmin {
		return domain.Dispute{}, apierr.Forbidden("forbidden")
	}
	dispute, err := s.store.GetDispute(disputeID)
	if err != nil {
		return domain.Dispute{}, apierr.NotFound("dispute not found")
	}
	closed, err := s.disputes.Close(actor, dispute.OrderID, resolution)
	if err != nil {
		return domain.Dispute{}, err
	}
	_, _ = s.store.CreateModerationAction(domain.ModerationAction{
		AdminUserID: actor.ID,
		TargetType:  "dispute",
		TargetID:    disputeID,
		Action:      "close",
		Reason:      strings.TrimSpace(reason),
	})
	return closed, nil
}

func (s *AdminService) ListModerationActions(actor domain.User, targetType, targetID string, limit int) ([]domain.ModerationAction, error) {
	if actor.Role != domain.RoleAdmin {
		return nil, apierr.Forbidden("forbidden")
	}
	actions, err := s.store.ListModerationActions(strings.TrimSpace(targetType), strings.TrimSpace(targetID), limit)
	if err != nil {
		return nil, err
	}
	if actions == nil {
		return []domain.ModerationAction{}, nil
	}
	return actions, nil
}

func (s *AdminService) attachPreviewURLs(_ context.Context, cards []domain.Card) error {
	for i := range cards {
		mediaFiles, err := s.store.ListMediaByCardAndRole(cards[i].ID, domain.MediaRolePreview)
		if err != nil {
			return err
		}
		cards[i].PreviewURLs = make([]string, 0, len(mediaFiles))
		for _, media := range mediaFiles {
			cards[i].PreviewURLs = append(cards[i].PreviewURLs, s.storage.PublicURL(media.FileKey))
		}
	}
	return nil
}
