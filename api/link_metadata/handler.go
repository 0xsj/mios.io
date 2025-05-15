// api/link_metadata/handler.go
package link_metadata

import (
	"github.com/0xsj/gin-sqlc/log"
	"github.com/0xsj/gin-sqlc/pkg/response"
	"github.com/0xsj/gin-sqlc/service"
	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for link metadata operations
type Handler struct {
	metadataService service.LinkMetadataService
	logger          log.Logger
}

// NewHandler creates a new link metadata handler
func NewHandler(metadataService service.LinkMetadataService, logger log.Logger) *Handler {
	return &Handler{
		metadataService: metadataService,
		logger:          logger,
	}
}

// RegisterRoutes registers link metadata routes on the given router
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	metadataGroup := r.Group("/api/link-metadata")
	{
		metadataGroup.GET("/url", h.GetLinkMetadata)
		metadataGroup.POST("/fetch", h.FetchLinkMetadata)
		metadataGroup.GET("/platforms", h.ListPlatforms)
	}

	h.logger.Info("Link Metadata routes registered successfully")
}

// GetLinkMetadata retrieves metadata for a URL
func (h *Handler) GetLinkMetadata(c *gin.Context) {
	h.logger.Info("GetLinkMetadata handler called")

	url := c.Query("url")
	if url == "" {
		h.logger.Warn("Missing URL parameter")
		response.Error(c, response.ErrBadRequestResponse, "URL parameter is required")
		return
	}

	metadata, err := h.metadataService.GetMetadata(c, url)
	if err != nil {
		h.logger.Errorf("Failed to get link metadata: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Infof("Link metadata retrieved successfully for URL: %s", url)
	response.Success(c, metadata, "Link metadata retrieved successfully")
}

// FetchLinkMetadata forcibly fetches fresh metadata for a URL
func (h *Handler) FetchLinkMetadata(c *gin.Context) {
	h.logger.Info("FetchLinkMetadata handler called")

	var req FetchLinkMetadataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Invalid request format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, err.Error())
		return
	}

	h.logger.Debugf("Received fetch link metadata request for URL: %s", req.URL)

	metadata, err := h.metadataService.FetchAndStoreMetadata(c, req.URL)
	if err != nil {
		h.logger.Errorf("Failed to fetch link metadata: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Infof("Link metadata fetched successfully for URL: %s", req.URL)
	response.Success(c, metadata, "Link metadata fetched successfully")
}

// ListPlatforms returns a list of known platforms
func (h *Handler) ListPlatforms(c *gin.Context) {
	h.logger.Info("ListPlatforms handler called")

	platforms, err := h.metadataService.ListKnownPlatforms(c)
	if err != nil {
		h.logger.Errorf("Failed to list platforms: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Infof("Platforms listed successfully, found %d platforms", len(platforms))
	response.Success(c, platforms, "Platforms listed successfully")
}
