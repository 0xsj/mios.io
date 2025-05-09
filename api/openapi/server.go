package openapi

import (
	"net/http"

	"github.com/0xsj/gin-sqlc/log"
	"github.com/0xsj/gin-sqlc/service"
	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime/types"
)

type Handler struct {
	authService      service.AuthService       
	userService      service.UserService       
	contentService   service.ContentService    
	analyticsService service.AnalyticsService 
	logger           log.Logger
}

// GetAnalytics implements ServerInterface.
func (h *Handler) GetAnalytics(c *gin.Context) {
	panic("unimplemented")
}

// GetContent implements ServerInterface.
func (h *Handler) GetContent(c *gin.Context) {
	panic("unimplemented")
}

// GetContentContentId implements ServerInterface.
func (h *Handler) GetContentContentId(c *gin.Context, contentId types.UUID) {
	panic("unimplemented")
}

// GetUsers implements ServerInterface.
func (h *Handler) GetUsers(c *gin.Context) {
	panic("unimplemented")
}

// GetUsersUserId implements ServerInterface.
func (h *Handler) GetUsersUserId(c *gin.Context, userId types.UUID) {
	panic("unimplemented")
}

// PostAuthLogin implements ServerInterface.
func (h *Handler) PostAuthLogin(c *gin.Context) {
	panic("unimplemented")
}

// PostAuthRefresh implements ServerInterface.
func (h *Handler) PostAuthRefresh(c *gin.Context) {
	panic("unimplemented")
}

// PostContent implements ServerInterface.
func (h *Handler) PostContent(c *gin.Context) {
	panic("unimplemented")
}

// PostUsers implements ServerInterface.
func (h *Handler) PostUsers(c *gin.Context) {
	panic("unimplemented")
}

func NewHandler(
	authService service.AuthService,    
	userService service.UserService,    
	contentService service.ContentService,
	analyticsService service.AnalyticsService,
	logger log.Logger,
) *Handler {
	return &Handler{
		authService:      authService,
		userService:      userService,
		contentService:   contentService,
		analyticsService: analyticsService,
		logger:           logger,
	}
}
func RegisterOpenAPIHandlers(r *gin.Engine, h *Handler) {
	apiGroup := r.Group("/api/v1")
	
	RegisterHandlersWithOptions(apiGroup, h, GinServerOptions{})
	
}

func (h *Handler) NotImplementedYet(c *gin.Context) {
	h.logger.Warn("Endpoint not implemented yet")
	c.JSON(http.StatusNotImplemented, Error{
		Code:    http.StatusNotImplemented,
		Message: "This endpoint is not implemented yet",
	})
}
