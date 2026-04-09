package repository

import "github.com/soundmarket/backend/internal/domain"

type Store interface {
	WithTx(fn func(Store) error) error

	CreateUser(email, passwordHash string, role domain.Role) (domain.User, domain.Profile, error)
	FindUserByEmail(email string) (domain.User, error)
	GetUser(userID string) (domain.User, error)
	GetProfile(userID string) (domain.Profile, error)
	UpdateProfile(userID, displayName, bio string) (domain.Profile, error)
	ListCardsByAuthor(authorID string) ([]domain.Card, error)

	CreateCard(card domain.Card) (domain.Card, error)
	UpdateCard(cardID string, payload domain.Card) (domain.Card, error)
	ListCards(cardType, query string) ([]domain.Card, error)
	GetCard(cardID string) (domain.Card, error)

	CreateBid(bid domain.Bid) (domain.Bid, error)
	ListBidsByRequest(requestID string) ([]domain.Bid, error)
	ListBidsByRequestForAuthor(requestID, authorID string) ([]domain.Bid, error)
	GetBid(bidID string) (domain.Bid, error)
	GetBidByRequestAndEngineer(requestID, engineerID string) (domain.Bid, error)

	CreateOrder(order domain.Order) (domain.Order, error)
	GetOrder(orderID string) (domain.Order, error)
	GetOrderByBidID(bidID string) (domain.Order, error)
	GetOrderByCardAndCustomer(cardID, customerID string) (domain.Order, error)
	ListOrdersByCustomer(customerID string) ([]domain.Order, error)
	ListOrdersByEngineer(engineerID string) ([]domain.Order, error)
	ListOrders() ([]domain.Order, error)
	UpdateOrder(order domain.Order) (domain.Order, error)

	CreateTransaction(tx domain.Transaction) (domain.Transaction, error)
	GetBalance(userID string) (int64, error)

	CreatePayment(payment domain.Payment) (domain.Payment, error)
	GetPaymentByExternalID(externalID string) (domain.Payment, error)
	MarkPaymentSucceeded(externalID string) (domain.Payment, error)

	CreateDispute(dispute domain.Dispute) (domain.Dispute, error)
	GetDisputeByOrderID(orderID string) (domain.Dispute, error)
	GetOpenDisputeByOrderID(orderID string) (domain.Dispute, error)
	CloseDispute(disputeID string, resolution domain.DisputeResolution) (domain.Dispute, error)

	CreateReview(review domain.Review) (domain.Review, error)
	GetReviewByOrderAndAuthor(orderID, authorID string) (domain.Review, error)
	ListReviewsByTargetUser(targetUserID string) ([]domain.Review, error)
	RefreshProfileRating(userID string) (domain.Profile, error)

	CreateNotification(userID, eventType, message string) error
}
