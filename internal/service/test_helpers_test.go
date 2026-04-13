package service

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/soundmarket/backend/internal/domain"
	"github.com/soundmarket/backend/internal/payments"
	"github.com/soundmarket/backend/internal/repository"
	"github.com/soundmarket/backend/internal/storage"
)

type fakeStore struct {
	users              map[string]domain.User
	cards              map[string]domain.Card
	orders             map[string]domain.Order
	disputes           map[string]domain.Dispute
	payments           map[string]domain.Payment
	deliverables       map[string][]domain.Deliverable
	transactions       []domain.Transaction
	moderationActions  []domain.ModerationAction
	nextDeliverableID  int
}

func newFakeStore() *fakeStore {
	return &fakeStore{
		users:        map[string]domain.User{},
		cards:        map[string]domain.Card{},
		orders:       map[string]domain.Order{},
		disputes:     map[string]domain.Dispute{},
		payments:     map[string]domain.Payment{},
		deliverables: map[string][]domain.Deliverable{},
		transactions: []domain.Transaction{},
	}
}

func (s *fakeStore) WithTx(fn func(repository.Store) error) error { return fn(s) }

func (s *fakeStore) CreateUser(email, passwordHash string, role domain.Role) (domain.User, domain.Profile, error) {
	return domain.User{}, domain.Profile{}, errors.New("not implemented")
}
func (s *fakeStore) FindUserByEmail(email string) (domain.User, error) { return domain.User{}, repository.ErrNotFound }
func (s *fakeStore) GetUser(userID string) (domain.User, error) {
	user, ok := s.users[userID]
	if !ok {
		return domain.User{}, repository.ErrNotFound
	}
	return user, nil
}
func (s *fakeStore) ListUsers(role, status string) ([]domain.User, error) { return nil, nil }
func (s *fakeStore) SetUserSuspended(userID string, suspended bool, reason string) (domain.User, error) {
	user, err := s.GetUser(userID)
	if err != nil {
		return domain.User{}, err
	}
	user.IsSuspended = suspended
	user.SuspensionReason = reason
	s.users[userID] = user
	return user, nil
}
func (s *fakeStore) GetProfile(userID string) (domain.Profile, error) { return domain.Profile{}, repository.ErrNotFound }
func (s *fakeStore) UpdateProfile(userID, displayName, bio string) (domain.Profile, error) {
	return domain.Profile{}, errors.New("not implemented")
}
func (s *fakeStore) ListCardsByAuthor(authorID string, query domain.CardQuery) (domain.CardList, error) {
	return domain.CardList{}, nil
}
func (s *fakeStore) CreateCard(card domain.Card) (domain.Card, error) {
	card.ID = "card-created"
	s.cards[card.ID] = card
	return card, nil
}
func (s *fakeStore) UpdateCard(cardID string, payload domain.Card) (domain.Card, error) {
	return domain.Card{}, errors.New("not implemented")
}
func (s *fakeStore) ListCards(query domain.CardQuery) (domain.CardList, error) { return domain.CardList{}, nil }
func (s *fakeStore) ListCardsForAdmin(query domain.CardQuery) (domain.CardList, error) {
	return domain.CardList{}, nil
}
func (s *fakeStore) GetCard(cardID string) (domain.Card, error) {
	card, ok := s.cards[cardID]
	if !ok {
		return domain.Card{}, repository.ErrNotFound
	}
	return card, nil
}
func (s *fakeStore) SetCardHidden(cardID string, hidden bool, reason string) (domain.Card, error) {
	card, err := s.GetCard(cardID)
	if err != nil {
		return domain.Card{}, err
	}
	card.IsHidden = hidden
	card.ModerationReason = reason
	s.cards[cardID] = card
	return card, nil
}
func (s *fakeStore) CreateMedia(media domain.MediaFile) (domain.MediaFile, error) { return media, nil }
func (s *fakeStore) ListMediaByCardAndRole(cardID string, role domain.MediaRole) ([]domain.MediaFile, error) {
	return []domain.MediaFile{}, nil
}
func (s *fakeStore) GetLatestMediaByCardAndRole(cardID string, role domain.MediaRole) (domain.MediaFile, error) {
	return domain.MediaFile{}, repository.ErrNotFound
}
func (s *fakeStore) CreateDeliverable(deliverable domain.Deliverable) (domain.Deliverable, error) {
	s.nextDeliverableID++
	deliverable.ID = "deliverable-" + string(rune('0'+s.nextDeliverableID))
	if deliverable.CreatedAt.IsZero() {
		deliverable.CreatedAt = time.Now().UTC()
	}
	s.deliverables[deliverable.OrderID] = append(s.deliverables[deliverable.OrderID], deliverable)
	return deliverable, nil
}
func (s *fakeStore) ListDeliverablesByOrder(orderID string) ([]domain.Deliverable, error) {
	return append([]domain.Deliverable{}, s.deliverables[orderID]...), nil
}
func (s *fakeStore) GetDeliverable(deliverableID string) (domain.Deliverable, error) {
	for _, items := range s.deliverables {
		for _, item := range items {
			if item.ID == deliverableID {
				return item, nil
			}
		}
	}
	return domain.Deliverable{}, repository.ErrNotFound
}
func (s *fakeStore) GetLatestDeliverableByOrder(orderID string) (domain.Deliverable, error) {
	items := s.deliverables[orderID]
	if len(items) == 0 {
		return domain.Deliverable{}, repository.ErrNotFound
	}
	latest := items[0]
	for _, item := range items[1:] {
		if item.Version > latest.Version {
			latest = item
		}
	}
	return latest, nil
}
func (s *fakeStore) DeactivateDeliverablesByOrder(orderID string) error {
	items := s.deliverables[orderID]
	for i := range items {
		items[i].IsActive = false
	}
	s.deliverables[orderID] = items
	return nil
}
func (s *fakeStore) UserHasCompletedCardAccess(cardID, userID string) (bool, error) { return false, nil }
func (s *fakeStore) GetChatRoomByOrderID(orderID string) (string, error) { return "", repository.ErrNotFound }
func (s *fakeStore) CreateMessage(orderID, senderID, body string) (domain.ChatMessage, error) {
	return domain.ChatMessage{}, errors.New("not implemented")
}
func (s *fakeStore) ListMessages(orderID, userID string, limit int, beforeID string) ([]domain.ChatMessage, error) {
	return nil, nil
}
func (s *fakeStore) CountUnreadMessages(orderID, userID string) (int64, error) { return 0, nil }
func (s *fakeStore) MarkChatRead(orderID, userID string, readAt time.Time) error { return nil }
func (s *fakeStore) ListConversationsByCustomer(userID string, limit int) ([]domain.Conversation, error) {
	return nil, nil
}
func (s *fakeStore) ListConversationsByEngineer(userID string, limit int) ([]domain.Conversation, error) {
	return nil, nil
}
func (s *fakeStore) ListConversations(limit int) ([]domain.Conversation, error) { return nil, nil }
func (s *fakeStore) CreateBid(bid domain.Bid) (domain.Bid, error) { return domain.Bid{}, errors.New("not implemented") }
func (s *fakeStore) ListBidsByRequest(requestID string) ([]domain.Bid, error)   { return nil, nil }
func (s *fakeStore) ListBidsByRequestForAuthor(requestID, authorID string) ([]domain.Bid, error) {
	return nil, nil
}
func (s *fakeStore) GetBid(bidID string) (domain.Bid, error) { return domain.Bid{}, repository.ErrNotFound }
func (s *fakeStore) GetBidByRequestAndEngineer(requestID, engineerID string) (domain.Bid, error) {
	return domain.Bid{}, repository.ErrNotFound
}
func (s *fakeStore) CreateOrder(order domain.Order) (domain.Order, error) { return domain.Order{}, errors.New("not implemented") }
func (s *fakeStore) GetOrder(orderID string) (domain.Order, error) {
	order, ok := s.orders[orderID]
	if !ok {
		return domain.Order{}, repository.ErrNotFound
	}
	return order, nil
}
func (s *fakeStore) GetOrderByBidID(bidID string) (domain.Order, error) { return domain.Order{}, repository.ErrNotFound }
func (s *fakeStore) GetOrderByCardAndCustomer(cardID, customerID string) (domain.Order, error) {
	return domain.Order{}, repository.ErrNotFound
}
func (s *fakeStore) ListOrdersByCustomer(customerID string) ([]domain.Order, error) { return nil, nil }
func (s *fakeStore) ListOrdersByEngineer(engineerID string) ([]domain.Order, error)  { return nil, nil }
func (s *fakeStore) ListOrders() ([]domain.Order, error)                              { return nil, nil }
func (s *fakeStore) UpdateOrder(order domain.Order) (domain.Order, error) {
	s.orders[order.ID] = order
	return order, nil
}
func (s *fakeStore) CreateTransaction(tx domain.Transaction) (domain.Transaction, error) {
	s.transactions = append(s.transactions, tx)
	return tx, nil
}
func (s *fakeStore) GetBalance(userID string) (int64, error) { return 0, nil }
func (s *fakeStore) CreatePayment(payment domain.Payment) (domain.Payment, error) {
	s.payments[payment.ExternalID] = payment
	return payment, nil
}
func (s *fakeStore) GetPaymentByExternalID(externalID string) (domain.Payment, error) {
	p, ok := s.payments[externalID]
	if !ok {
		return domain.Payment{}, repository.ErrNotFound
	}
	return p, nil
}
func (s *fakeStore) MarkPaymentSucceeded(externalID string) (domain.Payment, error) {
	p, err := s.GetPaymentByExternalID(externalID)
	if err != nil {
		return domain.Payment{}, err
	}
	p.Status = "succeeded"
	s.payments[externalID] = p
	return p, nil
}
func (s *fakeStore) CreateDispute(dispute domain.Dispute) (domain.Dispute, error) {
	if dispute.ID == "" {
		dispute.ID = "dispute-created"
	}
	s.disputes[dispute.OrderID] = dispute
	return dispute, nil
}
func (s *fakeStore) GetDisputeByOrderID(orderID string) (domain.Dispute, error) {
	dispute, ok := s.disputes[orderID]
	if !ok {
		return domain.Dispute{}, repository.ErrNotFound
	}
	return dispute, nil
}
func (s *fakeStore) GetOpenDisputeByOrderID(orderID string) (domain.Dispute, error) {
	dispute, err := s.GetDisputeByOrderID(orderID)
	if err != nil {
		return domain.Dispute{}, err
	}
	if dispute.Status != domain.DisputeStatusOpen {
		return domain.Dispute{}, repository.ErrNotFound
	}
	return dispute, nil
}
func (s *fakeStore) GetDispute(disputeID string) (domain.Dispute, error) {
	for _, dispute := range s.disputes {
		if dispute.ID == disputeID {
			return dispute, nil
		}
	}
	return domain.Dispute{}, repository.ErrNotFound
}
func (s *fakeStore) ListDisputes(status string) ([]domain.Dispute, error) { return nil, nil }
func (s *fakeStore) CloseDispute(disputeID string, resolution domain.DisputeResolution) (domain.Dispute, error) {
	dispute, err := s.GetDispute(disputeID)
	if err != nil {
		return domain.Dispute{}, err
	}
	now := time.Now().UTC()
	dispute.Status = domain.DisputeStatusClosed
	dispute.Resolution = resolution
	dispute.ClosedAt = &now
	s.disputes[dispute.OrderID] = dispute
	return dispute, nil
}
func (s *fakeStore) CreateReview(review domain.Review) (domain.Review, error) { return domain.Review{}, errors.New("not implemented") }
func (s *fakeStore) GetReviewByOrderAndAuthor(orderID, authorID string) (domain.Review, error) {
	return domain.Review{}, repository.ErrNotFound
}
func (s *fakeStore) ListReviewsByTargetUser(targetUserID string) ([]domain.Review, error) { return nil, nil }
func (s *fakeStore) RefreshProfileRating(userID string) (domain.Profile, error) {
	return domain.Profile{}, errors.New("not implemented")
}
func (s *fakeStore) CreateNotification(userID, eventType, message string) (domain.Notification, error) {
	return domain.Notification{}, nil
}
func (s *fakeStore) ListNotifications(userID string, limit int, beforeID string) ([]domain.Notification, error) {
	return nil, nil
}
func (s *fakeStore) MarkNotificationsRead(userID string, ids []string) error { return nil }
func (s *fakeStore) CountUnreadNotifications(userID string) (int64, error)    { return 0, nil }
func (s *fakeStore) CreateModerationAction(action domain.ModerationAction) (domain.ModerationAction, error) {
	s.moderationActions = append(s.moderationActions, action)
	return action, nil
}
func (s *fakeStore) ListModerationActions(targetType, targetID string, limit int) ([]domain.ModerationAction, error) {
	return s.moderationActions, nil
}

