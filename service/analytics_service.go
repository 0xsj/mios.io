package service

import "context"

type AnalyticsService interface {
	// Recording data
	RecordClick(ctx context.Context, input RecordClickInput) error
	RecordPageView(ctx context.Context, input RecordPageViewInput) error
	
	// Basic analytics
	GetContentItemAnalytics(ctx context.Context, itemID string, page, pageSize int) (*ContentItemAnalyticsDTO, error)
	GetUserAnalytics(ctx context.Context, userID string, page, pageSize int) (*UserAnalyticsDTO, error)
	
	// Time range analytics
	GetUserAnalyticsByTimeRange(ctx context.Context, userID string, input TimeRangeInput) (*TimeRangeAnalyticsDTO, error)
	GetItemAnalyticsByTimeRange(ctx context.Context, itemID string, input TimeRangeInput) (*ItemTimeRangeAnalyticsDTO, error)
	GetProfilePageViewsByTimeRange(ctx context.Context, userID string, input TimeRangeInput) (*PageViewAnalyticsDTO, error)
	
	// Dashboard analytics
	GetProfileDashboard(ctx context.Context, userID string, days int) (*ProfileDashboardDTO, error)
	
	// Referrer analytics
	GetReferrerAnalytics(ctx context.Context, userID string, input TimeRangeInput) (*ReferrerAnalyticsDTO, error)
}

type RecordClickInput struct {
	ItemID    string `json:"item_id" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	Referrer  string `json:"referrer"`
}

type RecordPageViewInput struct {
	ProfileID string `json:"profile_id" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	Referrer  string `json:"referrer"`
}

type TimeRangeInput struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	Limit     int    `json:"limit"`
}

// Output types (DTOs)
type ContentItemAnalyticsDTO struct {
	ItemID      string          `json:"item_id"`
	TotalClicks int64           `json:"total_clicks"`
	ClickData   []*AnalyticsDTO `json:"click_data"`
}

type UserAnalyticsDTO struct {
	UserID      string          `json:"user_id"`
	TotalClicks int64           `json:"total_clicks"`
	ClickData   []*AnalyticsDTO `json:"click_data"`
}

type AnalyticsDTO struct {
	ID        string `json:"id"`
	ItemID    string `json:"item_id"`
	UserID    string `json:"user_id"`
	IPAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	Referrer  string `json:"referrer,omitempty"`
	PageView  bool   `json:"page_view"`
	ClickedAt string `json:"clicked_at"`
}

type TimeRangeAnalyticsDTO struct {
	UserID      string               `json:"user_id"`
	StartDate   string               `json:"start_date"`
	EndDate     string               `json:"end_date"`
	TotalClicks int64                `json:"total_clicks"`
	DailyClicks []*DailyAnalyticsDTO `json:"daily_clicks"`
}

type ItemTimeRangeAnalyticsDTO struct {
	ItemID      string               `json:"item_id"`
	StartDate   string               `json:"start_date"`
	EndDate     string               `json:"end_date"`
	TotalClicks int64                `json:"total_clicks"`
	DailyClicks []*DailyAnalyticsDTO `json:"daily_clicks"`
}

type PageViewAnalyticsDTO struct {
	UserID      string               `json:"user_id"`
	StartDate   string               `json:"start_date"`
	EndDate     string               `json:"end_date"`
	TotalViews  int64                `json:"total_views"`
	DailyViews  []*DailyAnalyticsDTO `json:"daily_views"`
}

type DailyAnalyticsDTO struct {
	Date      string `json:"date"`
	DayOfWeek string `json:"day_of_week"`
	Count     int64  `json:"count"`
}

type ProfileDashboardDTO struct {
	UserID         string                `json:"user_id"`
	Period         string                `json:"period"`
	TotalViews     int64                 `json:"total_views"`
	TotalClicks    int64                 `json:"total_clicks"`
	UniqueVisitors int64                 `json:"unique_visitors"`
	ConversionRate float64               `json:"conversion_rate"`
	DailyViews     []*DailyAnalyticsDTO  `json:"daily_views"`
	DailyVisitors  []*DailyAnalyticsDTO  `json:"daily_visitors"`
	TopItems       []*TopContentItemDTO  `json:"top_items"`
	TopReferrers   []*ReferrerStatsDTO   `json:"top_referrers"`
}

type TopContentItemDTO struct {
	ItemID      string `json:"item_id"`
	ContentType string `json:"content_type"`
	Title       string `json:"title"`
	ClickCount  int64  `json:"click_count"`
}

type ReferrerAnalyticsDTO struct {
	UserID     string             `json:"user_id"`
	StartDate  string             `json:"start_date"`
	EndDate    string             `json:"end_date"`
	TotalCount int64              `json:"total_count"`
	Referrers  []*ReferrerStatsDTO `json:"referrers"`
}

type ReferrerStatsDTO struct {
	Referrer   string  `json:"referrer"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}
