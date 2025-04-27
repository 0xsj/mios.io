package content

import (
	"net/http"

	"github.com/0xsj/gin-sqlc/api"
	"github.com/0xsj/gin-sqlc/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	contentService service.ContentService
}

func NewHandler(contentService service.ContentService) *Handler {
	return &Handler{
		contentService: contentService,
	}
}

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
}

func (h *Handler) CreateContentItem(c *gin.Context) {
	var req CreateContentItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

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
		api.HandleError(c, err)
		return
	}

	api.RespondWithSuccess(c, contentItem, "Content item created successfully", http.StatusCreated)
}


func (h *Handler) GetContentItem(c *gin.Context) {
	itemID := c.Param("id")
	if _, err := uuid.Parse(itemID); err != nil {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	contentItem, err := h.contentService.GetContentItem(c, itemID)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	api.RespondWithSuccess(c, contentItem, "Content item retrieved successfully")
}


func (h *Handler) GetUserContentItems(c *gin.Context) {
	userID := c.Param("user_id")
	if _, err := uuid.Parse(userID); err != nil {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	contentItems, err := h.contentService.GetUserContentItems(c, userID)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	api.RespondWithSuccess(c, contentItems, "User content items retrieved successfully")
}

func (h *Handler) UpdateContentItem(c *gin.Context) {
	itemID := c.Param("id")
	if _, err := uuid.Parse(itemID); err != nil {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	var req UpdateContentItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.HandleError(c, api.ErrInvalidInput)
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
		api.HandleError(c, err)
		return
	}

	api.RespondWithSuccess(c, contentItem, "Content item updated successfully")
}

func (h *Handler) UpdateContentItemPosition(c *gin.Context) {
	itemID := c.Param("id")
	if _, err := uuid.Parse(itemID); err != nil {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	var req UpdatePositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.HandleError(c, api.ErrInvalidInput)
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
		api.HandleError(c, err)
		return
	}

	api.RespondWithSuccess(c, contentItem, "Content item position updated successfully")
}

func (h *Handler) DeleteContentItem(c *gin.Context) {
	itemID := c.Param("id")
	if _, err := uuid.Parse(itemID); err != nil {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	err := h.contentService.DeleteContentItem(c, itemID)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	api.RespondWithSuccess(c, nil, "Content item deleted successfully")
}