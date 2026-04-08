package notifications

import "github.com/soundmarket/backend/internal/repository"

type Service interface {
	Publish(userID, eventType, message string)
}

type RepositoryBackedService struct {
	store repository.Store
}

func NewRepositoryBackedService(store repository.Store) *RepositoryBackedService {
	return &RepositoryBackedService{store: store}
}

func (s *RepositoryBackedService) Publish(userID, eventType, message string) {
	_ = s.store.CreateNotification(userID, eventType, message)
}
