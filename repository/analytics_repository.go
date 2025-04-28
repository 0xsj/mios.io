package repository

import (
	"context"
	"time"

	db "github.com/0xsj/gin-sqlc/db/sqlc"
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
	ItemID     uuid.UUID
	UserID     uuid.UUID
	IPAddress  string
	UserAgent  string
	Referrer   string
}

type CreatePageViewParams struct {
	ItemID     uuid.UUID
	UserID     uuid.UUID
	IPAddress  string
	UserAgent  string
	Referrer   string
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
	db *db.Queries
}

func NewAnalyticsRepository(db *db.Queries) AnalyticsRepository {
	return &SQLCAnalyticsRepository{
		db: db,
	}
}

func (r *SQLCAnalyticsRepository) CreateAnalyticsEntry(ctx context.Context, params CreateAnalyticsParams) (*db.Analytic, error) {
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
		ItemID: params.ItemID,
		UserID: params.UserID,
		IpAddress: ipAddressPtr,
		UserAgent: userAgentPtr,
		Referrer: referrerPtr,
	}
	entry, err := r.db.CreateAnalyticsEntry(ctx, sqlcParams)
	if err != nil {
		return nil, ErrDatabase
	}
	
	return entry, nil
}

func ( r*SQLCAnalyticsRepository) CreatePageViewEntry(ctx context.Context, params CreatePageViewParams) (*db.Analytic, error) {
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
        ItemID:     params.ItemID,
        UserID:     params.UserID,
        IpAddress:  ipAddressPtr,
        UserAgent:  userAgentPtr,
        Referrer:   referrerPtr,
    }
    
    entry, err := r.db.CreatePageViewEntry(ctx, sqlcParams)
    if err != nil {
        return nil, ErrDatabase
    }
    
    return entry, nil
}

func (r *SQLCAnalyticsRepository) GetUserAnalyticsByTimeRange(ctx context.Context, params TimeRangeParams) ([]DailyAnalytics, error) {
    sqlcParams := db.GetUserAnalyticsByTimeRangeParams{
        UserID:      params.UserID,
        ClickedAt:   &params.StartDate,
        ClickedAt_2: &params.EndDate,
    }
    
    rows, err := r.db.GetUserAnalyticsByTimeRange(ctx, sqlcParams)
    if err != nil {
        return nil, ErrDatabase
    }
    
    result := make([]DailyAnalytics, len(rows))
    for i, row := range rows {
        result[i] = DailyAnalytics{
            Day:    row.Day,
            Clicks: row.Clicks,
        }
    }
    
    return result, nil
}

func (r *SQLCAnalyticsRepository) GetItemAnalyticsByTimeRange(ctx context.Context, params ItemTimeRangeParams) ([]DailyAnalytics, error) {
    sqlcParams := db.GetItemAnalyticsByTimeRangeParams{
        ItemID:      params.ItemID,
        ClickedAt:   &params.StartDate,
        ClickedAt_2: &params.EndDate,
    }
    
    rows, err := r.db.GetItemAnalyticsByTimeRange(ctx, sqlcParams)
    if err != nil {
        return nil, ErrDatabase
    }
    
    result := make([]DailyAnalytics, len(rows))
    for i, row := range rows {
        result[i] = DailyAnalytics{
            Day:    row.Day,
            Clicks: row.Clicks,
        }
    }
    
    return result, nil
}

func (r *SQLCAnalyticsRepository) GetProfilePageViewsByDate(ctx context.Context, params TimeRangeParams) ([]DailyAnalytics, error) {
    sqlcParams := db.GetProfilePageViewsByDateParams{
        UserID:      params.UserID,
        ClickedAt:   &params.StartDate,
        ClickedAt_2: &params.EndDate,
    }
    
    rows, err := r.db.GetProfilePageViewsByDate(ctx, sqlcParams)
    if err != nil {
        return nil, ErrDatabase
    }
    
    result := make([]DailyAnalytics, len(rows))
    for i, row := range rows {
        result[i] = DailyAnalytics{
            Day:    row.Day,
            Clicks: row.Views, // Repurposing 'Clicks' field for views
        }
    }
    
    return result, nil
}

