package repository

import (
	"context"
	"time"

	db "github.com/0xsj/gin-sqlc/db/sqlc"
	"github.com/0xsj/gin-sqlc/log"
	"github.com/0xsj/gin-sqlc/pkg/errors"
	"github.com/google/uuid"
)

type AnalyticsRepository interface {
	// Recording data
	CreateAnalyticsEntry(ctx context.Context, params CreateAnalyticsParams) (*db.Analytic, error)
	CreatePageViewEntry(ctx context.Context, params CreatePageViewParams) (*db.Analytic, error)

	// Basic analytics
	GetItemAnalytics(ctx context.Context, itemID uuid.UUID, limit, offset int) ([]*db.Analytic, error)
	GetUserAnalytics(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*db.Analytic, error)

	// Count queries
	GetContentItemClickCount(ctx context.Context, itemID uuid.UUID) (int64, error)
	GetUserItemClickCount(ctx context.Context, userID uuid.UUID) (int64, error)
	GetProfilePageViews(ctx context.Context, userID uuid.UUID) (int64, error)

	// Time range analytics
	GetUserAnalyticsByTimeRange(ctx context.Context, params TimeRangeParams) ([]DailyAnalytics, error)
	GetItemAnalyticsByTimeRange(ctx context.Context, params ItemTimeRangeParams) ([]DailyAnalytics, error)
	GetProfilePageViewsByDate(ctx context.Context, params TimeRangeParams) ([]DailyAnalytics, error)

	// Insight queries
	GetTopContentItemsByClicks(ctx context.Context, params TopItemsParams) ([]TopContentItem, error)
	GetReferrerAnalytics(ctx context.Context, params ReferrerParams) ([]ReferrerStats, error)

	// Visitor analytics
	GetUniqueVisitors(ctx context.Context, params TimeRangeParams) (int64, error)
	GetUniqueVisitorsByDay(ctx context.Context, params TimeRangeParams) ([]VisitorAnalytics, error)
}

type CreateAnalyticsParams struct {
	ItemID    uuid.UUID
	UserID    uuid.UUID
	IPAddress string
	UserAgent string
	Referrer  string
}

type CreatePageViewParams struct {
	ItemID    uuid.UUID
	UserID    uuid.UUID
	IPAddress string
	UserAgent string
	Referrer  string
}

type TimeRangeParams struct {
	UserID    uuid.UUID
	StartDate time.Time
	EndDate   time.Time
}

type ItemTimeRangeParams struct {
	ItemID    uuid.UUID
	StartDate time.Time
	EndDate   time.Time
}

type TopItemsParams struct {
	UserID    uuid.UUID
	StartDate time.Time
	EndDate   time.Time
	Limit     int
}

type ReferrerParams struct {
	UserID    uuid.UUID
	StartDate time.Time
	EndDate   time.Time
	Limit     int
}

// Output data types
type DailyAnalytics struct {
	Day    time.Time `json:"day"`
	Clicks int64     `json:"clicks"`
}

type TopContentItem struct {
	ItemID      string `json:"item_id"`
	ContentType string `json:"content_type"`
	Title       string `json:"title"`
	ClickCount  int64  `json:"click_count"`
}

type ReferrerStats struct {
	Referrer string `json:"referrer"`
	Count    int64  `json:"count"`
}

type VisitorAnalytics struct {
	Day      time.Time `json:"day"`
	Visitors int64     `json:"visitors"`
}

// Implementation
type SQLCAnalyticsRepository struct {
	db     *db.Queries
	logger log.Logger
}

func NewAnalyticsRepository(db *db.Queries, logger log.Logger) AnalyticsRepository {
	return &SQLCAnalyticsRepository{
		db:     db,
		logger: logger,
	}
}

