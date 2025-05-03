package analytics

import (
	"strconv"

	"github.com/0xsj/gin-sqlc/api"
	"github.com/0xsj/gin-sqlc/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	analyticsService service.AnalyticsService
}

func NewHandler(analyticsService service.AnalyticsService) *Handler {
	return &Handler{
		analyticsService: analyticsService,
	}
}

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
}

func (h *Handler) RecordClick(c *gin.Context) {
	var req RecordClickRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	input := service.RecordClickInput{
		ItemID:    req.ItemID,
		UserID:    req.UserID,
		IPAddress: req.IPAddress,
		UserAgent: req.UserAgent,
		Referrer:  req.Referrer,
	}

	err := h.analyticsService.RecordClick(c, input)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	api.RespondWithSuccess(c, nil, "Click record successfully")
}

// record page view
func (h *Handler) RecordPageView(c *gin.Context) {
	var req RecordPageViewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	input := service.RecordPageViewInput{
		ProfileID: req.ProfileID,
		UserID:    req.UserID,
		IPAddress: req.IPAddress,
		UserAgent: req.UserAgent,
		Referrer:  req.Referrer,
	}

	err := h.analyticsService.RecordPageView(c, input)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	api.RespondWithSuccess(c, nil, "Page view recorded successfully")
}

// get content item analytics
func (h *Handler) GetContentItemAnalytics(c *gin.Context) {
	itemID := c.Param("id")
	page, pageSize := getPaginationParams(c)

	analytics, err := h.analyticsService.GetContentItemAnalytics(c, itemID, page, pageSize)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	api.RespondWithSuccess(c, analytics, "Content item analytics retreived successfully")
}

// get user analytics
func (h *Handler) GetUserAnalytics(c *gin.Context) {
	userID := c.Param("id")
	page, pageSize := getPaginationParams(c)
	analytics, err := h.analyticsService.GetUserAnalytics(c, userID, page, pageSize)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	api.RespondWithSuccess(c, analytics, "User analytics retreived successfully")
}

// get user analytics by time range
func (h *Handler) GetUserAnalyticsByTimeRange(c *gin.Context) {
	userID := c.Param("id")
	var req TimeRangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	input := service.TimeRangeInput{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Limit:     req.Limit,
	}

	analytics, err := h.analyticsService.GetUserAnalyticsByTimeRange(c, userID, input)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	api.RespondWithSuccess(c, analytics, "User time range analytics retreived successfully")
}

// item
func (h *Handler) GetItemAnalyticsByTimeRange(c *gin.Context) {
	itemID := c.Param("id")
	var req TimeRangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	input := service.TimeRangeInput{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Limit:     req.Limit,
	}

	analytics, err := h.analyticsService.GetItemAnalyticsByTimeRange(c, itemID, input)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	api.RespondWithSuccess(c, analytics, "Item time range analytics retrieved successfully")
}

func (h *Handler) GetProfilePageViewsByTimeRange(c *gin.Context) {
	userID := c.Param("id")
	var req TimeRangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	input := service.TimeRangeInput{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Limit:     req.Limit,
	}

	analytics, err := h.analyticsService.GetProfilePageViewsByTimeRange(c, userID, input)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	api.RespondWithSuccess(c, analytics, "Profile page views retrieved successfully")
}

func (h *Handler) GetProfileDashboard(c *gin.Context) {
	userID := c.Param("id")
	days, err := strconv.Atoi(c.DefaultQuery("days", "30"))
	if err != nil {
		days = 30
	}

	dashboard, err := h.analyticsService.GetProfileDashboard(c, userID, days)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	api.RespondWithSuccess(c, dashboard, "Profile dashboard retrieved successfully")
}

func (h *Handler) GetReferrerAnalytics(c *gin.Context) {
	userID := c.Param("id")
	var req TimeRangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	input := service.TimeRangeInput{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Limit:     req.Limit,
	}

	analytics, err := h.analyticsService.GetReferrerAnalytics(c, userID, input)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	api.RespondWithSuccess(c, analytics, "Referrer analytics retrieved successfully")
}

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