func (r *SQLCAnalyticsRepository) GetTopContentItemsByClicks(ctx context.Context, params TopItemsParams) ([]TopContentItem, error) {
    sqlcParams := db.GetTopContentItemsByClicksParams{
        UserID:      params.UserID,
        ClickedAt:   &params.StartDate,
        ClickedAt_2: &params.EndDate,
        Limit:       int64(params.Limit),
    }
    
    rows, err := r.db.GetTopContentItemsByClicks(ctx, sqlcParams)
    if err != nil {
        return nil, ErrDatabase
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
    
    return result, nil
}

func (r *SQLCAnalyticsRepository) GetReferrerAnalytics(ctx context.Context, params ReferrerParams) ([]ReferrerStats, error) {
    sqlcParams := db.GetReferrerAnalyticsParams{
        UserID:      params.UserID,
        ClickedAt:   &params.StartDate,
        ClickedAt_2: &params.EndDate,
        Limit:       int64(params.Limit),
    }
    
    rows, err := r.db.GetReferrerAnalytics(ctx, sqlcParams)
    if err != nil {
        return nil, ErrDatabase
    }
    
    result := make([]ReferrerStats, len(rows))
    for i, row := range rows {
        result[i] = ReferrerStats{
            Referrer: row.Referrer,
            Count:    row.Count,
        }
    }
    
    return result, nil
}

func (r *SQLCAnalyticsRepository) GetUniqueVisitors(ctx context.Context, params TimeRangeParams) (int64, error) {
    sqlcParams := db.GetUniqueVisitorsParams{
        UserID:      params.UserID,
        ClickedAt:   &params.StartDate,
        ClickedAt_2: &params.EndDate,
    }
    
    count, err := r.db.GetUniqueVisitors(ctx, sqlcParams)
    if err != nil {
        return 0, ErrDatabase
    }
    
    return count, nil
}

func (r *SQLCAnalyticsRepository) GetUniqueVisitorsByDay(ctx context.Context, params TimeRangeParams) ([]VisitorAnalytics, error) {
    sqlcParams := db.GetUniqueVisitorsByDayParams{
        UserID:      params.UserID,
        ClickedAt:   &params.StartDate,
        ClickedAt_2: &params.EndDate,
    }
    
    rows, err := r.db.GetUniqueVisitorsByDay(ctx, sqlcParams)
    if err != nil {
        return nil, ErrDatabase
    }
    
    result := make([]VisitorAnalytics, len(rows))
    for i, row := range rows {
        result[i] = VisitorAnalytics{
            Day:      row.Day,
            Visitors: row.Visitors,
        }
    }
    
    return result, nil
}

func (r *SQLCAnalyticsRepository) GetItemAnalytics(ctx context.Context, itemID uuid.UUID, limit, offset int) ([]*db.Analytic, error) {
    params := db.GetItemAnalyticsParams{
        ItemID: itemID,
        Limit:  int64(limit),
        Offset: int64(offset),
    }
    
    analytics, err := r.db.GetItemAnalytics(ctx, params)
    if err != nil {
        return nil, ErrDatabase
    }
    
    return analytics, nil
}

func (r *SQLCAnalyticsRepository) GetUserAnalytics(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*db.Analytic, error) {
    params := db.GetUserAnalyticsParams{
        UserID: userID,
        Limit:  int64(limit),
        Offset: int64(offset),
    }
    
    analytics, err := r.db.GetUserAnalytics(ctx, params)
    if err != nil {
        return nil, ErrDatabase
    }
    
    return analytics, nil
}


// GetContentItemClickCount implements AnalyticsRepository.
func (r *SQLCAnalyticsRepository) GetContentItemClickCount(ctx context.Context, itemID uuid.UUID) (int64, error) {
    count, err := r.db.GetContentItemClickCount(ctx, itemID)
    if err != nil {
        return 0, ErrDatabase
    }
    
    return count, nil
}

// GetProfilePageViews implements AnalyticsRepository.
func (r *SQLCAnalyticsRepository) GetProfilePageViews(ctx context.Context, userID uuid.UUID) (int64, error) {
    count, err := r.db.GetProfilePageViews(ctx, userID)
    if err != nil {
        return 0, ErrDatabase
    }
    
    return count, nil
}

// GetUserItemClickCount implements AnalyticsRepository.
func (r *SQLCAnalyticsRepository) GetUserItemClickCount(ctx context.Context, userID uuid.UUID) (int64, error) {
    count, err := r.db.GetUserItemClickCount(ctx, userID)
    if err != nil {
        return 0, ErrDatabase
    }
    
    return count, nil
}