func (r *SQLCAnalyticsRepository) CreateAnalyticsEntry(ctx context.Context, params CreateAnalyticsParams) (*db.Analytic, error) {
	r.logger.Infof("Creating analytics entry for user ID: %s, item ID: %s", params.UserID, params.ItemID)

	ipAddressPtr := &params.IPAddress
	userAgentPtr := &params.UserAgent
	referrerPtr := &params.Referrer

	if params.IPAddress == "" {
		ipAddressPtr = nil
	}

	if params.UserAgent == "" {
		userAgentPtr = nil
	}

	if params.Referrer == "" {
		referrerPtr = nil
	}

	sqlcParams := db.CreateAnalyticsEntryParams{
		ItemID:    params.ItemID,
		UserID:    params.UserID,
		IpAddress: ipAddressPtr,
		UserAgent: userAgentPtr,
		Referrer:  referrerPtr,
	}

	start := time.Now()
	entry, err := r.db.CreateAnalyticsEntry(ctx, sqlcParams)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "analytics entry")
		appErr.Log(r.logger)
		return nil, appErr
	}

	r.logger.Infof("Analytics entry created successfully for user ID: %s, item ID: %s in %v", params.UserID, params.ItemID, duration)
	return entry, nil
}

func (r *SQLCAnalyticsRepository) CreatePageViewEntry(ctx context.Context, params CreatePageViewParams) (*db.Analytic, error) {
	r.logger.Infof("Creating page view entry for user ID: %s, item ID: %s", params.UserID, params.ItemID)

	ipAddressPtr := &params.IPAddress
	userAgentPtr := &params.UserAgent
	referrerPtr := &params.Referrer

	if params.IPAddress == "" {
		ipAddressPtr = nil
	}

	if params.UserAgent == "" {
		userAgentPtr = nil
	}

	if params.Referrer == "" {
		referrerPtr = nil
	}

	sqlcParams := db.CreatePageViewEntryParams{
		ItemID:    params.ItemID,
		UserID:    params.UserID,
		IpAddress: ipAddressPtr,
		UserAgent: userAgentPtr,
		Referrer:  referrerPtr,
	}

	start := time.Now()
	entry, err := r.db.CreatePageViewEntry(ctx, sqlcParams)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "page view entry")
		appErr.Log(r.logger)
		return nil, appErr
	}

	r.logger.Infof("Page view entry created successfully for user ID: %s, item ID: %s in %v", params.UserID, params.ItemID, duration)
	return entry, nil
}

func (r *SQLCAnalyticsRepository) GetUserAnalyticsByTimeRange(ctx context.Context, params TimeRangeParams) ([]DailyAnalytics, error) {
	r.logger.Debugf("Getting user analytics for user ID: %s from %s to %s", 
		params.UserID, params.StartDate.Format(time.RFC3339), params.EndDate.Format(time.RFC3339))

	sqlcParams := db.GetUserAnalyticsByTimeRangeParams{
		UserID:      params.UserID,
		ClickedAt:   &params.StartDate,
		ClickedAt_2: &params.EndDate,
	}

	start := time.Now()
	rows, err := r.db.GetUserAnalyticsByTimeRange(ctx, sqlcParams)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "user analytics")
		appErr.Log(r.logger)
		return nil, appErr
	}

	result := make([]DailyAnalytics, len(rows))
	for i, row := range rows {
		result[i] = DailyAnalytics{
			Day:    row.Day,
			Clicks: row.Clicks,
		}
	}

	r.logger.Debugf("Retrieved %d daily analytics for user ID: %s in %v", len(result), params.UserID, duration)
	return result, nil
}

func (r *SQLCAnalyticsRepository) GetItemAnalyticsByTimeRange(ctx context.Context, params ItemTimeRangeParams) ([]DailyAnalytics, error) {
	r.logger.Debugf("Getting item analytics for item ID: %s from %s to %s", 
		params.ItemID, params.StartDate.Format(time.RFC3339), params.EndDate.Format(time.RFC3339))

	sqlcParams := db.GetItemAnalyticsByTimeRangeParams{
		ItemID:      params.ItemID,
		ClickedAt:   &params.StartDate,
		ClickedAt_2: &params.EndDate,
	}

	start := time.Now()
	rows, err := r.db.GetItemAnalyticsByTimeRange(ctx, sqlcParams)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "item analytics")
		appErr.Log(r.logger)
		return nil, appErr
	}

	result := make([]DailyAnalytics, len(rows))
	for i, row := range rows {
		result[i] = DailyAnalytics{
			Day:    row.Day,
			Clicks: row.Clicks,
		}
	}

	r.logger.Debugf("Retrieved %d daily analytics for item ID: %s in %v", len(result), params.ItemID, duration)
	return result, nil
}

