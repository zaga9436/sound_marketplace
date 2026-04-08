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

func (s *ProfileService) Update(userID, displayName, bio string) (domain.Profile, error) {
	profile, err := s.store.UpdateProfile(userID, displayName, bio)
	if err != nil {
		return domain.Profile{}, apierr.NotFound("profile not found")
	}
	return profile, nil
}
