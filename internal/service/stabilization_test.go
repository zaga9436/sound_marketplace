package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/soundmarket/backend/internal/apierr"
	"github.com/soundmarket/backend/internal/config"
	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/payments"
)

func TestOrderLifecycleReleaseOnComplete(t *testing.T) {
	store := newFakeStore()
	notifier := &fakeNotifier{}
	service := NewOrderService(store, notifier)
	store.orders["order-1"] = domain.Order{
		ID:         "order-1",
		CustomerID: "customer-1",
		EngineerID: "engineer-1",
		Amount:     2500,
		Status:     domain.OrderStatusReview,
	}

	order, err := service.UpdateStatus(domain.User{ID: "customer-1", Role: domain.RoleCustomer}, "order-1", domain.OrderStatusCompleted)
	if err != nil {
		t.Fatalf("expected complete status update to succeed: %v", err)
	}
	if order.Status != domain.OrderStatusCompleted {
		t.Fatalf("expected completed status, got %s", order.Status)
	}
	if len(store.transactions) != 1 || store.transactions[0].Type != domain.TransactionTypeRelease {
		t.Fatalf("expected one release transaction, got %#v", store.transactions)
	}
}

func TestDisputeCancelFlowRefundsCustomer(t *testing.T) {
	store := newFakeStore()
	notifier := &fakeNotifier{}
	service := NewDisputeService(store, notifier)
	store.users["customer-1"] = domain.User{ID: "customer-1", Role: domain.RoleCustomer}
	store.orders["order-1"] = domain.Order{
		ID:         "order-1",
		CustomerID: "customer-1",
		EngineerID: "engineer-1",
		Amount:     3000,
		Status:     domain.OrderStatusOnHold,
	}

	dispute, err := service.Open(domain.User{ID: "customer-1", Role: domain.RoleCustomer}, "order-1", "Need resolution")
	if err != nil {
		t.Fatalf("expected dispute open to succeed: %v", err)
	}
	if dispute.Status != domain.DisputeStatusOpen {
		t.Fatalf("expected open dispute, got %s", dispute.Status)
	}
	closed, err := service.Close(domain.User{ID: "customer-1", Role: domain.RoleCustomer}, "order-1", domain.DisputeResolutionCancelOrder)
	if err != nil {
		t.Fatalf("expected dispute close to succeed: %v", err)
	}
	if closed.Status != domain.DisputeStatusClosed || closed.Resolution != domain.DisputeResolutionCancelOrder {
		t.Fatalf("unexpected closed dispute: %#v", closed)
	}
	if len(store.transactions) != 1 || store.transactions[0].Type != domain.TransactionTypeRefund {
		t.Fatalf("expected refund transaction, got %#v", store.transactions)
	}
}

func TestPaymentSyncIsIdempotent(t *testing.T) {
	store := newFakeStore()
	notifier := &fakeNotifier{}
	store.payments["pay-1"] = domain.Payment{
		ID:         "payment-1",
		UserID:     "customer-1",
		ExternalID: "pay-1",
		Amount:     100,
		Status:     "pending",
		Provider:   "yookassa",
	}
	service := NewPaymentService(&config.Config{AppEnv: "development"}, store, &fakeProvider{
		info: &payments.PaymentInfo{
			ExternalID: "pay-1",
			Status:     "succeeded",
			Paid:       true,
			Provider:   "yookassa",
		},
	}, notifier)

	first, err := service.SyncPayment(context.Background(), domain.User{ID: "customer-1", Role: domain.RoleCustomer}, "pay-1")
	if err != nil {
		t.Fatalf("expected first sync to succeed: %v", err)
	}
	if !first.DepositCreated {
		t.Fatal("expected first sync to create deposit")
	}
	second, err := service.SyncPayment(context.Background(), domain.User{ID: "customer-1", Role: domain.RoleCustomer}, "pay-1")
	if err != nil {
		t.Fatalf("expected second sync to succeed: %v", err)
	}
	if second.DepositCreated {
		t.Fatal("expected second sync to be idempotent")
	}
	if len(store.transactions) != 1 {
		t.Fatalf("expected exactly one deposit transaction, got %d", len(store.transactions))
	}
}

