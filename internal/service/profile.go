package service

import (
	"github.com/soundmarket/backend/internal/apierr"
	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/repository"
)

type ProfileService struct {
	store repository.Store
}

func NewProfileService(store repository.Store) *ProfileService {
	return &ProfileService{store: store}
}

func (s *ProfileService) Get(userID string) (domain.Profile, error) {
	profile, err := s.store.GetProfile(userID)
	if err != nil {
		return domain.Profile{}, apierr.NotFound("profile not found")
	}
	return profile, nil
}

func (s *ProfileService) ListCards(userID string) ([]domain.Card, error) {
	if _, err := s.store.GetProfile(userID); err != nil {
		return nil, apierr.NotFound("profile not found")
	}
	cards, err := s.store.ListCardsByAuthor(userID)
	if err != nil {
		return nil, err
	}
	if cards == nil {
		return []domain.Card{}, nil
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
