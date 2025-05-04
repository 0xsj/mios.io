package service

import (
	"context"
	"fmt"
	"time"

	db "github.com/0xsj/gin-sqlc/db/sqlc"
	"github.com/0xsj/gin-sqlc/log"
	"github.com/0xsj/gin-sqlc/pkg/errors"
	"github.com/0xsj/gin-sqlc/repository"
	"github.com/google/uuid"
)

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
	UserID     string               `json:"user_id"`
	StartDate  string               `json:"start_date"`
	EndDate    string               `json:"end_date"`
	TotalViews int64                `json:"total_views"`
	DailyViews []*DailyAnalyticsDTO `json:"daily_views"`
}

type DailyAnalyticsDTO struct {
	Date      string `json:"date"`
	DayOfWeek string `json:"day_of_week"`
	Count     int64  `json:"count"`
}

type ProfileDashboardDTO struct {
	UserID         string               `json:"user_id"`
	Period         string               `json:"period"`
	TotalViews     int64                `json:"total_views"`
	TotalClicks    int64                `json:"total_clicks"`
	UniqueVisitors int64                `json:"unique_visitors"`
	ConversionRate float64              `json:"conversion_rate"`
	DailyViews     []*DailyAnalyticsDTO `json:"daily_views"`
	DailyVisitors  []*DailyAnalyticsDTO `json:"daily_visitors"`
	TopItems       []*TopContentItemDTO `json:"top_items"`
	TopReferrers   []*ReferrerStatsDTO  `json:"top_referrers"`
}

type TopContentItemDTO struct {
	ItemID      string `json:"item_id"`
	ContentType string `json:"content_type"`
	Title       string `json:"title"`
	ClickCount  int64  `json:"click_count"`
}

type ReferrerAnalyticsDTO struct {
	UserID     string              `json:"user_id"`
	StartDate  string              `json:"start_date"`
	EndDate    string              `json:"end_date"`
	TotalCount int64               `json:"total_count"`
	Referrers  []*ReferrerStatsDTO `json:"referrers"`
}

