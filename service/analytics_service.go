package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/0xsj/gin-sqlc/api"
	db "github.com/0xsj/gin-sqlc/db/sqlc"
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
}

func NewAnalyticsService(
	analyticsRepo repository.AnalyticsRepository,
	contentRepo repository.ContentRepository,
	userRepo repository.UserRepository,
) AnalyticsService {
	return &analyticsService{
		analyticsRepo: analyticsRepo,
		contentRepo:   contentRepo,
		userRepo:      userRepo,
	}
}

func (s *analyticsService) RecordClick(ctx context.Context, input RecordClickInput) error {
	itemID, err := uuid.Parse(input.ItemID)
	if err != nil {
		return api.ErrInvalidInput
	}

	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return api.ErrInvalidInput
	}

	// Verify content item exists
	_, err = s.contentRepo.GetContentItem(ctx, itemID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return api.ErrNotFound
		}
		return api.ErrInternalServer
	}

	// Verify user exists
	_, err = s.userRepo.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return api.ErrNotFound
		}
		return api.ErrInternalServer
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
		return api.ErrInternalServer
	}

	return nil
}

func (s *analyticsService) RecordPageView(ctx context.Context, input RecordPageViewInput) error {
	profileID, err := uuid.Parse(input.ProfileID)
	if err != nil {
		return api.ErrInvalidInput
	}

	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return api.ErrInvalidInput
	}

	// Verify profile (content item) exists
	_, err = s.contentRepo.GetContentItem(ctx, profileID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return api.ErrNotFound
		}
		return api.ErrInternalServer
	}

	// Verify user exists
	_, err = s.userRepo.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return api.ErrNotFound
		}
		return api.ErrInternalServer
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
		return api.ErrInternalServer
	}

	return nil
}

func (s *analyticsService) GetContentItemAnalytics(ctx context.Context, itemIDStr string, page, pageSize int) (*ContentItemAnalyticsDTO, error) {
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		return nil, api.ErrInvalidInput
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
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, api.ErrNotFound
		}
		return nil, api.ErrInternalServer
	}

	// Get click count
	totalClicks, err := s.analyticsRepo.GetContentItemClickCount(ctx, itemID)
	if err != nil {
		return nil, api.ErrInternalServer
	}

	// Get analytics entries
	entries, err := s.analyticsRepo.GetItemAnalytics(ctx, itemID, pageSize, offset)
	if err != nil {
		return nil, api.ErrInternalServer
	}

	// Map to DTOs
	clickData := make([]*AnalyticsDTO, len(entries))
	for i, entry := range entries {
		clickData[i] = mapAnalyticToDTO(entry)
	}

	return &ContentItemAnalyticsDTO{
		ItemID:      itemIDStr,
		TotalClicks: totalClicks,
		ClickData:   clickData,
	}, nil
}

func (s *analyticsService) GetUserAnalytics(ctx context.Context, userIDStr string, page, pageSize int) (*UserAnalyticsDTO, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, api.ErrInvalidInput
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
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, api.ErrNotFound
		}
		return nil, api.ErrInternalServer
	}

	// Get click count
	totalClicks, err := s.analyticsRepo.GetUserItemClickCount(ctx, userID)
	if err != nil {
		return nil, api.ErrInternalServer
	}

	// Get analytics entries
	entries, err := s.analyticsRepo.GetUserAnalytics(ctx, userID, pageSize, offset)
	if err != nil {
		return nil, api.ErrInternalServer
	}

	// Map to DTOs
	clickData := make([]*AnalyticsDTO, len(entries))
	for i, entry := range entries {
		clickData[i] = mapAnalyticToDTO(entry)
	}

	return &UserAnalyticsDTO{
		UserID:      userIDStr,
		TotalClicks: totalClicks,
		ClickData:   clickData,
	}, nil
}