func (r *SQLCAnalyticsRepository) GetProfilePageViewsByDate(ctx context.Context, params TimeRangeParams) ([]DailyAnalytics, error) {
	r.logger.Debugf("Getting profile page views for user ID: %s from %s to %s", 
		params.UserID, params.StartDate.Format(time.RFC3339), params.EndDate.Format(time.RFC3339))

	sqlcParams := db.GetProfilePageViewsByDateParams{
		UserID:      params.UserID,
		ClickedAt:   &params.StartDate,
		ClickedAt_2: &params.EndDate,
	}

	start := time.Now()
	rows, err := r.db.GetProfilePageViewsByDate(ctx, sqlcParams)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "profile page views")
		appErr.Log(r.logger)
		return nil, appErr
	}

	result := make([]DailyAnalytics, len(rows))
	for i, row := range rows {
		result[i] = DailyAnalytics{
			Day:    row.Day,
			Clicks: row.Views,
		}
	}

	r.logger.Debugf("Retrieved %d daily profile page views for user ID: %s in %v", len(result), params.UserID, duration)
	return result, nil
}

func (r *SQLCAnalyticsRepository) GetTopContentItemsByClicks(ctx context.Context, params TopItemsParams) ([]TopContentItem, error) {
	r.logger.Debugf("Getting top content items by clicks for user ID: %s, limit: %d", params.UserID, params.Limit)

	sqlcParams := db.GetTopContentItemsByClicksParams{
		UserID:      params.UserID,
		ClickedAt:   &params.StartDate,
		ClickedAt_2: &params.EndDate,
		Limit:       int64(params.Limit),
	}

	start := time.Now()
	rows, err := r.db.GetTopContentItemsByClicks(ctx, sqlcParams)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "top content items")
		appErr.Log(r.logger)
		return nil, appErr
	}

	result := make([]TopContentItem, len(rows))
	for i, row := range rows {
		result[i] = TopContentItem{
			ItemID:      row.ItemID.String(),
			ContentType: row.ContentType,
			Title:       row.Title,
			ClickCount:  row.ClickCount,
		}
	}

	r.logger.Debugf("Retrieved %d top content items for user ID: %s in %v", len(result), params.UserID, duration)
	return result, nil
}

func (r *SQLCAnalyticsRepository) GetReferrerAnalytics(ctx context.Context, params ReferrerParams) ([]ReferrerStats, error) {
	r.logger.Debugf("Getting referrer analytics for user ID: %s, limit: %d", params.UserID, params.Limit)

	sqlcParams := db.GetReferrerAnalyticsParams{
		UserID:      params.UserID,
		ClickedAt:   &params.StartDate,
		ClickedAt_2: &params.EndDate,
		Limit:       int64(params.Limit),
	}

	start := time.Now()
	rows, err := r.db.GetReferrerAnalytics(ctx, sqlcParams)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "referrer analytics")
		appErr.Log(r.logger)
		return nil, appErr
	}

	result := make([]ReferrerStats, len(rows))
	for i, row := range rows {
		result[i] = ReferrerStats{
			Referrer: row.Referrer,
			Count:    row.Count,
		}
	}

	r.logger.Debugf("Retrieved %d referrer stats for user ID: %s in %v", len(result), params.UserID, duration)
	return result, nil
}

func (r *SQLCAnalyticsRepository) GetUniqueVisitors(ctx context.Context, params TimeRangeParams) (int64, error) {
	r.logger.Debugf("Getting unique visitors count for user ID: %s from %s to %s", 
		params.UserID, params.StartDate.Format(time.RFC3339), params.EndDate.Format(time.RFC3339))

	sqlcParams := db.GetUniqueVisitorsParams{
		UserID:      params.UserID,
		ClickedAt:   &params.StartDate,
		ClickedAt_2: &params.EndDate,
	}

	start := time.Now()
	count, err := r.db.GetUniqueVisitors(ctx, sqlcParams)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "unique visitors")
		appErr.Log(r.logger)
		return 0, appErr
	}

	r.logger.Debugf("Retrieved unique visitors count: %d for user ID: %s in %v", count, params.UserID, duration)
	return count, nil
}

