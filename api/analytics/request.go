package analytics

// Request types

// RecordClickRequest represents the payload for recording a click event
type RecordClickRequest struct {
	ItemID    string `json:"item_id" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	Referrer  string `json:"referrer"`
}

// RecordPageViewRequest represents the payload for recording a page view event
type RecordPageViewRequest struct {
	ProfileID string `json:"profile_id" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	Referrer  string `json:"referrer"`
}

// TimeRangeRequest represents the payload for time-range based analytics queries
type TimeRangeRequest struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	Limit     int    `json:"limit"`
}

// Response types

// AnalyticsEntry represents a single analytics event in responses
type AnalyticsEntry struct {
	ID        string `json:"id"`
	ItemID    string `json:"item_id"`
	UserID    string `json:"user_id"`
	IPAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	Referrer  string `json:"referrer,omitempty"`
	PageView  bool   `json:"page_view"`
	ClickedAt string `json:"clicked_at"`
}

// ContentItemAnalyticsResponse represents analytics for a content item
type ContentItemAnalyticsResponse struct {
	ItemID      string           `json:"item_id"`
	TotalClicks int64            `json:"total_clicks"`
	ClickData   []*AnalyticsEntry `json:"click_data"`
}

// UserAnalyticsResponse represents analytics for a user
type UserAnalyticsResponse struct {
	UserID      string           `json:"user_id"`
	TotalClicks int64            `json:"total_clicks"`
	ClickData   []*AnalyticsEntry `json:"click_data"`
}

// DailyAnalyticsEntry represents analytics aggregated by day
type DailyAnalyticsEntry struct {
	Date      string `json:"date"`
	DayOfWeek string `json:"day_of_week"`
	Count     int64  `json:"count"`
}

// TimeRangeAnalyticsResponse represents analytics over a time range
type TimeRangeAnalyticsResponse struct {
	UserID      string                `json:"user_id"`
	StartDate   string                `json:"start_date"`
	EndDate     string                `json:"end_date"`
	TotalClicks int64                 `json:"total_clicks"`
	DailyClicks []*DailyAnalyticsEntry `json:"daily_clicks"`
}

// ItemTimeRangeAnalyticsResponse represents item analytics over a time range
type ItemTimeRangeAnalyticsResponse struct {
	ItemID      string                `json:"item_id"`
	StartDate   string                `json:"start_date"`
	EndDate     string                `json:"end_date"`
	TotalClicks int64                 `json:"total_clicks"`
	DailyClicks []*DailyAnalyticsEntry `json:"daily_clicks"`
}

// PageViewAnalyticsResponse represents page view analytics over a time range
type PageViewAnalyticsResponse struct {
	UserID     string                `json:"user_id"`
	StartDate  string                `json:"start_date"`
	EndDate    string                `json:"end_date"`
	TotalViews int64                 `json:"total_views"`
	DailyViews []*DailyAnalyticsEntry `json:"daily_views"`
}

// TopContentItemEntry represents a popular content item
type TopContentItemEntry struct {
	ItemID      string `json:"item_id"`
	ContentType string `json:"content_type"`
	Title       string `json:"title"`
	ClickCount  int64  `json:"click_count"`
}

// ReferrerEntry represents a traffic source
type ReferrerEntry struct {
	Referrer   string  `json:"referrer"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

// ReferrerAnalyticsResponse represents analytics about traffic sources
type ReferrerAnalyticsResponse struct {
	UserID     string           `json:"user_id"`
	StartDate  string           `json:"start_date"`
	EndDate    string           `json:"end_date"`
	TotalCount int64            `json:"total_count"`
	Referrers  []*ReferrerEntry `json:"referrers"`
}

// ProfileDashboardResponse represents a comprehensive analytics dashboard
type ProfileDashboardResponse struct {
	UserID         string                `json:"user_id"`
	Period         string                `json:"period"`
	TotalViews     int64                 `json:"total_views"`
	TotalClicks    int64                 `json:"total_clicks"`
	UniqueVisitors int64                 `json:"unique_visitors"`
	ConversionRate float64               `json:"conversion_rate"`
	DailyViews     []*DailyAnalyticsEntry `json:"daily_views"`
	DailyVisitors  []*DailyAnalyticsEntry `json:"daily_visitors"`
	TopItems       []*TopContentItemEntry `json:"top_items"`
	TopReferrers   []*ReferrerEntry      `json:"top_referrers"`
}