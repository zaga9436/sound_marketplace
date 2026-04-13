package service

import (
	"context"

	"github.com/soundmarket/backend/internal/apierr"
	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/repository"
	"github.com/soundmarket/backend/internal/storage"
)

type ProfileService struct {
	store   repository.Store
	storage storage.Adapter
}

func NewProfileService(store repository.Store, storageAdapter storage.Adapter) *ProfileService {
	return &ProfileService{store: store, storage: storageAdapter}
}

func (s *ProfileService) Get(userID string) (domain.Profile, error) {
	profile, err := s.store.GetProfile(userID)
	if err != nil {
		return domain.Profile{}, apierr.NotFound("profile not found")
	}
	return profile, nil
}

func (s *ProfileService) ListCards(userID string, query domain.CardQuery) (domain.CardList, error) {
	if _, err := s.store.GetProfile(userID); err != nil {
		return domain.CardList{}, apierr.NotFound("profile not found")
	}
	cards, err := s.store.ListCardsByAuthor(userID, query)
	if err != nil {
		return domain.CardList{}, err
	}
	if err := s.attachPreviewURLs(context.Background(), cards.Items); err != nil {
		return domain.CardList{}, err
	}
	if cards.Items == nil {
		cards.Items = []domain.Card{}
	}
	return cards, nil
}

func (s *ProfileService) ListReviews(userID string) ([]domain.Review, error) {
	if _, err := s.store.GetProfile(userID); err != nil {
		return nil, apierr.NotFound("profile not found")
	}
	reviews, err := s.store.ListReviewsByTargetUser(userID)
	if err != nil {
		return nil, err
	}
	if reviews == nil {
		return []domain.Review{}, nil
	}
	return reviews, nil
}

func (s *ProfileService) Update(userID, displayName, bio string) (domain.Profile, error) {
	profile, err := s.store.UpdateProfile(userID, displayName, bio)
	if err != nil {
		return domain.Profile{}, apierr.NotFound("profile not found")
	}
	return profile, nil
}

func (s *ProfileService) attachPreviewURLs(_ context.Context, cards []domain.Card) error {
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
