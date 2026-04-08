package notifications

import "sync"

type Event struct {
	UserID  string `json:"user_id"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Service interface {
	Publish(userID, eventType, message string)
	List(userID string) []Event
}

type InMemoryService struct {
	mu     sync.RWMutex
	events map[string][]Event
}

func NewInMemoryService() *InMemoryService {
	return &InMemoryService{events: map[string][]Event{}}
}

func (s *InMemoryService) Publish(userID, eventType, message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events[userID] = append(s.events[userID], Event{UserID: userID, Type: eventType, Message: message})
}

func (s *InMemoryService) List(userID string) []Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]Event(nil), s.events[userID]...)
}