type fakeNotifier struct {
	events []string
}

func (n *fakeNotifier) Publish(userID, eventType, message string) {
	n.events = append(n.events, userID+":"+eventType+":"+message)
}

type fakeProvider struct {
	info *payments.PaymentInfo
}

func (p *fakeProvider) CreatePayment(ctx context.Context, input payments.CreatePaymentInput) (*payments.PaymentSession, error) {
	return &payments.PaymentSession{}, nil
}
func (p *fakeProvider) GetPayment(ctx context.Context, externalID string) (*payments.PaymentInfo, error) {
	return p.info, nil
}
func (p *fakeProvider) HandleWebhook(ctx context.Context, payload []byte, headers http.Header) (*payments.WebhookResult, error) {
	return nil, payments.ErrInvalidWebhook
}

type fakeStorage struct {
	uploads    []string
	deletes    []string
	signedURLs map[string]string
}

func (s *fakeStorage) Upload(ctx context.Context, key, contentType string, body io.Reader, opts storage.UploadOptions) (storage.StoredObject, error) {
	s.uploads = append(s.uploads, key)
	return storage.StoredObject{Key: key}, nil
}
func (s *fakeStorage) Delete(ctx context.Context, key string) error {
	s.deletes = append(s.deletes, key)
	return nil
}
func (s *fakeStorage) GenerateSignedURL(ctx context.Context, key string, ttl time.Duration) (string, error) {
	if s.signedURLs != nil {
		if url, ok := s.signedURLs[key]; ok {
			return url, nil
		}
	}
	return "signed://" + key, nil
}
func (s *fakeStorage) PublicURL(key string) string { return "public://" + key }

