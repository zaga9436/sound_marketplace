package domain

import "time"

type Role string

const (
	RoleCustomer Role = "customer"
	RoleEngineer Role = "engineer"
	RoleAdmin    Role = "admin"
)

type CardType string

const (
	CardTypeOffer   CardType = "offer"
	CardTypeRequest CardType = "request"
)

type CardKind string

const (
	CardKindProduct CardKind = "product"
	CardKindService CardKind = "service"
)

type OrderStatus string

const (
	OrderStatusCreated    OrderStatus = "created"
	OrderStatusOnHold     OrderStatus = "on_hold"
	OrderStatusInProgress OrderStatus = "in_progress"
	OrderStatusReview     OrderStatus = "review"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusDispute    OrderStatus = "dispute"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

type TransactionType string

const (
	TransactionTypeDeposit       TransactionType = "deposit"
	TransactionTypeHold          TransactionType = "hold"
	TransactionTypeRelease       TransactionType = "release"
	TransactionTypeRefund        TransactionType = "refund"
	TransactionTypePartialRefund TransactionType = "partial_refund"
)

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type Profile struct {
	UserID       string    `json:"user_id"`
	DisplayName  string    `json:"display_name"`
	Bio          string    `json:"bio"`
	Rating       float64   `json:"rating"`
	ReviewsCount int       `json:"reviews_count"`
	CreatedAt    time.Time `json:"created_at"`
}

type Card struct {
	ID          string    `json:"id"`
	AuthorID    string    `json:"author_id"`
	CardType    CardType  `json:"card_type"`
	Kind        CardKind  `json:"kind"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	Tags        []string  `json:"tags"`
	PreviewURLs []string  `json:"preview_urls"`
	IsPublished bool      `json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
}

type Bid struct {
	ID         string    `json:"id"`
	RequestID  string    `json:"request_id"`
	EngineerID string    `json:"engineer_id"`
	Price      int64     `json:"price"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
}

type Order struct {
	ID            string      `json:"id"`
	CardID         string      `json:"card_id,omitempty"`
	RequestID      string      `json:"request_id,omitempty"`
	BidID          string      `json:"bid_id,omitempty"`
	CustomerID     string      `json:"customer_id"`
	EngineerID     string      `json:"engineer_id"`
	Amount         int64       `json:"amount"`
	Status         OrderStatus `json:"status"`
	DeliveryNotes  string      `json:"delivery_notes,omitempty"`
	DisputeReason  string      `json:"dispute_reason,omitempty"`
	CreatedAt      time.Time   `json:"created_at"`
	LastStatusTime time.Time   `json:"last_status_time"`
}

type Transaction struct {
	ID         string          `json:"id"`
	UserID     string          `json:"user_id"`
	OrderID    string          `json:"order_id,omitempty"`
	Type       TransactionType `json:"type"`
	Amount     int64           `json:"amount"`
	ExternalID string          `json:"external_id,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
}

type Payment struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	ExternalID      string    `json:"external_id"`
	Amount          int64     `json:"amount"`
	Status          string    `json:"status"`
	Provider        string    `json:"provider"`
	RedirectURL     string    `json:"redirect_url"`
	ConfirmationURL string    `json:"confirmation_url,omitempty"`
	CallbackData    string    `json:"callback_data,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type DisputeStatus string

const (
	DisputeStatusOpen   DisputeStatus = "open"
	DisputeStatusClosed DisputeStatus = "closed"
)

type DisputeResolution string

const (
	DisputeResolutionCompleteOrder DisputeResolution = "complete_order"
	DisputeResolutionCancelOrder   DisputeResolution = "cancel_order"
)

type Dispute struct {
	ID             string             `json:"id"`
	OrderID        string             `json:"order_id"`
	OpenedByUserID string             `json:"opened_by_user_id"`
	Reason         string             `json:"reason"`
	Status         DisputeStatus      `json:"status"`
	Resolution     DisputeResolution  `json:"resolution,omitempty"`
	CreatedAt      time.Time          `json:"created_at"`
	ClosedAt       *time.Time         `json:"closed_at,omitempty"`
}

type Review struct {
	ID           string    `json:"id"`
	OrderID      string    `json:"order_id"`
	AuthorID     string    `json:"author_id"`
	TargetUserID string    `json:"target_user_id"`
	Rating       int       `json:"rating"`
	Text         string    `json:"text"`
	CreatedAt    time.Time `json:"created_at"`
}

type MediaRole string

const (
	MediaRolePreview MediaRole = "preview"
	MediaRoleFull    MediaRole = "full"
)

type MediaFile struct {
	ID               string    `json:"id"`
	CardID           string    `json:"card_id,omitempty"`
	OwnerUserID      string    `json:"owner_user_id"`
	FileKey          string    `json:"file_key"`
	OriginalFilename string    `json:"original_filename"`
	ContentType      string    `json:"content_type"`
	SizeBytes        int64     `json:"size_bytes"`
	MediaRole        MediaRole `json:"media_role"`
	URL              string    `json:"url,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}
