export type Role = "customer" | "engineer" | "admin";

export interface User {
  id: string;
  email: string;
  role: Role;
  is_suspended?: boolean;
  suspension_reason?: string;
  suspended_at?: string;
}

export interface Profile {
  user_id: string;
  display_name: string;
  bio: string;
  rating: number;
  reviews_count: number;
  created_at?: string;
}

export type CardType = "offer" | "request";
export type CardKind = "product" | "service";
export type OrderStatus = "created" | "on_hold" | "in_progress" | "review" | "completed" | "dispute" | "cancelled";
export type DisputeStatus = "open" | "closed";
export type DisputeResolution = "complete_order" | "cancel_order";

export interface Card {
  id: string;
  author_id: string;
  card_type: CardType;
  kind: CardKind;
  title: string;
  description: string;
  price: number;
  tags: string[];
  cover_url?: string;
  preview_urls: string[];
  is_published: boolean;
  is_hidden?: boolean;
  moderation_reason?: string;
  created_at: string;
}

export interface CardListResponse {
  items: Card[];
  total: number;
  limit: number;
  offset: number;
}

export interface Bid {
  id: string;
  request_id: string;
  engineer_id: string;
  price: number;
  message: string;
  created_at: string;
}

export interface Order {
  id: string;
  card_id?: string;
  request_id?: string;
  bid_id?: string;
  customer_id: string;
  engineer_id: string;
  amount: number;
  status: OrderStatus;
  delivery_notes?: string;
  dispute_reason?: string;
  created_at: string;
  last_status_time: string;
}

export interface Conversation {
  order_id: string;
  chat_room_id: string;
  customer_id: string;
  engineer_id: string;
  last_message?: string;
  last_message_at?: string;
  unread_count: number;
}

export interface ChatMessage {
  id: string;
  chat_room_id: string;
  order_id: string;
  sender_id: string;
  body: string;
  created_at: string;
  read_at?: string;
}

export interface AuthResponse {
  token?: string;
  user: User;
  profile: Profile;
}

export interface Review {
  id: string;
  order_id: string;
  author_id: string;
  target_user_id: string;
  rating: number;
  text: string;
  created_at: string;
}

export interface Dispute {
  id: string;
  order_id: string;
  opened_by_user_id: string;
  reason: string;
  status: DisputeStatus;
  resolution?: DisputeResolution;
  created_at: string;
  closed_at?: string;
}

export interface ProfileUpdatePayload {
  display_name: string;
  bio: string;
}

export interface CardPayload {
  card_type: CardType;
  kind: CardKind;
  title: string;
  description: string;
  price: number;
  tags: string[];
  is_published: boolean;
}

export interface MediaFile {
  id: string;
  card_id?: string;
  owner_user_id: string;
  file_key: string;
  original_filename: string;
  content_type: string;
  size_bytes: number;
  media_role: "cover" | "preview" | "full";
  url?: string;
  created_at?: string;
}

export interface DownloadResponse {
  url: string;
}

export interface Deliverable {
  id: string;
  order_id: string;
  uploaded_by: string;
  storage_key: string;
  original_filename: string;
  content_type: string;
  size_bytes: number;
  version: number;
  is_active: boolean;
  created_at: string;
}

export interface BidCreatePayload {
  price: number;
  message: string;
}

export interface NotificationItem {
  id: string;
  user_id: string;
  type: string;
  message: string;
  is_read: boolean;
  created_at: string;
}

export interface NotificationListResponse {
  items: NotificationItem[];
  unread_count: number;
}

export interface Payment {
  id: string;
  user_id: string;
  external_id: string;
  amount: number;
  status: string;
  provider: string;
  redirect_url: string;
  confirmation_url?: string;
  callback_data?: string;
  created_at: string;
}

export interface BalanceResponse {
  balance: number;
}

export interface PaymentSyncResponse {
  payment: Payment;
  deposit_created: boolean;
}

export interface ApiErrorPayload {
  error: string;
}