func (s *analyticsService) GetUserAnalyticsByTimeRange(ctx context.Context, userIDStr string, input TimeRangeInput) (*TimeRangeAnalyticsDTO, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, api.ErrInvalidInput
	}

	// Parse date strings to time.Time
	startDate, err := time.Parse(time.RFC3339, input.StartDate)
	if err != nil {
		return nil, api.ErrInvalidInput
	}

	endDate, err := time.Parse(time.RFC3339, input.EndDate)
	if err != nil {
		return nil, api.ErrInvalidInput
	}

	// Verify user exists
	_, err = s.userRepo.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, api.ErrNotFound
		}
		return nil, api.ErrInternalServer
	}

	// Get daily analytics
	params := repository.TimeRangeParams{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	dailyAnalytics, err := s.analyticsRepo.GetUserAnalyticsByTimeRange(ctx, params)
	if err != nil {
		return nil, api.ErrInternalServer
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

	return &TimeRangeAnalyticsDTO{
		UserID:      userIDStr,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		TotalClicks: totalClicks,
		DailyClicks: dailyClicks,
	}, nil
}

func (s *analyticsService) GetItemAnalyticsByTimeRange(ctx context.Context, itemIDStr string, input TimeRangeInput) (*ItemTimeRangeAnalyticsDTO, error) {
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		return nil, api.ErrInvalidInput
	}

	// Parse date strings to time.Time
	startDate, err := time.Parse(time.RFC3339, input.StartDate)
	if err != nil {
		return nil, api.ErrInvalidInput
	}

	endDate, err := time.Parse(time.RFC3339, input.EndDate)
	if err != nil {
		return nil, api.ErrInvalidInput
	}

	// Verify item exists
	_, err = s.contentRepo.GetContentItem(ctx, itemID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, api.ErrNotFound
		}
		return nil, api.ErrInternalServer
	}

	// Get daily analytics
	params := repository.ItemTimeRangeParams{
		ItemID:    itemID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	dailyAnalytics, err := s.analyticsRepo.GetItemAnalyticsByTimeRange(ctx, params)
	if err != nil {
		return nil, api.ErrInternalServer
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

	return &ItemTimeRangeAnalyticsDTO{
		ItemID:      itemIDStr,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		TotalClicks: totalClicks,
		DailyClicks: dailyClicks,
	}, nil
}

func (s *analyticsService) GetProfilePageViewsByTimeRange(ctx context.Context, userIDStr string, input TimeRangeInput) (*PageViewAnalyticsDTO, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, api.ErrInvalidInput
	}

	// Parse date strings to time.Time
	startDate, err := time.Parse(time.RFC3339, input.StartDate)
	if err != nil {
		return nil, api.ErrInvalidInput
	}

	endDate, err := time.Parse(time.RFC3339, input.EndDate)
	if err != nil {
		return nil, api.ErrInvalidInput
	}

	// Verify user exists
	_, err = s.userRepo.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, api.ErrNotFound
		}
		return nil, api.ErrInternalServer
	}

	// Get page view analytics
	params := repository.TimeRangeParams{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	dailyViews, err := s.analyticsRepo.GetProfilePageViewsByDate(ctx, params)
	if err != nil {
		return nil, api.ErrInternalServer
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
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, api.ErrInvalidInput
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
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, api.ErrNotFound
		}
		return nil, api.ErrInternalServer
	}

	// Get total page views
	totalViews, err := s.analyticsRepo.GetProfilePageViews(ctx, userID)
	if err != nil {
		return nil, api.ErrInternalServer
	}

	// Get total clicks
	totalClicks, err := s.analyticsRepo.GetUserItemClickCount(ctx, userID)
	if err != nil {
		return nil, api.ErrInternalServer
	}

	// Get unique visitors
	visitorParams := repository.TimeRangeParams{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	uniqueVisitors, err := s.analyticsRepo.GetUniqueVisitors(ctx, visitorParams)
	if err != nil {
		return nil, api.ErrInternalServer
	}

	// Get daily page views
	dailyViewsData, err := s.analyticsRepo.GetProfilePageViewsByDate(ctx, visitorParams)
	if err != nil {
		return nil, api.ErrInternalServer
	}

	// Get daily visitors
	dailyVisitorsData, err := s.analyticsRepo.GetUniqueVisitorsByDay(ctx, visitorParams)
	if err != nil {
		return nil, api.ErrInternalServer
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
		return nil, api.ErrInternalServer
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
		return nil, api.ErrInternalServer
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
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, api.ErrInvalidInput
	}

	// Parse date strings to time.Time
	startDate, err := time.Parse(time.RFC3339, input.StartDate)
	if err != nil {
		return nil, api.ErrInvalidInput
	}

	endDate, err := time.Parse(time.RFC3339, input.EndDate)
	if err != nil {
		return nil, api.ErrInvalidInput
	}

	// Set default limit if not specified
	limit := input.Limit
	if limit <= 0 {
		limit = 10
	}

	// Verify user exists
	_, err = s.userRepo.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, api.ErrNotFound
		}
		return nil, api.ErrInternalServer
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
		return nil, api.ErrInternalServer
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
