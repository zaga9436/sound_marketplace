ALTER TABLE users
ADD COLUMN IF NOT EXISTS is_suspended BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE users
ADD COLUMN IF NOT EXISTS suspension_reason TEXT NOT NULL DEFAULT '';

ALTER TABLE users
ADD COLUMN IF NOT EXISTS suspended_at TIMESTAMPTZ NULL;

ALTER TABLE cards
ADD COLUMN IF NOT EXISTS is_hidden BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE cards
ADD COLUMN IF NOT EXISTS moderation_reason TEXT NOT NULL DEFAULT '';

CREATE TABLE IF NOT EXISTS moderation_actions (
    id TEXT PRIMARY KEY,
    admin_user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_type TEXT NOT NULL,
    target_id TEXT NOT NULL,
    action TEXT NOT NULL,
    reason TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_role_suspended ON users(role, is_suspended);
CREATE INDEX IF NOT EXISTS idx_cards_hidden_created_at ON cards(is_hidden, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_disputes_status_created_at ON disputes(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_moderation_actions_target ON moderation_actions(target_type, target_id, created_at DESC);
