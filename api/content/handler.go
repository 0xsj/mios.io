package content

import (
	"net/http"

	"github.com/0xsj/gin-sqlc/log"
	"github.com/0xsj/gin-sqlc/pkg/response"
	"github.com/0xsj/gin-sqlc/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for content operations
type Handler struct {
	contentService service.ContentService
	logger         log.Logger
}

// NewHandler creates a new content handler
func NewHandler(contentService service.ContentService, logger log.Logger) *Handler {
	return &Handler{
		contentService: contentService,
		logger:         logger,
	}
}

// RegisterRoutes registers content routes on the given router
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	contentGroup := r.Group("/api/content")
	{
		contentGroup.POST("", h.CreateContentItem)
		contentGroup.GET("/:id", h.GetContentItem)
		contentGroup.GET("/user/:user_id", h.GetUserContentItems)
		contentGroup.PUT("/:id", h.UpdateContentItem)
		contentGroup.PATCH("/:id/position", h.UpdateContentItemPosition)
		contentGroup.DELETE("/:id", h.DeleteContentItem)
	}

	h.logger.Info("Content routes registered successfully")
}

// CreateContentItem creates a new content item
func (h *Handler) CreateContentItem(c *gin.Context) {
	h.logger.Info("CreateContentItem handler called")

	var req CreateContentItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Invalid request format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, err.Error())
		return
	}

	h.logger.Debugf("Received create content item request for user ID: %s", req.UserID)

	input := service.CreateContentItemInput{
		UserID:       req.UserID,
		ContentID:    req.ContentID,
		ContentType:  req.ContentType,
		Title:        req.Title,
		Href:         req.Href,
		URL:          req.URL,
		MediaType:    req.MediaType,
		DesktopX:     req.DesktopX,
		DesktopY:     req.DesktopY,
		DesktopStyle: req.DesktopStyle,
		MobileX:      req.MobileX,
		MobileY:      req.MobileY,
		MobileStyle:  req.MobileStyle,
		HAlign:       req.HAlign,
		VAlign:       req.VAlign,
		ContentData:  req.ContentData,
		Overrides:    req.Overrides,
	}

	contentItem, err := h.contentService.CreateContentItem(c, input)
	if err != nil {
		h.logger.Errorf("Failed to create content item: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Infof("Content item created successfully with ID: %s", contentItem.ID)
	response.Success(c, contentItem, "Content item created successfully", http.StatusCreated)
}

// GetContentItem retrieves a content item by ID
func (h *Handler) GetContentItem(c *gin.Context) {
	itemID := c.Param("id")
	h.logger.Debugf("GetContentItem handler called for item ID: %s", itemID)

	if _, err := uuid.Parse(itemID); err != nil {
		h.logger.Warnf("Invalid item ID format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, "Invalid item ID format")
		return
	}

	contentItem, err := h.contentService.GetContentItem(c, itemID)
	if err != nil {
		h.logger.Warnf("Failed to retrieve content item: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Debugf("Content item retrieved successfully with ID: %s", itemID)
	response.Success(c, contentItem, "Content item retrieved successfully")
}

// GetUserContentItems retrieves all content items for a user
func (h *Handler) GetUserContentItems(c *gin.Context) {
	userID := c.Param("user_id")
	h.logger.Debugf("GetUserContentItems handler called for user ID: %s", userID)

	if _, err := uuid.Parse(userID); err != nil {
		h.logger.Warnf("Invalid user ID format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, "Invalid user ID format")
		return
	}

	contentItems, err := h.contentService.GetUserContentItems(c, userID)
	if err != nil {
		h.logger.Warnf("Failed to retrieve user content items: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Debugf("Retrieved %d content items for user ID: %s", len(contentItems), userID)
	response.Success(c, contentItems, "User content items retrieved successfully")
}

// UpdateContentItem updates a content item
func (h *Handler) UpdateContentItem(c *gin.Context) {
	itemID := c.Param("id")
	h.logger.Infof("UpdateContentItem handler called for item ID: %s", itemID)

	if _, err := uuid.Parse(itemID); err != nil {
		h.logger.Warnf("Invalid item ID format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, "Invalid item ID format")
		return
	}

	var req UpdateContentItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Invalid request format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, err.Error())
		return
	}

	input := service.UpdateContentItemInput{
		Title:        req.Title,
		Href:         req.Href,
		URL:          req.URL,
		MediaType:    req.MediaType,
		DesktopStyle: req.DesktopStyle,
		MobileStyle:  req.MobileStyle,
		HAlign:       req.HAlign,
		VAlign:       req.VAlign,
		ContentData:  req.ContentData,
		Overrides:    req.Overrides,
		IsActive:     req.IsActive,
	}

	contentItem, err := h.contentService.UpdateContentItem(c, itemID, input)
	if err != nil {
		h.logger.Errorf("Failed to update content item: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Infof("Content item updated successfully with ID: %s", itemID)
	response.Success(c, contentItem, "Content item updated successfully")
}

// UpdateContentItemPosition updates the position of a content item
func (h *Handler) UpdateContentItemPosition(c *gin.Context) {
	itemID := c.Param("id")
	h.logger.Infof("UpdateContentItemPosition handler called for item ID: %s", itemID)

	if _, err := uuid.Parse(itemID); err != nil {
		h.logger.Warnf("Invalid item ID format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, "Invalid item ID format")
		return
	}

	var req UpdatePositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Invalid request format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, err.Error())
		return
	}

	input := service.UpdatePositionInput{
		DesktopX: req.DesktopX,
		DesktopY: req.DesktopY,
		MobileX:  req.MobileX,
		MobileY:  req.MobileY,
	}

	contentItem, err := h.contentService.UpdateContentItemPosition(c, itemID, input)
	if err != nil {
		h.logger.Errorf("Failed to update content item position: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Infof("Content item position updated successfully with ID: %s", itemID)
	response.Success(c, contentItem, "Content item position updated successfully")
}

// DeleteContentItem deletes a content item
func (h *Handler) DeleteContentItem(c *gin.Context) {
	itemID := c.Param("id")
	h.logger.Infof("DeleteContentItem handler called for item ID: %s", itemID)

	if _, err := uuid.Parse(itemID); err != nil {
		h.logger.Warnf("Invalid item ID format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, "Invalid item ID format")
		return
	}

	err := h.contentService.DeleteContentItem(c, itemID)
	if err != nil {
		h.logger.Errorf("Failed to delete content item: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Infof("Content item deleted successfully with ID: %s", itemID)
	response.Success(c, nil, "Content item deleted successfully")
}
