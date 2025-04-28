package analytics

import (
	"github.com/0xsj/gin-sqlc/api"
	"github.com/0xsj/gin-sqlc/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	analyticsService service.AnalyticsService
}

func NewHandler(analyticsService service.AnalyticsService) *Handler{
	return &Handler{
		analyticsService: analyticsService,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	analyticsGroup := r.Group("/api/analytics")
	{
		analyticsGroup.POST("/clicks", h.RecordClick)
	}
}

func (h *Handler) RecordClick(c *gin.Context) {
	var req RecordClickRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	input := service.RecordClickInput{
		ItemID: req.ItemID,
		UserID: req.UserID,
		IPAddress: req.IPAddress,
		UserAgent: req.UserAgent,
		Referrer: req.Referrer,
	}

	err := h.analyticsService.RecordClick(c, input)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	api.RespondWithSuccess(c, nil,"Click record successfully")
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
		UserID: req.UserID,
		IPAddress: req.IPAddress,
		UserAgent: req.UserAgent,
		Referrer: req.Referrer,
	}

	err := h.analyticsService.RecordPageView(c, input)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	api.RespondWithSuccess(c, nil, "Page view recorded successfully")
}