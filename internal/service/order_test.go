package service

import (
	"testing"

	"github.com/soundmarket/backend/internal/domain"
)

func TestStatusTransitionRules(t *testing.T) {
	if !isStatusTransitionAllowed(domain.OrderStatusOnHold, domain.OrderStatusInProgress) {
		t.Fatal("expected on_hold -> in_progress to be allowed")
	}
	if isStatusTransitionAllowed(domain.OrderStatusCompleted, domain.OrderStatusInProgress) {
		t.Fatal("expected completed -> in_progress to be rejected")
	}
}
