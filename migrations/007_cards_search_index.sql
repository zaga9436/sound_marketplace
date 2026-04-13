CREATE INDEX IF NOT EXISTS idx_cards_search_tsv ON cards
USING GIN (to_tsvector('simple'::regconfig, COALESCE(title, '') || ' ' || COALESCE(description, '')));