func (r *SQLCAnalyticsRepository) GetUniqueVisitorsByDay(ctx context.Context, params TimeRangeParams) ([]VisitorAnalytics, error) {
	r.logger.Debugf("Getting unique visitors by day for user ID: %s from %s to %s", 
		params.UserID, params.StartDate.Format(time.RFC3339), params.EndDate.Format(time.RFC3339))

	sqlcParams := db.GetUniqueVisitorsByDayParams{
		UserID:      params.UserID,
		ClickedAt:   &params.StartDate,
		ClickedAt_2: &params.EndDate,
	}

	start := time.Now()
	rows, err := r.db.GetUniqueVisitorsByDay(ctx, sqlcParams)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "unique visitors by day")
		appErr.Log(r.logger)
		return nil, appErr
	}

	result := make([]VisitorAnalytics, len(rows))
	for i, row := range rows {
		result[i] = VisitorAnalytics{
			Day:      row.Day,
			Visitors: row.Visitors,
		}
	}

	r.logger.Debugf("Retrieved %d days of visitor data for user ID: %s in %v", len(result), params.UserID, duration)
	return result, nil
}

func (r *SQLCAnalyticsRepository) GetItemAnalytics(ctx context.Context, itemID uuid.UUID, limit, offset int) ([]*db.Analytic, error) {
	r.logger.Debugf("Getting item analytics for item ID: %s, limit: %d, offset: %d", itemID, limit, offset)

	params := db.GetItemAnalyticsParams{
		ItemID: itemID,
		Limit:  int64(limit),
		Offset: int64(offset),
	}

	start := time.Now()
	analytics, err := r.db.GetItemAnalytics(ctx, params)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "item analytics")
		appErr.Log(r.logger)
		return nil, appErr
	}

	r.logger.Debugf("Retrieved %d analytics entries for item ID: %s in %v", len(analytics), itemID, duration)
	return analytics, nil
}

func (r *SQLCAnalyticsRepository) GetUserAnalytics(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*db.Analytic, error) {
	r.logger.Debugf("Getting user analytics for user ID: %s, limit: %d, offset: %d", userID, limit, offset)

	params := db.GetUserAnalyticsParams{
		UserID: userID,
		Limit:  int64(limit),
		Offset: int64(offset),
	}

	start := time.Now()
	analytics, err := r.db.GetUserAnalytics(ctx, params)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "user analytics")
		appErr.Log(r.logger)
		return nil, appErr
	}

	r.logger.Debugf("Retrieved %d analytics entries for user ID: %s in %v", len(analytics), userID, duration)
	return analytics, nil
}

func (r *SQLCAnalyticsRepository) GetContentItemClickCount(ctx context.Context, itemID uuid.UUID) (int64, error) {
	r.logger.Debugf("Getting click count for content item ID: %s", itemID)

	start := time.Now()
	count, err := r.db.GetContentItemClickCount(ctx, itemID)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "content item click count")
		appErr.Log(r.logger)
		return 0, appErr
	}

	r.logger.Debugf("Retrieved click count: %d for content item ID: %s in %v", count, itemID, duration)
	return count, nil
}

func (r *SQLCAnalyticsRepository) GetProfilePageViews(ctx context.Context, userID uuid.UUID) (int64, error) {
	r.logger.Debugf("Getting profile page views for user ID: %s", userID)

	start := time.Now()
	count, err := r.db.GetProfilePageViews(ctx, userID)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "profile page views")
		appErr.Log(r.logger)
		return 0, appErr
	}

	r.logger.Debugf("Retrieved profile page views: %d for user ID: %s in %v", count, userID, duration)
	return count, nil
}

func (r *SQLCAnalyticsRepository) GetUserItemClickCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	r.logger.Debugf("Getting total clicks for user ID: %s", userID)

	start := time.Now()
	count, err := r.db.GetUserItemClickCount(ctx, userID)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "user item click count")
		appErr.Log(r.logger)
		return 0, appErr
	}

	r.logger.Debugf("Retrieved total click count: %d for user ID: %s in %v", count, userID, duration)
	return count, nil
}