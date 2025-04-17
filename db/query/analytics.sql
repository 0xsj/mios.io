-- name: RecordLinkClick :one
INSERT INTO analytics (
    link_id, user_id, ip_address, user_agent, referrer
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: RecordSocialLinkClick :one
INSERT INTO analytics (
    social_link_id, user_id, ip_address, user_agent, referrer
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetLinkClicksByUserID :many
SELECT 
    l.link_id,
    l.title,
    l.url,
    COUNT(a.analytics_id) as click_count
FROM links l
LEFT JOIN analytics a ON l.link_id = a.link_id
WHERE l.user_id = $1
GROUP BY l.link_id, l.title, l.url
ORDER BY click_count DESC;

-- name: GetSocialLinkClicksByUserID :many
SELECT 
    sl.social_link_id,
    sl.platform,
    sl.username,
    COUNT(a.analytics_id) as click_count
FROM social_links sl
LEFT JOIN analytics a ON sl.social_link_id = a.social_link_id
WHERE sl.user_id = $1
GROUP BY sl.social_link_id, sl.platform, sl.username
ORDER BY click_count DESC;

-- name: GetClickAnalyticsByDateRange :many
SELECT 
    DATE(a.clicked_at) as date,
    COUNT(a.analytics_id) as click_count
FROM analytics a
WHERE a.user_id = $1 
AND a.clicked_at >= $2
AND a.clicked_at <= $3
GROUP BY DATE(a.clicked_at)
ORDER BY DATE(a.clicked_at);

-- name: GetClickAnalyticsByReferrer :many
SELECT 
    COALESCE(a.referrer, 'direct') as referrer,
    COUNT(a.analytics_id) as click_count
FROM analytics a
WHERE a.user_id = $1
GROUP BY COALESCE(a.referrer, 'direct')
ORDER BY click_count DESC;

-- name: DeleteUserAnalytics :exec
DELETE FROM analytics
WHERE user_id = $1;