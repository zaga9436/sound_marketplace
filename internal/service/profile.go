package service

import (
	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/repository"
)

type ProfileService struct {
	store *repository.MemoryStore
}

func NewProfileService(store *repository.MemoryStore) *ProfileService {
	return &ProfileService{store: store}
}

func (s *ProfileService) Get(userID string) (domain.Profile, error) {
	return s.store.GetProfile(userID)
}

func (s *ProfileService) Update(userID, displayName, bio string) (domain.Profile, error) {
	return s.store.UpdateProfile(userID, displayName, bio)
}
