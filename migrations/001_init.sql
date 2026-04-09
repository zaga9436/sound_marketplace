CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('customer', 'engineer', 'admin')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS profiles (
    user_id TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    display_name TEXT NOT NULL,
    bio TEXT NOT NULL DEFAULT '',
    rating DOUBLE PRECISION NOT NULL DEFAULT 0,
    reviews_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS cards (
    id TEXT PRIMARY KEY,
    author_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    card_type TEXT NOT NULL CHECK (card_type IN ('offer', 'request')),
    kind TEXT NOT NULL CHECK (kind IN ('product', 'service')),
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    price BIGINT NOT NULL CHECK (price >= 0),
    tags TEXT[] NOT NULL DEFAULT '{}',
    is_published BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS bids (
    id TEXT PRIMARY KEY,
    request_id TEXT NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    engineer_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    price BIGINT NOT NULL CHECK (price >= 0),
    message TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS orders (
    id TEXT PRIMARY KEY,
    card_id TEXT NULL REFERENCES cards(id) ON DELETE SET NULL,
    request_id TEXT NULL REFERENCES cards(id) ON DELETE SET NULL,
    bid_id TEXT NULL UNIQUE REFERENCES bids(id) ON DELETE SET NULL,
    customer_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    engineer_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount BIGINT NOT NULL CHECK (amount >= 0),
    status TEXT NOT NULL CHECK (status IN ('created', 'on_hold', 'in_progress', 'review', 'completed', 'dispute', 'cancelled')),
    delivery_notes TEXT NOT NULL DEFAULT '',
    dispute_reason TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transactions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_id TEXT NULL REFERENCES orders(id) ON DELETE SET NULL,
    type TEXT NOT NULL CHECK (type IN ('deposit', 'hold', 'release', 'refund', 'partial_refund')),
    amount BIGINT NOT NULL CHECK (amount >= 0),
    external_id TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS payments (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    external_id TEXT NOT NULL UNIQUE,
    amount BIGINT NOT NULL CHECK (amount > 0),
    status TEXT NOT NULL,
    provider TEXT NOT NULL,
    redirect_url TEXT NOT NULL,
    callback_data TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS media_files (
    id TEXT PRIMARY KEY,
    card_id TEXT NULL REFERENCES cards(id) ON DELETE CASCADE,
    order_id TEXT NULL REFERENCES orders(id) ON DELETE CASCADE,
    uploaded_by TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    storage_key TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    purpose TEXT NOT NULL,
    visibility TEXT NOT NULL DEFAULT 'private',
    is_processed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS chat_rooms (
    id TEXT PRIMARY KEY,
    order_id TEXT NOT NULL UNIQUE REFERENCES orders(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS messages (
    id TEXT PRIMARY KEY,
    chat_room_id TEXT NOT NULL REFERENCES chat_rooms(id) ON DELETE CASCADE,
    sender_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS notifications (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS reviews (
    id TEXT PRIMARY KEY,
    card_id TEXT NULL REFERENCES cards(id) ON DELETE SET NULL,
    order_id TEXT NULL REFERENCES orders(id) ON DELETE SET NULL,
    author_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rating INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_reviews_order_author_unique ON reviews(order_id, author_id) WHERE order_id IS NOT NULL;

CREATE TABLE IF NOT EXISTS disputes (
    id TEXT PRIMARY KEY,
    order_id TEXT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    opened_by TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'open',
    resolution TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMPTZ NULL
);
