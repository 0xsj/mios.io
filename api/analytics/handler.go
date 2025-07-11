package analytics

import (
	"strconv"

	"github.com/0xsj/mios.io/log"
	"github.com/0xsj/mios.io/pkg/response"
	"github.com/0xsj/mios.io/service"
	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for analytics operations
type Handler struct {
	analyticsService service.AnalyticsService
	logger           log.Logger
}

// NewHandler creates a new analytics handler
func NewHandler(analyticsService service.AnalyticsService, logger log.Logger) *Handler {
	return &Handler{
		analyticsService: analyticsService,
		logger:           logger,
	}
}

// RegisterRoutes registers analytics routes on the given router
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	analyticsGroup := r.Group("/api/analytics")
	{
		analyticsGroup.POST("/clicks", h.RecordClick)
		analyticsGroup.POST("/page-views", h.RecordPageView)

		analyticsGroup.GET("/items/:id", h.GetContentItemAnalytics)
		analyticsGroup.POST("/items/:id/time-range", h.GetItemAnalyticsByTimeRange)

		analyticsGroup.GET("/users/:id", h.GetUserAnalytics)
		analyticsGroup.POST("/users/:id/time-range", h.GetUserAnalyticsByTimeRange)
		analyticsGroup.POST("/users/:id/page-views", h.GetProfilePageViewsByTimeRange)
		analyticsGroup.GET("/users/:id/dashboard", h.GetProfileDashboard)
		analyticsGroup.POST("/users/:id/referrers", h.GetReferrerAnalytics)
	}

	h.logger.Info("Analytics routes registered successfully")
}

// RecordClick records an analytics click event
func (h *Handler) RecordClick(c *gin.Context) {
	h.logger.Info("RecordClick handler called")

	var req RecordClickRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Invalid request format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, err.Error())
		return
	}

	h.logger.Debugf("Received click analytics for item ID: %s, user ID: %s", req.ItemID, req.UserID)

	input := service.RecordClickInput{
		ItemID:    req.ItemID,
		UserID:    req.UserID,
		IPAddress: req.IPAddress,
		UserAgent: req.UserAgent,
		Referrer:  req.Referrer,
	}

	err := h.analyticsService.RecordClick(c, input)
	if err != nil {
		h.logger.Errorf("Failed to record click: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Infof("Click recorded successfully for item ID: %s", req.ItemID)
	response.Success(c, nil, "Click recorded successfully")
}

// RecordPageView records a page view analytics event
func (h *Handler) RecordPageView(c *gin.Context) {
	h.logger.Info("RecordPageView handler called")

	var req RecordPageViewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Invalid request format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, err.Error())
		return
	}

	h.logger.Debugf("Received page view analytics for profile ID: %s, user ID: %s", req.ProfileID, req.UserID)

	input := service.RecordPageViewInput{
		ProfileID: req.ProfileID,
		UserID:    req.UserID,
		IPAddress: req.IPAddress,
		UserAgent: req.UserAgent,
		Referrer:  req.Referrer,
	}

	err := h.analyticsService.RecordPageView(c, input)
	if err != nil {
		h.logger.Errorf("Failed to record page view: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Infof("Page view recorded successfully for profile ID: %s", req.ProfileID)
	response.Success(c, nil, "Page view recorded successfully")
}

// GetContentItemAnalytics retrieves analytics for a specific content item
func (h *Handler) GetContentItemAnalytics(c *gin.Context) {
	itemID := c.Param("id")
	h.logger.Debugf("GetContentItemAnalytics handler called for item ID: %s", itemID)

	page, pageSize := getPaginationParams(c)
	h.logger.Debugf("Pagination parameters: page=%d, pageSize=%d", page, pageSize)

	analytics, err := h.analyticsService.GetContentItemAnalytics(c, itemID, page, pageSize)
	if err != nil {
		h.logger.Warnf("Failed to retrieve content item analytics: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Debugf("Retrieved analytics for content item ID: %s with %d entries",
		itemID, len(analytics.ClickData))
	response.Success(c, analytics, "Content item analytics retrieved successfully")
}

// GetUserAnalytics retrieves analytics for a specific user
func (h *Handler) GetUserAnalytics(c *gin.Context) {
	userID := c.Param("id")
	h.logger.Debugf("GetUserAnalytics handler called for user ID: %s", userID)

	page, pageSize := getPaginationParams(c)
	h.logger.Debugf("Pagination parameters: page=%d, pageSize=%d", page, pageSize)

	analytics, err := h.analyticsService.GetUserAnalytics(c, userID, page, pageSize)
	if err != nil {
		h.logger.Warnf("Failed to retrieve user analytics: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Debugf("Retrieved analytics for user ID: %s with %d entries",
		userID, len(analytics.ClickData))
	response.Success(c, analytics, "User analytics retrieved successfully")
}

// GetUserAnalyticsByTimeRange retrieves user analytics within a specific time range
func (h *Handler) GetUserAnalyticsByTimeRange(c *gin.Context) {
	userID := c.Param("id")
	h.logger.Debugf("GetUserAnalyticsByTimeRange handler called for user ID: %s", userID)

	var req TimeRangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Invalid request format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, err.Error())
		return
	}

	h.logger.Debugf("Time range parameters: start=%s, end=%s", req.StartDate, req.EndDate)

	input := service.TimeRangeInput{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Limit:     req.Limit,
	}

	analytics, err := h.analyticsService.GetUserAnalyticsByTimeRange(c, userID, input)
	if err != nil {
		h.logger.Warnf("Failed to retrieve user time range analytics: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Debugf("Retrieved time range analytics for user ID: %s with %d daily entries",
		userID, len(analytics.DailyClicks))
	response.Success(c, analytics, "User time range analytics retrieved successfully")
}

// GetItemAnalyticsByTimeRange retrieves content item analytics within a specific time range
func (h *Handler) GetItemAnalyticsByTimeRange(c *gin.Context) {
	itemID := c.Param("id")
	h.logger.Debugf("GetItemAnalyticsByTimeRange handler called for item ID: %s", itemID)

	var req TimeRangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Invalid request format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, err.Error())
		return
	}

	h.logger.Debugf("Time range parameters: start=%s, end=%s", req.StartDate, req.EndDate)

	input := service.TimeRangeInput{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Limit:     req.Limit,
	}

	analytics, err := h.analyticsService.GetItemAnalyticsByTimeRange(c, itemID, input)
	if err != nil {
		h.logger.Warnf("Failed to retrieve item time range analytics: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Debugf("Retrieved time range analytics for item ID: %s with %d daily entries",
		itemID, len(analytics.DailyClicks))
	response.Success(c, analytics, "Item time range analytics retrieved successfully")
}

// GetProfilePageViewsByTimeRange retrieves profile page view analytics within a time range
func (h *Handler) GetProfilePageViewsByTimeRange(c *gin.Context) {
	userID := c.Param("id")
	h.logger.Debugf("GetProfilePageViewsByTimeRange handler called for user ID: %s", userID)

	var req TimeRangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Invalid request format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, err.Error())
		return
	}

	h.logger.Debugf("Time range parameters: start=%s, end=%s", req.StartDate, req.EndDate)

	input := service.TimeRangeInput{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Limit:     req.Limit,
	}

	analytics, err := h.analyticsService.GetProfilePageViewsByTimeRange(c, userID, input)
	if err != nil {
		h.logger.Warnf("Failed to retrieve profile page views: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Debugf("Retrieved profile page views for user ID: %s with %d daily entries",
		userID, len(analytics.DailyViews))
	response.Success(c, analytics, "Profile page views retrieved successfully")
}