func TestDeliverablesVersioningAndAccess(t *testing.T) {
	store := newFakeStore()
	notifier := &fakeNotifier{}
	storageAdapter := &fakeStorage{signedURLs: map[string]string{}}
	service := NewDeliverableService(&config.Config{
		MaxUploadSize:       1024 * 1024,
		AllowedAudioFormats: []string{"mp3"},
		SignedURLTTL:        time.Minute,
	}, store, storageAdapter, notifier)

	store.users["engineer-1"] = domain.User{ID: "engineer-1", Role: domain.RoleEngineer}
	store.orders["order-1"] = domain.Order{
		ID:         "order-1",
		CustomerID: "customer-1",
		EngineerID: "engineer-1",
		Status:     domain.OrderStatusInProgress,
	}

	first, err := service.Upload(context.Background(), domain.User{ID: "engineer-1", Role: domain.RoleEngineer}, "order-1", MediaUploadInput{
		Filename:    "mix.mp3",
		ContentType: "audio/mpeg",
		SizeBytes:   100,
		Reader:      strings.NewReader("v1"),
	})
	if err != nil {
		t.Fatalf("expected first upload to succeed: %v", err)
	}
	second, err := service.Upload(context.Background(), domain.User{ID: "engineer-1", Role: domain.RoleEngineer}, "order-1", MediaUploadInput{
		Filename:    "mix-v2.mp3",
		ContentType: "audio/mpeg",
		SizeBytes:   120,
		Reader:      strings.NewReader("v2"),
	})
	if err != nil {
		t.Fatalf("expected second upload to succeed: %v", err)
	}
	items, err := service.List(domain.User{ID: "customer-1", Role: domain.RoleCustomer}, "order-1")
	if err != nil {
		t.Fatalf("expected list to succeed: %v", err)
	}
	if first.Version != 1 || second.Version != 2 {
		t.Fatalf("expected versions 1 and 2, got %d and %d", first.Version, second.Version)
	}
	if len(items) != 2 || items[0].Version != 1 || items[0].IsActive || !items[1].IsActive {
		t.Fatalf("unexpected deliverable history: %#v", items)
	}
	url, err := service.Download(context.Background(), domain.User{ID: "customer-1", Role: domain.RoleCustomer}, "order-1", items[1].ID)
	if err != nil || !strings.HasPrefix(url, "signed://") {
		t.Fatalf("expected signed download, got %q, err=%v", url, err)
	}
	_, err = service.Download(context.Background(), domain.User{ID: "intruder", Role: domain.RoleCustomer}, "order-1", items[1].ID)
	if !errors.Is(err, apierr.ErrForbidden) {
		t.Fatalf("expected forbidden for outsider, got %v", err)
	}
}

func TestAdminSuspendAndHideEnforcement(t *testing.T) {
	store := newFakeStore()
	notifier := &fakeNotifier{}
	cardService := NewCardService(store, notifier, &fakeStorage{})
	store.users["engineer-1"] = domain.User{ID: "engineer-1", Role: domain.RoleEngineer, IsSuspended: true}
	if _, err := cardService.Create(domain.User{ID: "engineer-1", Role: domain.RoleEngineer}, domain.Card{
		CardType:    domain.CardTypeOffer,
		Kind:        domain.CardKindService,
		Title:       "Offer",
		Description: "Desc",
		Price:       100,
	}); !errors.Is(err, apierr.ErrForbidden) {
		t.Fatalf("expected suspended user to be blocked, got %v", err)
	}

	store.cards["hidden-card"] = domain.Card{
		ID:       "hidden-card",
		AuthorID: "engineer-1",
		CardType: domain.CardTypeOffer,
		Kind:     domain.CardKindService,
		Title:    "Hidden",
		IsHidden: true,
	}
	if _, err := cardService.Get("hidden-card"); !errors.Is(err, apierr.ErrNotFound) {
		t.Fatalf("expected hidden card to be invisible publicly, got %v", err)
	}
}
