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
	ID               string     `json:"id"`
	Email            string     `json:"email"`
	PasswordHash     string     `json:"-"`
	Role             Role       `json:"role"`
	IsSuspended      bool       `json:"is_suspended"`
	SuspensionReason string     `json:"suspension_reason,omitempty"`
	SuspendedAt      *time.Time `json:"suspended_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

type Profile struct {
	UserID       string    `json:"user_id"`
	DisplayName  string    `json:"display_name"`
	Bio          string    `json:"bio"`
	AvatarURL    string    `json:"avatar_url,omitempty"`
	Rating       float64   `json:"rating"`
	ReviewsCount int       `json:"reviews_count"`
	CreatedAt    time.Time `json:"created_at"`
}

type Card struct {
	ID               string    `json:"id"`
	AuthorID         string    `json:"author_id"`
	CardType         CardType  `json:"card_type"`
	Kind             CardKind  `json:"kind"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Price            int64     `json:"price"`
	Tags             []string  `json:"tags"`
	CoverURL         string    `json:"cover_url,omitempty"`
	PreviewURLs      []string  `json:"preview_urls"`
	IsPublished      bool      `json:"is_published"`
	IsHidden         bool      `json:"is_hidden"`
	ModerationReason string    `json:"moderation_reason,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

type CardQuery struct {
	CardType    CardType `json:"card_type,omitempty"`
	Kind        CardKind `json:"kind,omitempty"`
	AuthorID    string   `json:"author_id,omitempty"`
	Query       string   `json:"q,omitempty"`
	Tag         string   `json:"tag,omitempty"`
	MinPrice    *int64   `json:"min_price,omitempty"`
	MaxPrice    *int64   `json:"max_price,omitempty"`
	IsPublished *bool    `json:"is_published,omitempty"`
	Visibility  string   `json:"visibility,omitempty"`
	SortBy      string   `json:"sort_by,omitempty"`
	SortOrder   string   `json:"sort_order,omitempty"`
	Limit       int      `json:"limit,omitempty"`
	Offset      int      `json:"offset,omitempty"`
}

type CardList struct {
	Items  []Card `json:"items"`
	Total  int64  `json:"total"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
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
	ID             string      `json:"id"`
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
	ID             string            `json:"id"`
	OrderID        string            `json:"order_id"`
	OpenedByUserID string            `json:"opened_by_user_id"`
	Reason         string            `json:"reason"`
	Status         DisputeStatus     `json:"status"`
	Resolution     DisputeResolution `json:"resolution,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
	ClosedAt       *time.Time        `json:"closed_at,omitempty"`
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
	MediaRoleAvatar   MediaRole = "avatar"
	MediaRoleCover    MediaRole = "cover"
	MediaRolePreview  MediaRole = "preview"
	MediaRoleFull     MediaRole = "full"
	MediaRoleMaterial MediaRole = "material"
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

type Deliverable struct {
	ID               string    `json:"id"`
	OrderID          string    `json:"order_id"`
	UploadedBy       string    `json:"uploaded_by"`
	StorageKey       string    `json:"storage_key"`
	OriginalFilename string    `json:"original_filename"`
	ContentType      string    `json:"content_type"`
	SizeBytes        int64     `json:"size_bytes"`
	Version          int       `json:"version"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
}

type ChatMessage struct {
	ID         string     `json:"id"`
	ChatRoomID string     `json:"chat_room_id"`
	OrderID    string     `json:"order_id"`
	SenderID   string     `json:"sender_id"`
	Body       string     `json:"body"`
	CreatedAt  time.Time  `json:"created_at"`
	ReadAt     *time.Time `json:"read_at,omitempty"`
}

type Conversation struct {
	OrderID       string     `json:"order_id"`
	ChatRoomID    string     `json:"chat_room_id"`
	CustomerID    string     `json:"customer_id"`
	EngineerID    string     `json:"engineer_id"`
	LastMessage   string     `json:"last_message,omitempty"`
	LastMessageAt *time.Time `json:"last_message_at,omitempty"`
	UnreadCount   int64      `json:"unread_count"`
}

type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

type ModerationAction struct {
	ID          string    `json:"id"`
	AdminUserID string    `json:"admin_user_id"`
	TargetType  string    `json:"target_type"`
	TargetID    string    `json:"target_id"`
	Action      string    `json:"action"`
	Reason      string    `json:"reason,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
