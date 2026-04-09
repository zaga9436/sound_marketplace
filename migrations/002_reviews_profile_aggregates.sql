ALTER TABLE profiles
ADD COLUMN IF NOT EXISTS reviews_count INT NOT NULL DEFAULT 0;

CREATE UNIQUE INDEX IF NOT EXISTS idx_reviews_order_author_unique
ON reviews(order_id, author_id)
WHERE order_id IS NOT NULL;

UPDATE profiles p
SET rating = COALESCE(agg.avg_rating, 0),
    reviews_count = COALESCE(agg.reviews_count, 0),
    updated_at = NOW()
FROM (
    SELECT target_user_id, AVG(rating)::double precision AS avg_rating, COUNT(*)::int AS reviews_count
    FROM reviews
    GROUP BY target_user_id
) AS agg
WHERE p.user_id = agg.target_user_id;

UPDATE profiles
SET rating = 0,
    reviews_count = 0,
    updated_at = NOW()
WHERE user_id NOT IN (SELECT DISTINCT target_user_id FROM reviews);
