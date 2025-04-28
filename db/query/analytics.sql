-- db/query/analytics.sql

-- Recording clicks and page views
-- name: CreateAnalyticsEntry :one
INSERT INTO analytics (
    item_id, user_id, ip_address, user_agent, referrer, page_view
) VALUES (
    $1, $2, $3, $4, $5, false
) RETURNING *;

-- name: CreatePageViewEntry :one
INSERT INTO analytics (
    item_id, user_id, ip_address, user_agent, referrer, page_view
) VALUES (
    $1, $2, $3, $4, $5, true
) RETURNING *;

-- Basic analytics queries
-- name: GetItemAnalytics :many
SELECT * FROM analytics
WHERE item_id = $1
ORDER BY clicked_at DESC
LIMIT $2 OFFSET $3;

-- name: GetUserAnalytics :many
SELECT * FROM analytics
WHERE user_id = $1
ORDER BY clicked_at DESC
LIMIT $2 OFFSET $3;

-- Count queries
-- name: GetContentItemClickCount :one
SELECT COUNT(*) FROM analytics
WHERE item_id = $1 AND page_view = false;

-- name: GetUserItemClickCount :one
SELECT COUNT(*) FROM analytics
WHERE user_id = $1 AND page_view = false;

-- name: GetProfilePageViews :one
SELECT COUNT(*) FROM analytics
WHERE user_id = $1 AND page_view = true;

-- Time range analytics
-- name: GetUserAnalyticsByTimeRange :many
SELECT 
    DATE_TRUNC('day', clicked_at) AS day,
    COUNT(*) AS clicks
FROM analytics
WHERE user_id = $1
AND clicked_at >= $2
AND clicked_at <= $3
AND page_view = false
GROUP BY DATE_TRUNC('day', clicked_at)
ORDER BY day;

-- name: GetItemAnalyticsByTimeRange :many
SELECT 
    DATE_TRUNC('day', clicked_at) AS day,
    COUNT(*) AS clicks
FROM analytics
WHERE item_id = $1
AND clicked_at >= $2
AND clicked_at <= $3
AND page_view = false
GROUP BY DATE_TRUNC('day', clicked_at)
ORDER BY day;

-- name: GetProfilePageViewsByDate :many
SELECT 
    DATE_TRUNC('day', clicked_at) AS day,
    COUNT(*) AS views
FROM analytics
WHERE user_id = $1
AND clicked_at >= $2
AND clicked_at <= $3
AND page_view = true
GROUP BY DATE_TRUNC('day', clicked_at)
ORDER BY day;

-- Insight queries
-- name: GetTopContentItemsByClicks :many
SELECT 
    a.item_id,
    c.content_type,
    COALESCE(c.title, '') AS title,
    COUNT(*) AS click_count
FROM analytics a
JOIN content_items c ON a.item_id = c.item_id
WHERE c.user_id = $1
AND a.clicked_at >= $2
AND a.clicked_at <= $3
AND a.page_view = false
GROUP BY a.item_id, c.content_type, c.title
ORDER BY click_count DESC
LIMIT $4;

-- name: GetReferrerAnalytics :many
SELECT 
    COALESCE(referrer, '') AS referrer,
    COUNT(*) AS count
FROM analytics
WHERE user_id = $1
AND clicked_at >= $2
AND clicked_at <= $3
AND referrer IS NOT NULL
GROUP BY referrer
ORDER BY count DESC
LIMIT $4;

-- Visitor analytics
-- name: GetUniqueVisitors :one
SELECT COUNT(DISTINCT ip_address) 
FROM analytics
WHERE user_id = $1
AND clicked_at >= $2
AND clicked_at <= $3
AND ip_address IS NOT NULL
AND page_view = true;

-- name: GetUniqueVisitorsByDay :many
SELECT 
    DATE_TRUNC('day', clicked_at) AS day,
    COUNT(DISTINCT ip_address) AS visitors
FROM analytics
WHERE user_id = $1
AND clicked_at >= $2
AND clicked_at <= $3
AND ip_address IS NOT NULL
AND page_view = true
GROUP BY DATE_TRUNC('day', clicked_at)
ORDER BY day;