type ReferrerStatsDTO struct {
	Referrer   string  `json:"referrer"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

type analyticsService struct {
	analyticsRepo repository.AnalyticsRepository
	contentRepo   repository.ContentRepository
	userRepo      repository.UserRepository
	logger        log.Logger
}

func NewAnalyticsService(
	analyticsRepo repository.AnalyticsRepository,
	contentRepo repository.ContentRepository,
	userRepo repository.UserRepository,
	logger log.Logger,
) AnalyticsService {
	return &analyticsService{
		analyticsRepo: analyticsRepo,
		contentRepo:   contentRepo,
		userRepo:      userRepo,
		logger:        logger,
	}
}

func (s *analyticsService) RecordClick(ctx context.Context, input RecordClickInput) error {
	s.logger.Infof("Recording click for item ID: %s from user ID: %s", input.ItemID, input.UserID)

	itemID, err := uuid.Parse(input.ItemID)
	if err != nil {
		s.logger.Warnf("Invalid item ID format: %v", err)
		return errors.NewBadRequestError("Invalid item ID format", err)
	}

	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		s.logger.Warnf("Invalid user ID format: %v", err)
		return errors.NewBadRequestError("Invalid user ID format", err)
	}

	// Verify content item exists
	_, err = s.contentRepo.GetContentItem(ctx, itemID)
	if err != nil {
		if errors.IsNotFound(err) {
			s.logger.Warnf("Content item not found with ID: %s", input.ItemID)
			return errors.NewNotFoundError("Content item not found", err)
		}
		s.logger.Errorf("Error retrieving content item: %v", err)
		return errors.Wrap(err, "Failed to retrieve content item")
	}

	// Verify user exists
	_, err = s.userRepo.GetUser(ctx, userID)
	if err != nil {
		if errors.IsNotFound(err) {
			s.logger.Warnf("User not found with ID: %s", input.UserID)
			return errors.NewNotFoundError("User not found", err)
		}
		s.logger.Errorf("Error retrieving user: %v", err)
		return errors.Wrap(err, "Failed to retrieve user")
	}

	params := repository.CreateAnalyticsParams{
		ItemID:    itemID,
		UserID:    userID,
		IPAddress: input.IPAddress,
		UserAgent: input.UserAgent,
		Referrer:  input.Referrer,
	}

	_, err = s.analyticsRepo.CreateAnalyticsEntry(ctx, params)
	if err != nil {
		s.logger.Errorf("Failed to create analytics entry: %v", err)
		return errors.Wrap(err, "Failed to record click")
	}

	s.logger.Infof("Click recorded successfully for item ID: %s from user ID: %s", input.ItemID, input.UserID)
	return nil
}

func (s *analyticsService) RecordPageView(ctx context.Context, input RecordPageViewInput) error {
	s.logger.Infof("Recording page view for profile ID: %s by user ID: %s", input.ProfileID, input.UserID)

	profileID, err := uuid.Parse(input.ProfileID)
	if err != nil {
		s.logger.Warnf("Invalid profile ID format: %v", err)
		return errors.NewBadRequestError("Invalid profile ID format", err)
	}

	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		s.logger.Warnf("Invalid user ID format: %v", err)
		return errors.NewBadRequestError("Invalid user ID format", err)
	}

	// Verify profile (content item) exists
	_, err = s.contentRepo.GetContentItem(ctx, profileID)
	if err != nil {
		if errors.IsNotFound(err) {
			s.logger.Warnf("Profile not found with ID: %s", input.ProfileID)
			return errors.NewNotFoundError("Profile not found", err)
		}
		s.logger.Errorf("Error retrieving profile: %v", err)
		return errors.Wrap(err, "Failed to retrieve profile")
	}

	// Verify user exists
	_, err = s.userRepo.GetUser(ctx, userID)
	if err != nil {
		if errors.IsNotFound(err) {
			s.logger.Warnf("User not found with ID: %s", input.UserID)
			return errors.NewNotFoundError("User not found", err)
		}
		s.logger.Errorf("Error retrieving user: %v", err)
		return errors.Wrap(err, "Failed to retrieve user")
	}

	params := repository.CreatePageViewParams{
		ItemID:    profileID,
		UserID:    userID,
		IPAddress: input.IPAddress,
		UserAgent: input.UserAgent,
		Referrer:  input.Referrer,
	}

	_, err = s.analyticsRepo.CreatePageViewEntry(ctx, params)
	if err != nil {
		s.logger.Errorf("Failed to create page view entry: %v", err)
		return errors.Wrap(err, "Failed to record page view")
	}

	s.logger.Infof("Page view recorded successfully for profile ID: %s by user ID: %s", input.ProfileID, input.UserID)
	return nil
}

func (s *analyticsService) GetContentItemAnalytics(ctx context.Context, itemIDStr string, page, pageSize int) (*ContentItemAnalyticsDTO, error) {
	s.logger.Debugf("Getting content item analytics for item ID: %s (page: %d, size: %d)", itemIDStr, page, pageSize)

	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		s.logger.Warnf("Invalid item ID format: %v", err)
		return nil, errors.NewBadRequestError("Invalid item ID format", err)
	}

	// Set default pagination values
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Verify content item exists
	_, err = s.contentRepo.GetContentItem(ctx, itemID)
	if err != nil {
		if errors.IsNotFound(err) {
			s.logger.Warnf("Content item not found with ID: %s", itemIDStr)
			return nil, errors.NewNotFoundError("Content item not found", err)
		}
		s.logger.Errorf("Error retrieving content item: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve content item")
	}

	// Get click count
	totalClicks, err := s.analyticsRepo.GetContentItemClickCount(ctx, itemID)
	if err != nil {
		s.logger.Errorf("Failed to get click count: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve click count")
	}

	// Get analytics entries
	entries, err := s.analyticsRepo.GetItemAnalytics(ctx, itemID, pageSize, offset)
	if err != nil {
		s.logger.Errorf("Failed to get analytics entries: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve analytics data")
	}

	// Map to DTOs
	clickData := make([]*AnalyticsDTO, len(entries))
	for i, entry := range entries {
		clickData[i] = mapAnalyticToDTO(entry)
	}

	s.logger.Debugf("Retrieved %d analytics entries for item ID: %s with total clicks: %d", len(clickData), itemIDStr, totalClicks)
	return &ContentItemAnalyticsDTO{
		ItemID:      itemIDStr,
		TotalClicks: totalClicks,
		ClickData:   clickData,
	}, nil
}

func (s *analyticsService) GetUserAnalytics(ctx context.Context, userIDStr string, page, pageSize int) (*UserAnalyticsDTO, error) {
	s.logger.Debugf("Getting user analytics for user ID: %s (page: %d, size: %d)", userIDStr, page, pageSize)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		s.logger.Warnf("Invalid user ID format: %v", err)
		return nil, errors.NewBadRequestError("Invalid user ID format", err)
	}

	// Set default pagination values
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Verify user exists
	_, err = s.userRepo.GetUser(ctx, userID)
	if err != nil {
		if errors.IsNotFound(err) {
			s.logger.Warnf("User not found with ID: %s", userIDStr)
			return nil, errors.NewNotFoundError("User not found", err)
		}
		s.logger.Errorf("Error retrieving user: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve user")
	}

	// Get click count
	totalClicks, err := s.analyticsRepo.GetUserItemClickCount(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to get click count: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve click count")
	}

	// Get analytics entries
	entries, err := s.analyticsRepo.GetUserAnalytics(ctx, userID, pageSize, offset)
	if err != nil {
		s.logger.Errorf("Failed to get analytics entries: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve analytics data")
	}

	// Map to DTOs
	clickData := make([]*AnalyticsDTO, len(entries))
	for i, entry := range entries {
		clickData[i] = mapAnalyticToDTO(entry)
	}

	s.logger.Debugf("Retrieved %d analytics entries for user ID: %s with total clicks: %d", len(clickData), userIDStr, totalClicks)
	return &UserAnalyticsDTO{
		UserID:      userIDStr,
		TotalClicks: totalClicks,
		ClickData:   clickData,
	}, nil
}

func (s *analyticsService) GetUserAnalyticsByTimeRange(ctx context.Context, userIDStr string, input TimeRangeInput) (*TimeRangeAnalyticsDTO, error) {
	s.logger.Debugf("Getting user analytics by time range for user ID: %s from %s to %s", 
		userIDStr, input.StartDate, input.EndDate)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		s.logger.Warnf("Invalid user ID format: %v", err)
		return nil, errors.NewBadRequestError("Invalid user ID format", err)
	}

	// Parse date strings to time.Time
	startDate, err := time.Parse(time.RFC3339, input.StartDate)
	if err != nil {
		s.logger.Warnf("Invalid start date format: %v", err)
		return nil, errors.NewValidationError("Invalid start date format, expected RFC3339", err)
	}

	endDate, err := time.Parse(time.RFC3339, input.EndDate)
	if err != nil {
		s.logger.Warnf("Invalid end date format: %v", err)
		return nil, errors.NewValidationError("Invalid end date format, expected RFC3339", err)
	}

	// Verify user exists
	_, err = s.userRepo.GetUser(ctx, userID)
	if err != nil {
		if errors.IsNotFound(err) {
			s.logger.Warnf("User not found with ID: %s", userIDStr)
			return nil, errors.NewNotFoundError("User not found", err)
		}
		s.logger.Errorf("Error retrieving user: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve user")
	}

	// Get daily analytics
	params := repository.TimeRangeParams{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	dailyAnalytics, err := s.analyticsRepo.GetUserAnalyticsByTimeRange(ctx, params)
	if err != nil {
		s.logger.Errorf("Failed to get daily analytics: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve analytics data")
	}

	// Calculate total clicks
	var totalClicks int64
	dailyClicks := make([]*DailyAnalyticsDTO, len(dailyAnalytics))

	for i, da := range dailyAnalytics {
		totalClicks += da.Clicks
		dailyClicks[i] = &DailyAnalyticsDTO{
			Date:      da.Day.Format("2006-01-02"),
			DayOfWeek: da.Day.Format("Monday"),
			Count:     da.Clicks,
		}
	}

	s.logger.Debugf("Retrieved %d days of analytics for user ID: %s with total clicks: %d", 
		len(dailyClicks), userIDStr, totalClicks)
	
	return &TimeRangeAnalyticsDTO{
		UserID:      userIDStr,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		TotalClicks: totalClicks,
		DailyClicks: dailyClicks,
	}, nil
}

func (s *analyticsService) GetItemAnalyticsByTimeRange(ctx context.Context, itemIDStr string, input TimeRangeInput) (*ItemTimeRangeAnalyticsDTO, error) {
	s.logger.Debugf("Getting item analytics by time range for item ID: %s from %s to %s", 
		itemIDStr, input.StartDate, input.EndDate)

	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		s.logger.Warnf("Invalid item ID format: %v", err)
		return nil, errors.NewBadRequestError("Invalid item ID format", err)
	}

	// Parse date strings to time.Time
	startDate, err := time.Parse(time.RFC3339, input.StartDate)
	if err != nil {
		s.logger.Warnf("Invalid start date format: %v", err)
		return nil, errors.NewValidationError("Invalid start date format, expected RFC3339", err)
	}

	endDate, err := time.Parse(time.RFC3339, input.EndDate)
	if err != nil {
		s.logger.Warnf("Invalid end date format: %v", err)
		return nil, errors.NewValidationError("Invalid end date format, expected RFC3339", err)
	}

	// Verify item exists
	_, err = s.contentRepo.GetContentItem(ctx, itemID)
	if err != nil {
		if errors.IsNotFound(err) {
			s.logger.Warnf("Content item not found with ID: %s", itemIDStr)
			return nil, errors.NewNotFoundError("Content item not found", err)
		}
		s.logger.Errorf("Error retrieving content item: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve content item")
	}

	// Get daily analytics
	params := repository.ItemTimeRangeParams{
		ItemID:    itemID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	dailyAnalytics, err := s.analyticsRepo.GetItemAnalyticsByTimeRange(ctx, params)
	if err != nil {
		s.logger.Errorf("Failed to get daily analytics: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve analytics data")
	}

	// Calculate total clicks
	var totalClicks int64
	dailyClicks := make([]*DailyAnalyticsDTO, len(dailyAnalytics))

	for i, da := range dailyAnalytics {
		totalClicks += da.Clicks
		dailyClicks[i] = &DailyAnalyticsDTO{
			Date:      da.Day.Format("2006-01-02"),
			DayOfWeek: da.Day.Format("Monday"),
			Count:     da.Clicks,
		}
	}

	s.logger.Debugf("Retrieved %d days of analytics for item ID: %s with total clicks: %d", 
		len(dailyClicks), itemIDStr, totalClicks)
	
	return &ItemTimeRangeAnalyticsDTO{
		ItemID:      itemIDStr,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		TotalClicks: totalClicks,
		DailyClicks: dailyClicks,
	}, nil
}

func (s *analyticsService) GetProfilePageViewsByTimeRange(ctx context.Context, userIDStr string, input TimeRangeInput) (*PageViewAnalyticsDTO, error) {
	s.logger.Debugf("Getting profile page views by time range for user ID: %s from %s to %s", 
		userIDStr, input.StartDate, input.EndDate)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		s.logger.Warnf("Invalid user ID format: %v", err)
		return nil, errors.NewBadRequestError("Invalid user ID format", err)
	}

	// Parse date strings to time.Time
	startDate, err := time.Parse(time.RFC3339, input.StartDate)
	if err != nil {
		s.logger.Warnf("Invalid start date format: %v", err)
		return nil, errors.NewValidationError("Invalid start date format, expected RFC3339", err)
	}

	endDate, err := time.Parse(time.RFC3339, input.EndDate)
	if err != nil {
		s.logger.Warnf("Invalid end date format: %v", err)
		return nil, errors.NewValidationError("Invalid end date format, expected RFC3339", err)
	}

	// Verify user exists
	_, err = s.userRepo.GetUser(ctx, userID)
	if err != nil {
		if errors.IsNotFound(err) {
			s.logger.Warnf("User not found with ID: %s", userIDStr)
			return nil, errors.NewNotFoundError("User not found", err)
		}
		s.logger.Errorf("Error retrieving user: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve user")
	}

	// Get page view analytics
	params := repository.TimeRangeParams{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	dailyViews, err := s.analyticsRepo.GetProfilePageViewsByDate(ctx, params)
	if err != nil {
		s.logger.Errorf("Failed to get daily page views: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve page view data")
	}

	// Calculate total views
	var totalViews int64
	views := make([]*DailyAnalyticsDTO, len(dailyViews))

	for i, dv := range dailyViews {
		totalViews += dv.Clicks // Clicks field repurposed for views
		views[i] = &DailyAnalyticsDTO{
			Date:      dv.Day.Format("2006-01-02"),
			DayOfWeek: dv.Day.Format("Monday"),
			Count:     dv.Clicks,
		}
	}

	s.logger.Debugf("Retrieved %d days of page views for user ID: %s with total views: %d", 
		len(views), userIDStr, totalViews)
	
	return &PageViewAnalyticsDTO{
		UserID:     userIDStr,
		StartDate:  input.StartDate,
		EndDate:    input.EndDate,
		TotalViews: totalViews,
		DailyViews: views,
	}, nil
}

// Dashboard analytics
func (s *analyticsService) GetProfileDashboard(ctx context.Context, userIDStr string, days int) (*ProfileDashboardDTO, error) {
	s.logger.Infof("Getting profile dashboard for user ID: %s over %d days", userIDStr, days)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		s.logger.Warnf("Invalid user ID format: %v", err)
		return nil, errors.NewBadRequestError("Invalid user ID format", err)
	}

	// Set default period if not specified
	if days <= 0 {
		days = 30
	}

	// Calculate date range
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	// Verify user exists
	_, err = s.userRepo.GetUser(ctx, userID)
	if err != nil {
		if errors.IsNotFound(err) {
			s.logger.Warnf("User not found with ID: %s", userIDStr)
			return nil, errors.NewNotFoundError("User not found", err)
		}
		s.logger.Errorf("Error retrieving user: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve user")
	}

	// Get total page views
	totalViews, err := s.analyticsRepo.GetProfilePageViews(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to get total page views: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve page view count")
	}

	// Get total clicks
	totalClicks, err := s.analyticsRepo.GetUserItemClickCount(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to get total clicks: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve click count")
	}

	// Get unique visitors
	visitorParams := repository.TimeRangeParams{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	uniqueVisitors, err := s.analyticsRepo.GetUniqueVisitors(ctx, visitorParams)
	if err != nil {
		s.logger.Errorf("Failed to get unique visitors: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve unique visitor count")
	}

	// Get daily page views
	dailyViewsData, err := s.analyticsRepo.GetProfilePageViewsByDate(ctx, visitorParams)
	if err != nil {
		s.logger.Errorf("Failed to get daily page views: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve daily page views")
	}

	// Get daily visitors
	dailyVisitorsData, err := s.analyticsRepo.GetUniqueVisitorsByDay(ctx, visitorParams)
	if err != nil {
		s.logger.Errorf("Failed to get daily visitors: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve daily visitors")
	}

	// Get top items
	topItemsParams := repository.TopItemsParams{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
		Limit:     10, // Top 10 items
	}

	topItems, err := s.analyticsRepo.GetTopContentItemsByClicks(ctx, topItemsParams)
	if err != nil {
		s.logger.Errorf("Failed to get top content items: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve top content items")
	}

	// Get top referrers
	referrerParams := repository.ReferrerParams{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
		Limit:     5, // Top 5 referrers
	}

	topReferrers, err := s.analyticsRepo.GetReferrerAnalytics(ctx, referrerParams)
	if err != nil {
		s.logger.Errorf("Failed to get referrer analytics: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve referrer analytics")
	}

	// Calculate conversion rate
	var conversionRate float64
	if totalViews > 0 {
		conversionRate = float64(totalClicks) / float64(totalViews) * 100
	}

	// Map to DTOs
	dailyViews := make([]*DailyAnalyticsDTO, len(dailyViewsData))
	for i, dv := range dailyViewsData {
		dailyViews[i] = &DailyAnalyticsDTO{
			Date:      dv.Day.Format("2006-01-02"),
			DayOfWeek: dv.Day.Format("Monday"),
			Count:     dv.Clicks, // Clicks field repurposed for views
		}
	}

	dailyVisitors := make([]*DailyAnalyticsDTO, len(dailyVisitorsData))
	for i, dv := range dailyVisitorsData {
		dailyVisitors[i] = &DailyAnalyticsDTO{
			Date:      dv.Day.Format("2006-01-02"),
			DayOfWeek: dv.Day.Format("Monday"),
			Count:     dv.Visitors,
		}
	}

	topItemsDTO := make([]*TopContentItemDTO, len(topItems))
	for i, item := range topItems {
		topItemsDTO[i] = &TopContentItemDTO{
			ItemID:      item.ItemID,
			ContentType: item.ContentType,
			Title:       item.Title,
			ClickCount:  item.ClickCount,
		}
	}

	// Calculate percentages for referrers
	referrersDTO := make([]*ReferrerStatsDTO, len(topReferrers))
	var totalReferrerCount int64

	for _, ref := range topReferrers {
		totalReferrerCount += ref.Count
	}

	for i, ref := range topReferrers {
		var percentage float64
		if totalReferrerCount > 0 {
			percentage = float64(ref.Count) / float64(totalReferrerCount) * 100
		}

		referrersDTO[i] = &ReferrerStatsDTO{
			Referrer:   ref.Referrer,
			Count:      ref.Count,
			Percentage: percentage,
		}
	}

	periodLabel := fmt.Sprintf("Last %d days", days)

	s.logger.Infof("Retrieved dashboard data for user ID: %s with %d views, %d clicks, and %d unique visitors", 
		userIDStr, totalViews, totalClicks, uniqueVisitors)
	
	return &ProfileDashboardDTO{
		UserID:         userIDStr,
		Period:         periodLabel,
		TotalViews:     totalViews,
		TotalClicks:    totalClicks,
		UniqueVisitors: uniqueVisitors,
		ConversionRate: conversionRate,
		DailyViews:     dailyViews,
		DailyVisitors:  dailyVisitors,
		TopItems:       topItemsDTO,
		TopReferrers:   referrersDTO,
	}, nil
}

// Referrer analytics
func (s *analyticsService) GetReferrerAnalytics(ctx context.Context, userIDStr string, input TimeRangeInput) (*ReferrerAnalyticsDTO, error) {
	s.logger.Debugf("Getting referrer analytics for user ID: %s from %s to %s", 
		userIDStr, input.StartDate, input.EndDate)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		s.logger.Warnf("Invalid user ID format: %v", err)
		return nil, errors.NewBadRequestError("Invalid user ID format", err)
	}

	// Parse date strings to time.Time
	startDate, err := time.Parse(time.RFC3339, input.StartDate)
	if err != nil {
		s.logger.Warnf("Invalid start date format: %v", err)
		return nil, errors.NewValidationError("Invalid start date format, expected RFC3339", err)
	}

	endDate, err := time.Parse(time.RFC3339, input.EndDate)
	if err != nil {
		s.logger.Warnf("Invalid end date format: %v", err)
		return nil, errors.NewValidationError("Invalid end date format, expected RFC3339", err)
	}

	// Set default limit if not specified
	limit := input.Limit
	if limit <= 0 {
		limit = 10
	}

	// Verify user exists
	_, err = s.userRepo.GetUser(ctx, userID)
	if err != nil {
		if errors.IsNotFound(err) {
			s.logger.Warnf("User not found with ID: %s", userIDStr)
			return nil, errors.NewNotFoundError("User not found", err)
		}
		s.logger.Errorf("Error retrieving user: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve user")
	}

	// Get referrer stats
	params := repository.ReferrerParams{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
		Limit:     limit,
	}

	referrers, err := s.analyticsRepo.GetReferrerAnalytics(ctx, params)
	if err != nil {
		s.logger.Errorf("Failed to get referrer analytics: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve referrer data")
	}

	// Calculate total count and percentages
	var totalCount int64
	referrersDTO := make([]*ReferrerStatsDTO, len(referrers))

	for _, ref := range referrers {
		totalCount += ref.Count
	}

	for i, ref := range referrers {
		var percentage float64
		if totalCount > 0 {
			percentage = float64(ref.Count) / float64(totalCount) * 100
		}

		referrersDTO[i] = &ReferrerStatsDTO{
			Referrer:   ref.Referrer,
			Count:      ref.Count,
			Percentage: percentage,
		}
	}

	s.logger.Debugf("Retrieved %d referrers for user ID: %s with total count: %d", 
		len(referrersDTO), userIDStr, totalCount)
	
	return &ReferrerAnalyticsDTO{
		UserID:     userIDStr,
		StartDate:  input.StartDate,
		EndDate:    input.EndDate,
		TotalCount: totalCount,
		Referrers:  referrersDTO,
	}, nil
}

// Helper function to map Analytic db model to DTO
func mapAnalyticToDTO(a *db.Analytic) *AnalyticsDTO {
	dto := &AnalyticsDTO{
		ID:       a.AnalyticsID.String(),
		ItemID:   a.ItemID.String(),
		UserID:   a.UserID.String(),
		PageView: *a.PageView,
	}

	if a.IpAddress != nil {
		dto.IPAddress = *a.IpAddress
	}

	if a.UserAgent != nil {
		dto.UserAgent = *a.UserAgent
	}

	if a.Referrer != nil {
		dto.Referrer = *a.Referrer
	}

	if a.ClickedAt != nil {
		dto.ClickedAt = a.ClickedAt.Format(time.RFC3339)
	}

	return dto
}