// GetProfileDashboard retrieves a comprehensive dashboard for a user profile
func (h *Handler) GetProfileDashboard(c *gin.Context) {
	userID := c.Param("id")
	h.logger.Debugf("GetProfileDashboard handler called for user ID: %s", userID)

	days, err := strconv.Atoi(c.DefaultQuery("days", "30"))
	if err != nil {
		h.logger.Warnf("Invalid days parameter: %v, using default of 30", err)
		days = 30
	}

	h.logger.Debugf("Dashboard time range: last %d days", days)

	dashboard, err := h.analyticsService.GetProfileDashboard(c, userID, days)
	if err != nil {
		h.logger.Warnf("Failed to retrieve profile dashboard: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Debugf("Retrieved profile dashboard for user ID: %s with %d daily views",
		userID, len(dashboard.DailyViews))
	response.Success(c, dashboard, "Profile dashboard retrieved successfully")
}

// GetReferrerAnalytics retrieves analytics about referrers to a user's content
func (h *Handler) GetReferrerAnalytics(c *gin.Context) {
	userID := c.Param("id")
	h.logger.Debugf("GetReferrerAnalytics handler called for user ID: %s", userID)

	var req TimeRangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Invalid request format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, err.Error())
		return
	}

	h.logger.Debugf("Time range parameters: start=%s, end=%s", req.StartDate, req.EndDate)

	input := service.TimeRangeInput{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Limit:     req.Limit,
	}

	analytics, err := h.analyticsService.GetReferrerAnalytics(c, userID, input)
	if err != nil {
		h.logger.Warnf("Failed to retrieve referrer analytics: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Debugf("Retrieved referrer analytics for user ID: %s with %d referrers",
		userID, len(analytics.Referrers))
	response.Success(c, analytics, "Referrer analytics retrieved successfully")
}

// getPaginationParams extracts and validates pagination parameters from the request
func getPaginationParams(c *gin.Context) (int, int) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return page, pageSize
}
