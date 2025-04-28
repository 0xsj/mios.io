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