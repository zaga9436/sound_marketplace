package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/soundmarket/backend/internal/domain"
)

type sqlRunner interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type PostgresStore struct {
	db *sql.DB
	tx *sql.Tx
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) WithTx(fn func(Store) error) error {
	if s.tx != nil {
		return fn(s)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	child := &PostgresStore{db: s.db, tx: tx}
	if err := fn(child); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *PostgresStore) runner() sqlRunner {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

func (s *PostgresStore) CreateUser(email, passwordHash string, role domain.Role) (domain.User, domain.Profile, error) {
	var (
		user    domain.User
		profile domain.Profile
	)

	err := s.WithTx(func(txStore Store) error {
		ps := txStore.(*PostgresStore)
		now := time.Now().UTC()
		user = domain.User{
			ID:           uuid.NewString(),
			Email:        strings.ToLower(email),
			PasswordHash: passwordHash,
			Role:         role,
			CreatedAt:    now,
		}
		profile = domain.Profile{
			UserID:       user.ID,
			DisplayName:  strings.Split(user.Email, "@")[0],
			Bio:          "",
			Rating:       0,
			ReviewsCount: 0,
			CreatedAt:    now,
		}

		if _, err := ps.runner().Exec(
			`INSERT INTO users (id, email, password_hash, role, created_at) VALUES ($1, $2, $3, $4, $5)`,
			user.ID, user.Email, user.PasswordHash, string(user.Role), user.CreatedAt,
		); err != nil {
			return err
		}
		if _, err := ps.runner().Exec(
			`INSERT INTO profiles (user_id, display_name, bio, rating, reviews_count, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $6)`,
			profile.UserID, profile.DisplayName, profile.Bio, profile.Rating, profile.ReviewsCount, profile.CreatedAt,
		); err != nil {
			return err
		}
		return nil
	})

	return user, profile, err
}

func (s *PostgresStore) FindUserByEmail(email string) (domain.User, error) {
	return s.scanUser(s.runner().QueryRow(
		`SELECT id, email, password_hash, role, created_at FROM users WHERE email = $1`,
		strings.ToLower(email),
	))
}

func (s *PostgresStore) GetUser(userID string) (domain.User, error) {
	return s.scanUser(s.runner().QueryRow(
		`SELECT id, email, password_hash, role, created_at FROM users WHERE id = $1`,
		userID,
	))
}

func (s *PostgresStore) GetProfile(userID string) (domain.Profile, error) {
	var profile domain.Profile
	err := s.runner().QueryRow(
		`SELECT user_id, display_name, bio, rating, reviews_count, created_at FROM profiles WHERE user_id = $1`,
		userID,
	).Scan(&profile.UserID, &profile.DisplayName, &profile.Bio, &profile.Rating, &profile.ReviewsCount, &profile.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Profile{}, ErrNotFound
		}
		return domain.Profile{}, err
	}
	return profile, nil
}

func (s *PostgresStore) UpdateProfile(userID, displayName, bio string) (domain.Profile, error) {
	_, err := s.runner().Exec(
		`UPDATE profiles SET display_name = $2, bio = $3, updated_at = NOW() WHERE user_id = $1`,
		userID, displayName, bio,
	)
	if err != nil {
		return domain.Profile{}, err
	}
	return s.GetProfile(userID)
}

func (s *PostgresStore) ListCardsByAuthor(authorID string) ([]domain.Card, error) {
	rows, err := s.runner().Query(
		`SELECT id, author_id, card_type, kind, title, description, price, tags, is_published, created_at
		 FROM cards
		 WHERE author_id = $1
		 ORDER BY created_at DESC`,
		authorID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cards := make([]domain.Card, 0)
	for rows.Next() {
		var card domain.Card
		if err := rows.Scan(&card.ID, &card.AuthorID, &card.CardType, &card.Kind, &card.Title, &card.Description, &card.Price, pq.Array(&card.Tags), &card.IsPublished, &card.CreatedAt); err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return cards, rows.Err()
}

func (s *PostgresStore) CreateCard(card domain.Card) (domain.Card, error) {
	if card.Tags == nil {
		card.Tags = []string{}
	}
	card.ID = uuid.NewString()
	card.CreatedAt = time.Now().UTC()
	_, err := s.runner().Exec(
		`INSERT INTO cards (id, author_id, card_type, kind, title, description, price, tags, is_published, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $10)`,
		card.ID, card.AuthorID, string(card.CardType), string(card.Kind), card.Title, card.Description, card.Price, pq.Array(card.Tags), card.IsPublished, card.CreatedAt,
	)
	if err != nil {
		return domain.Card{}, err
	}
	return card, nil
}

func (s *PostgresStore) UpdateCard(cardID string, payload domain.Card) (domain.Card, error) {
	if payload.Tags == nil {
		payload.Tags = []string{}
	}
	_, err := s.runner().Exec(
		`UPDATE cards SET kind = $2, title = $3, description = $4, price = $5, tags = $6, is_published = $7, updated_at = NOW() WHERE id = $1`,
		cardID, string(payload.Kind), payload.Title, payload.Description, payload.Price, pq.Array(payload.Tags), payload.IsPublished,
	)
	if err != nil {
		return domain.Card{}, err
	}
	return s.GetCard(cardID)
}

func (s *PostgresStore) ListCards(cardType, query string) ([]domain.Card, error) {
	base := `SELECT id, author_id, card_type, kind, title, description, price, tags, is_published, created_at FROM cards`
	args := make([]interface{}, 0, 2)
	conditions := make([]string, 0, 2)

	if cardType != "" {
		args = append(args, cardType)
		conditions = append(conditions, fmt.Sprintf("card_type = $%d", len(args)))
	}
	if query != "" {
		args = append(args, "%"+strings.ToLower(query)+"%")
		conditions = append(conditions, fmt.Sprintf("(LOWER(title) LIKE $%d OR LOWER(description) LIKE $%d)", len(args), len(args)))
	}
	if len(conditions) > 0 {
		base += " WHERE " + strings.Join(conditions, " AND ")
	}
	base += " ORDER BY created_at DESC"

	rows, err := s.runner().Query(base, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []domain.Card
	for rows.Next() {
		var card domain.Card
		if err := rows.Scan(&card.ID, &card.AuthorID, &card.CardType, &card.Kind, &card.Title, &card.Description, &card.Price, pq.Array(&card.Tags), &card.IsPublished, &card.CreatedAt); err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	if cards == nil {
		cards = make([]domain.Card, 0)
	}
	return cards, rows.Err()
}

func (s *PostgresStore) GetCard(cardID string) (domain.Card, error) {
	var card domain.Card
	err := s.runner().QueryRow(
		`SELECT id, author_id, card_type, kind, title, description, price, tags, is_published, created_at FROM cards WHERE id = $1`,
		cardID,
	).Scan(&card.ID, &card.AuthorID, &card.CardType, &card.Kind, &card.Title, &card.Description, &card.Price, pq.Array(&card.Tags), &card.IsPublished, &card.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Card{}, ErrNotFound
		}
		return domain.Card{}, err
	}
	return card, nil
}

func (s *PostgresStore) CreateMedia(media domain.MediaFile) (domain.MediaFile, error) {
	media.ID = uuid.NewString()
	media.CreatedAt = time.Now().UTC()
	visibility := "private"
	if media.MediaRole == domain.MediaRolePreview {
		visibility = "public"
	}

	_, err := s.runner().Exec(
		`INSERT INTO media_files (id, card_id, uploaded_by, storage_key, original_filename, mime_type, size_bytes, purpose, visibility, is_processed, created_at)
		 VALUES ($1, NULLIF($2, ''), $3, $4, $5, $6, $7, $8, $9, TRUE, $10)`,
		media.ID, media.CardID, media.OwnerUserID, media.FileKey, media.OriginalFilename, media.ContentType, media.SizeBytes, string(media.MediaRole), visibility, media.CreatedAt,
	)
	if err != nil {
		return domain.MediaFile{}, err
	}
	return media, nil
}

func (s *PostgresStore) ListMediaByCardAndRole(cardID string, role domain.MediaRole) ([]domain.MediaFile, error) {
	rows, err := s.runner().Query(
		`SELECT id, COALESCE(card_id, ''), uploaded_by, storage_key, COALESCE(original_filename, ''), mime_type, COALESCE(size_bytes, 0), purpose, created_at
		 FROM media_files
		 WHERE card_id = $1 AND purpose = $2
		 ORDER BY created_at ASC`,
		cardID, string(role),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mediaFiles := make([]domain.MediaFile, 0)
	for rows.Next() {
		var media domain.MediaFile
		var roleValue string
		if err := rows.Scan(&media.ID, &media.CardID, &media.OwnerUserID, &media.FileKey, &media.OriginalFilename, &media.ContentType, &media.SizeBytes, &roleValue, &media.CreatedAt); err != nil {
			return nil, err
		}
		media.MediaRole = domain.MediaRole(roleValue)
		mediaFiles = append(mediaFiles, media)
	}
	return mediaFiles, rows.Err()
}

func (s *PostgresStore) GetLatestMediaByCardAndRole(cardID string, role domain.MediaRole) (domain.MediaFile, error) {
	return s.scanMedia(s.runner().QueryRow(
		`SELECT id, COALESCE(card_id, ''), uploaded_by, storage_key, COALESCE(original_filename, ''), mime_type, COALESCE(size_bytes, 0), purpose, created_at
		 FROM media_files
		 WHERE card_id = $1 AND purpose = $2
		 ORDER BY created_at DESC
		 LIMIT 1`,
		cardID, string(role),
	))
}

func (s *PostgresStore) UserHasCompletedCardAccess(cardID, userID string) (bool, error) {
	var exists bool
	err := s.runner().QueryRow(
		`SELECT EXISTS (
			SELECT 1
			FROM orders
			WHERE status = 'completed'
			  AND (
			    (card_id = $1 AND (customer_id = $2 OR engineer_id = $2))
			    OR
			    (request_id = $1 AND (customer_id = $2 OR engineer_id = $2))
			  )
		)`,
		cardID, userID,
	).Scan(&exists)
	return exists, err
}

func (s *PostgresStore) CreateBid(bid domain.Bid) (domain.Bid, error) {
	bid.ID = uuid.NewString()
	bid.CreatedAt = time.Now().UTC()
	_, err := s.runner().Exec(
		`INSERT INTO bids (id, request_id, engineer_id, price, message, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		bid.ID, bid.RequestID, bid.EngineerID, bid.Price, bid.Message, bid.CreatedAt,
	)
	if err != nil {
		return domain.Bid{}, err
	}
	return bid, nil
}

func (s *PostgresStore) ListBidsByRequest(requestID string) ([]domain.Bid, error) {
	rows, err := s.runner().Query(
		`SELECT id, request_id, engineer_id, price, message, created_at FROM bids WHERE request_id = $1 ORDER BY created_at DESC`,
		requestID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []domain.Bid
	for rows.Next() {
		var bid domain.Bid
		if err := rows.Scan(&bid.ID, &bid.RequestID, &bid.EngineerID, &bid.Price, &bid.Message, &bid.CreatedAt); err != nil {
			return nil, err
		}
		bids = append(bids, bid)
	}
	if bids == nil {
		bids = make([]domain.Bid, 0)
	}
	return bids, rows.Err()
}

func (s *PostgresStore) ListBidsByRequestForAuthor(requestID, authorID string) ([]domain.Bid, error) {
	rows, err := s.runner().Query(
		`SELECT b.id, b.request_id, b.engineer_id, b.price, b.message, b.created_at
		 FROM bids b
		 JOIN cards c ON c.id = b.request_id
		 WHERE b.request_id = $1 AND c.author_id = $2 AND c.card_type = 'request'
		 ORDER BY b.created_at DESC`,
		requestID, authorID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []domain.Bid
	for rows.Next() {
		var bid domain.Bid
		if err := rows.Scan(&bid.ID, &bid.RequestID, &bid.EngineerID, &bid.Price, &bid.Message, &bid.CreatedAt); err != nil {
			return nil, err
		}
		bids = append(bids, bid)
	}
	if bids == nil {
		bids = make([]domain.Bid, 0)
	}
	return bids, rows.Err()
}

func (s *PostgresStore) GetBid(bidID string) (domain.Bid, error) {
	var bid domain.Bid
	err := s.runner().QueryRow(
		`SELECT id, request_id, engineer_id, price, message, created_at FROM bids WHERE id = $1`,
		bidID,
	).Scan(&bid.ID, &bid.RequestID, &bid.EngineerID, &bid.Price, &bid.Message, &bid.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Bid{}, ErrNotFound
		}
		return domain.Bid{}, err
	}
	return bid, nil
}

func (s *PostgresStore) GetBidByRequestAndEngineer(requestID, engineerID string) (domain.Bid, error) {
	var bid domain.Bid
	err := s.runner().QueryRow(
		`SELECT id, request_id, engineer_id, price, message, created_at FROM bids WHERE request_id = $1 AND engineer_id = $2 ORDER BY created_at DESC LIMIT 1`,
		requestID, engineerID,
	).Scan(&bid.ID, &bid.RequestID, &bid.EngineerID, &bid.Price, &bid.Message, &bid.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Bid{}, ErrNotFound
		}
		return domain.Bid{}, err
	}
	return bid, nil
}

func (s *PostgresStore) CreateOrder(order domain.Order) (domain.Order, error) {
	order.ID = uuid.NewString()
	now := time.Now().UTC()
	order.CreatedAt = now
	order.LastStatusTime = now
	_, err := s.runner().Exec(
		`INSERT INTO orders (id, card_id, request_id, bid_id, customer_id, engineer_id, amount, status, delivery_notes, dispute_reason, created_at, updated_at)
		 VALUES ($1, NULLIF($2, ''), NULLIF($3, ''), NULLIF($4, ''), $5, $6, $7, $8, $9, $10, $11, $11)`,
		order.ID, order.CardID, order.RequestID, order.BidID, order.CustomerID, order.EngineerID, order.Amount, string(order.Status), order.DeliveryNotes, order.DisputeReason, order.CreatedAt,
	)
	if err != nil {
		return domain.Order{}, err
	}
	if _, err := s.runner().Exec(
		`INSERT INTO chat_rooms (id, order_id, created_at) VALUES ($1, $2, $3)`,
		uuid.NewString(), order.ID, order.CreatedAt,
	); err != nil {
		return domain.Order{}, err
	}
	return order, nil
}

func (s *PostgresStore) GetOrder(orderID string) (domain.Order, error) {
	return s.scanOrder(s.runner().QueryRow(
		`SELECT id, COALESCE(card_id, ''), COALESCE(request_id, ''), COALESCE(bid_id, ''), customer_id, engineer_id, amount, status, COALESCE(delivery_notes, ''), COALESCE(dispute_reason, ''), created_at, updated_at
		 FROM orders WHERE id = $1`,
		orderID,
	))
}

func (s *PostgresStore) GetOrderByBidID(bidID string) (domain.Order, error) {
	return s.scanOrder(s.runner().QueryRow(
		`SELECT id, COALESCE(card_id, ''), COALESCE(request_id, ''), COALESCE(bid_id, ''), customer_id, engineer_id, amount, status, COALESCE(delivery_notes, ''), COALESCE(dispute_reason, ''), created_at, updated_at
		 FROM orders WHERE bid_id = $1`,
		bidID,
	))
}

func (s *PostgresStore) GetOrderByCardAndCustomer(cardID, customerID string) (domain.Order, error) {
	return s.scanOrder(s.runner().QueryRow(
		`SELECT id, COALESCE(card_id, ''), COALESCE(request_id, ''), COALESCE(bid_id, ''), customer_id, engineer_id, amount, status, COALESCE(delivery_notes, ''), COALESCE(dispute_reason, ''), created_at, updated_at
		 FROM orders WHERE card_id = $1 AND customer_id = $2 ORDER BY created_at DESC LIMIT 1`,
		cardID, customerID,
	))
}

func (s *PostgresStore) ListOrdersByCustomer(customerID string) ([]domain.Order, error) {
	return s.listOrders(
		`SELECT id, COALESCE(card_id, ''), COALESCE(request_id, ''), COALESCE(bid_id, ''), customer_id, engineer_id, amount, status, COALESCE(delivery_notes, ''), COALESCE(dispute_reason, ''), created_at, updated_at
		 FROM orders WHERE customer_id = $1 ORDER BY created_at DESC`,
		customerID,
	)
}

func (s *PostgresStore) ListOrdersByEngineer(engineerID string) ([]domain.Order, error) {
	return s.listOrders(
		`SELECT id, COALESCE(card_id, ''), COALESCE(request_id, ''), COALESCE(bid_id, ''), customer_id, engineer_id, amount, status, COALESCE(delivery_notes, ''), COALESCE(dispute_reason, ''), created_at, updated_at
		 FROM orders WHERE engineer_id = $1 ORDER BY created_at DESC`,
		engineerID,
	)
}

func (s *PostgresStore) ListOrders() ([]domain.Order, error) {
	return s.listOrders(
		`SELECT id, COALESCE(card_id, ''), COALESCE(request_id, ''), COALESCE(bid_id, ''), customer_id, engineer_id, amount, status, COALESCE(delivery_notes, ''), COALESCE(dispute_reason, ''), created_at, updated_at
		 FROM orders ORDER BY created_at DESC`,
	)
}

func (s *PostgresStore) UpdateOrder(order domain.Order) (domain.Order, error) {
	order.LastStatusTime = time.Now().UTC()
	_, err := s.runner().Exec(
		`UPDATE orders SET status = $2, delivery_notes = $3, dispute_reason = $4, updated_at = $5 WHERE id = $1`,
		order.ID, string(order.Status), order.DeliveryNotes, order.DisputeReason, order.LastStatusTime,
	)
	if err != nil {
		return domain.Order{}, err
	}
	return s.GetOrder(order.ID)
}

func (s *PostgresStore) CreateTransaction(tx domain.Transaction) (domain.Transaction, error) {
	tx.ID = uuid.NewString()
	tx.CreatedAt = time.Now().UTC()
	_, err := s.runner().Exec(
		`INSERT INTO transactions (id, user_id, order_id, type, amount, external_id, created_at)
		 VALUES ($1, $2, NULLIF($3, ''), $4, $5, NULLIF($6, ''), $7)`,
		tx.ID, tx.UserID, tx.OrderID, string(tx.Type), tx.Amount, tx.ExternalID, tx.CreatedAt,
	)
	if err != nil {
		return domain.Transaction{}, err
	}
	return tx, nil
}

func (s *PostgresStore) GetBalance(userID string) (int64, error) {
	var balance sql.NullInt64
	err := s.runner().QueryRow(
		`SELECT COALESCE(SUM(CASE
			WHEN type IN ('deposit', 'release', 'refund', 'partial_refund') THEN amount
			WHEN type = 'hold' THEN -amount
			ELSE 0
		END), 0) AS balance
		FROM transactions WHERE user_id = $1`,
		userID,
	).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance.Int64, nil
}

func (s *PostgresStore) CreatePayment(payment domain.Payment) (domain.Payment, error) {
	payment.ID = uuid.NewString()
	payment.CreatedAt = time.Now().UTC()
	_, err := s.runner().Exec(
		`INSERT INTO payments (id, user_id, external_id, amount, status, provider, redirect_url, callback_data, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $9)`,
		payment.ID, payment.UserID, payment.ExternalID, payment.Amount, payment.Status, payment.Provider, payment.RedirectURL, payment.CallbackData, payment.CreatedAt,
	)
	if err != nil {
		return domain.Payment{}, err
	}
	return payment, nil
}

func (s *PostgresStore) GetPaymentByExternalID(externalID string) (domain.Payment, error) {
	var payment domain.Payment
	err := s.runner().QueryRow(
		`SELECT id, user_id, external_id, amount, status, provider, redirect_url, COALESCE(callback_data, ''), created_at FROM payments WHERE external_id = $1`,
		externalID,
	).Scan(&payment.ID, &payment.UserID, &payment.ExternalID, &payment.Amount, &payment.Status, &payment.Provider, &payment.RedirectURL, &payment.CallbackData, &payment.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Payment{}, ErrNotFound
		}
		return domain.Payment{}, err
	}
	return payment, nil
}

func (s *PostgresStore) MarkPaymentSucceeded(externalID string) (domain.Payment, error) {
	_, err := s.runner().Exec(
		`UPDATE payments SET status = 'succeeded', updated_at = NOW() WHERE external_id = $1`,
		externalID,
	)
	if err != nil {
		return domain.Payment{}, err
	}
	return s.GetPaymentByExternalID(externalID)
}

func (s *PostgresStore) CreateDispute(dispute domain.Dispute) (domain.Dispute, error) {
	dispute.ID = uuid.NewString()
	dispute.CreatedAt = time.Now().UTC()
	_, err := s.runner().Exec(
		`INSERT INTO disputes (id, order_id, opened_by, reason, status, resolution, created_at, resolved_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, NULL)`,
		dispute.ID, dispute.OrderID, dispute.OpenedByUserID, dispute.Reason, string(dispute.Status), string(dispute.Resolution), dispute.CreatedAt,
	)
	if err != nil {
		return domain.Dispute{}, err
	}
	return dispute, nil
}

func (s *PostgresStore) GetDisputeByOrderID(orderID string) (domain.Dispute, error) {
	return s.scanDispute(s.runner().QueryRow(
		`SELECT id, order_id, opened_by, reason, status, resolution, created_at, resolved_at
		 FROM disputes WHERE order_id = $1 ORDER BY created_at DESC LIMIT 1`,
		orderID,
	))
}

func (s *PostgresStore) GetOpenDisputeByOrderID(orderID string) (domain.Dispute, error) {
	return s.scanDispute(s.runner().QueryRow(
		`SELECT id, order_id, opened_by, reason, status, resolution, created_at, resolved_at
		 FROM disputes WHERE order_id = $1 AND status = 'open' ORDER BY created_at DESC LIMIT 1`,
		orderID,
	))
}

func (s *PostgresStore) CloseDispute(disputeID string, resolution domain.DisputeResolution) (domain.Dispute, error) {
	now := time.Now().UTC()
	_, err := s.runner().Exec(
		`UPDATE disputes SET status = 'closed', resolution = $2, resolved_at = $3 WHERE id = $1`,
		disputeID, string(resolution), now,
	)
	if err != nil {
		return domain.Dispute{}, err
	}
	return s.scanDispute(s.runner().QueryRow(
		`SELECT id, order_id, opened_by, reason, status, resolution, created_at, resolved_at
		 FROM disputes WHERE id = $1`,
		disputeID,
	))
}

func (s *PostgresStore) CreateReview(review domain.Review) (domain.Review, error) {
	review.ID = uuid.NewString()
	review.CreatedAt = time.Now().UTC()
	_, err := s.runner().Exec(
		`INSERT INTO reviews (id, order_id, author_id, target_user_id, rating, comment, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		review.ID, review.OrderID, review.AuthorID, review.TargetUserID, review.Rating, review.Text, review.CreatedAt,
	)
	if err != nil {
		return domain.Review{}, err
	}
	return review, nil
}

func (s *PostgresStore) GetReviewByOrderAndAuthor(orderID, authorID string) (domain.Review, error) {
	return s.scanReview(s.runner().QueryRow(
		`SELECT id, order_id, author_id, target_user_id, rating, comment, created_at
		 FROM reviews
		 WHERE order_id = $1 AND author_id = $2
		 ORDER BY created_at DESC
		 LIMIT 1`,
		orderID, authorID,
	))
}

func (s *PostgresStore) ListReviewsByTargetUser(targetUserID string) ([]domain.Review, error) {
	rows, err := s.runner().Query(
		`SELECT id, order_id, author_id, target_user_id, rating, comment, created_at
		 FROM reviews
		 WHERE target_user_id = $1
		 ORDER BY created_at DESC`,
		targetUserID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := make([]domain.Review, 0)
	for rows.Next() {
		var review domain.Review
		if err := rows.Scan(&review.ID, &review.OrderID, &review.AuthorID, &review.TargetUserID, &review.Rating, &review.Text, &review.CreatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, rows.Err()
}

func (s *PostgresStore) RefreshProfileRating(userID string) (domain.Profile, error) {
	_, err := s.runner().Exec(
		`UPDATE profiles
		 SET rating = COALESCE(agg.avg_rating, 0),
		     reviews_count = COALESCE(agg.reviews_count, 0),
		     updated_at = NOW()
		 FROM (
		     SELECT target_user_id, AVG(rating)::double precision AS avg_rating, COUNT(*)::int AS reviews_count
		     FROM reviews
		     WHERE target_user_id = $1
		     GROUP BY target_user_id
		 ) AS agg
		 WHERE profiles.user_id = $1`,
		userID,
	)
	if err != nil {
		return domain.Profile{}, err
	}
	if _, err := s.runner().Exec(
		`UPDATE profiles
		 SET rating = 0, reviews_count = 0, updated_at = NOW()
		 WHERE user_id = $1
		   AND NOT EXISTS (SELECT 1 FROM reviews WHERE target_user_id = $1)`,
		userID,
	); err != nil {
		return domain.Profile{}, err
	}
	return s.GetProfile(userID)
}

func (s *PostgresStore) CreateNotification(userID, eventType, message string) error {
	_, err := s.runner().Exec(
		`INSERT INTO notifications (id, user_id, type, message, is_read, created_at) VALUES ($1, $2, $3, $4, FALSE, $5)`,
		uuid.NewString(), userID, eventType, message, time.Now().UTC(),
	)
	return err
}

func (s *PostgresStore) scanUser(row *sql.Row) (domain.User, error) {
	var user domain.User
	var role string
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &role, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, ErrNotFound
		}
		return domain.User{}, err
	}
	user.Role = domain.Role(role)
	return user, nil
}

func (s *PostgresStore) scanOrder(row *sql.Row) (domain.Order, error) {
	var order domain.Order
	var status string
	err := row.Scan(
		&order.ID,
		&order.CardID,
		&order.RequestID,
		&order.BidID,
		&order.CustomerID,
		&order.EngineerID,
		&order.Amount,
		&status,
		&order.DeliveryNotes,
		&order.DisputeReason,
		&order.CreatedAt,
		&order.LastStatusTime,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Order{}, ErrNotFound
		}
		return domain.Order{}, err
	}
	order.Status = domain.OrderStatus(status)
	return order, nil
}

func (s *PostgresStore) scanDispute(row *sql.Row) (domain.Dispute, error) {
	var dispute domain.Dispute
	var status string
	var resolution string
	var closedAt sql.NullTime
	err := row.Scan(
		&dispute.ID,
		&dispute.OrderID,
		&dispute.OpenedByUserID,
		&dispute.Reason,
		&status,
		&resolution,
		&dispute.CreatedAt,
		&closedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Dispute{}, ErrNotFound
		}
		return domain.Dispute{}, err
	}
	dispute.Status = domain.DisputeStatus(status)
	dispute.Resolution = domain.DisputeResolution(resolution)
	if closedAt.Valid {
		t := closedAt.Time
		dispute.ClosedAt = &t
	}
	return dispute, nil
}

func (s *PostgresStore) scanMedia(row *sql.Row) (domain.MediaFile, error) {
	var media domain.MediaFile
	var roleValue string
	err := row.Scan(
		&media.ID,
		&media.CardID,
		&media.OwnerUserID,
		&media.FileKey,
		&media.OriginalFilename,
		&media.ContentType,
		&media.SizeBytes,
		&roleValue,
		&media.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.MediaFile{}, ErrNotFound
		}
		return domain.MediaFile{}, err
	}
	media.MediaRole = domain.MediaRole(roleValue)
	return media, nil
}

func (s *PostgresStore) scanReview(row *sql.Row) (domain.Review, error) {
	var review domain.Review
	err := row.Scan(
		&review.ID,
		&review.OrderID,
		&review.AuthorID,
		&review.TargetUserID,
		&review.Rating,
		&review.Text,
		&review.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Review{}, ErrNotFound
		}
		return domain.Review{}, err
	}
	return review, nil
}

func (s *PostgresStore) listOrders(query string, args ...interface{}) ([]domain.Order, error) {
	rows, err := s.runner().Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]domain.Order, 0)
	for rows.Next() {
		var order domain.Order
		var status string
		if err := rows.Scan(
			&order.ID,
			&order.CardID,
			&order.RequestID,
			&order.BidID,
			&order.CustomerID,
			&order.EngineerID,
			&order.Amount,
			&status,
			&order.DeliveryNotes,
			&order.DisputeReason,
			&order.CreatedAt,
			&order.LastStatusTime,
		); err != nil {
			return nil, err
		}
		order.Status = domain.OrderStatus(status)
		orders = append(orders, order)
	}
	return orders, rows.Err()
}
