// api/file/handler.go
package file

import (
	"net/http"
	"strconv"
	"time"

	"github.com/0xsj/mios.io/log"
	"github.com/0xsj/mios.io/pkg/response"
	"github.com/0xsj/mios.io/service"
	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for file operations
type Handler struct {
	fileService service.FileService
	logger      log.Logger
}

// NewHandler creates a new file handler
func NewHandler(fileService service.FileService, logger log.Logger) *Handler {
	return &Handler{
		fileService: fileService,
		logger:      logger,
	}
}

// RegisterRoutes registers file routes on the given router
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	fileGroup := r.Group("/api/files")
	{
		fileGroup.POST("/upload", h.UploadFile)
		fileGroup.POST("/upload/avatar", h.UploadAvatar)
		fileGroup.POST("/upload/content", h.UploadContentMedia)
		fileGroup.POST("/presigned-upload", h.GetPresignedUploadURL)
		fileGroup.DELETE("/:key", h.DeleteFile)
		fileGroup.GET("/:key/url", h.GetFileURL)
	}

	h.logger.Info("File routes registered successfully")
}

// UploadFile handles general file uploads
func (h *Handler) UploadFile(c *gin.Context) {
	h.logger.Info("UploadFile handler called")

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Warn("User ID not found in context")
		response.Error(c, response.ErrUnauthorizedResponse, "User not authenticated")
		return
	}

	// Parse multipart form
	err := c.Request.ParseMultipartForm(32 << 20) // 32 MB max memory
	if err != nil {
		h.logger.Warnf("Failed to parse multipart form: %v", err)
		response.Error(c, response.ErrBadRequestResponse, "Failed to parse form data")
		return
	}

	// Get file from form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		h.logger.Warnf("Failed to get file from form: %v", err)
		response.Error(c, response.ErrBadRequestResponse, "File is required")
		return
	}
	defer file.Close()

	// Get optional category from form
	category := c.PostForm("category")
	if category == "" {
		category = "general"
	}

	// Detect content type from header or file extension
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	input := service.UploadFileInput{
		File:        file,
		Filename:    header.Filename,
		ContentType: contentType,
		Category:    category,
		UserID:      userID.(string),
	}

	result, err := h.fileService.UploadFile(c, input)
	if err != nil {
		h.logger.Errorf("Failed to upload file: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Infof("File uploaded successfully: %s", result.Key)
	response.Success(c, result, "File uploaded successfully", http.StatusCreated)
}

// UploadAvatar handles avatar uploads
func (h *Handler) UploadAvatar(c *gin.Context) {
	h.logger.Info("UploadAvatar handler called")

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Warn("User ID not found in context")
		response.Error(c, response.ErrUnauthorizedResponse, "User not authenticated")
		return
	}

	// Parse multipart form
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB max for avatars
	if err != nil {
		h.logger.Warnf("Failed to parse multipart form: %v", err)
		response.Error(c, response.ErrBadRequestResponse, "Failed to parse form data")
		return
	}

	// Get file from form
	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		h.logger.Warnf("Failed to get avatar file from form: %v", err)
		response.Error(c, response.ErrBadRequestResponse, "Avatar file is required")
		return
	}
	defer file.Close()

	// Detect content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	input := service.UploadFileInput{
		File:        file,
		Filename:    header.Filename,
		ContentType: contentType,
	}

	result, err := h.fileService.UploadUserAvatar(c, userID.(string), input)
	if err != nil {
		h.logger.Errorf("Failed to upload avatar: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Infof("Avatar uploaded successfully: %s", result.Key)
	response.Success(c, result, "Avatar uploaded successfully", http.StatusCreated)
}

// UploadContentMedia handles content media uploads
func (h *Handler) UploadContentMedia(c *gin.Context) {
	h.logger.Info("UploadContentMedia handler called")

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Warn("User ID not found in context")
		response.Error(c, response.ErrUnauthorizedResponse, "User not authenticated")
		return
	}

	// Parse multipart form
	err := c.Request.ParseMultipartForm(50 << 20) // 50 MB max for content media
	if err != nil {
		h.logger.Warnf("Failed to parse multipart form: %v", err)
		response.Error(c, response.ErrBadRequestResponse, "Failed to parse form data")
		return
	}

	// Get file from form
	file, header, err := c.Request.FormFile("media")
	if err != nil {
		h.logger.Warnf("Failed to get media file from form: %v", err)
		response.Error(c, response.ErrBadRequestResponse, "Media file is required")
		return
	}
	defer file.Close()

	// Detect content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	input := service.UploadFileInput{
		File:        file,
		Filename:    header.Filename,
		ContentType: contentType,
	}

	result, err := h.fileService.UploadContentMedia(c, userID.(string), input)
	if err != nil {
		h.logger.Errorf("Failed to upload content media: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Infof("Content media uploaded successfully: %s", result.Key)
	response.Success(c, result, "Content media uploaded successfully", http.StatusCreated)
}

// GetPresignedUploadURL generates a presigned upload URL
func (h *Handler) GetPresignedUploadURL(c *gin.Context) {
	h.logger.Info("GetPresignedUploadURL handler called")

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Warn("User ID not found in context")
		response.Error(c, response.ErrUnauthorizedResponse, "User not authenticated")
		return
	}

	var req PresignedUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Invalid request format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, err.Error())
		return
	}

	input := service.PresignedUploadInput{
		Filename:    req.Filename,
		ContentType: req.ContentType,
		Category:    req.Category,
		UserID:      userID.(string),
	}

	result, err := h.fileService.GetPresignedUploadURL(c, input)
	if err != nil {
		h.logger.Errorf("Failed to generate presigned upload URL: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Infof("Presigned upload URL generated successfully for file: %s", req.Filename)
	response.Success(c, result, "Presigned upload URL generated successfully")
}

// DeleteFile deletes a file
func (h *Handler) DeleteFile(c *gin.Context) {
	h.logger.Info("DeleteFile handler called")

	key := c.Param("key")
	if key == "" {
		h.logger.Warn("File key is required")
		response.Error(c, response.ErrBadRequestResponse, "File key is required")
		return
	}

	// TODO: Add authorization check to ensure user owns the file
	// This would require storing file ownership information

	err := h.fileService.DeleteFile(c, key)
	if err != nil {
		h.logger.Errorf("Failed to delete file: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Infof("File deleted successfully: %s", key)
	response.Success(c, nil, "File deleted successfully")
}

// GetFileURL gets a URL for accessing a file
func (h *Handler) GetFileURL(c *gin.Context) {
	h.logger.Info("GetFileURL handler called")

	key := c.Param("key")
	if key == "" {
		h.logger.Warn("File key is required")
		response.Error(c, response.ErrBadRequestResponse, "File key is required")
		return
	}

	// Parse expires query parameter (optional)
	var expires time.Duration
	if expiresStr := c.Query("expires"); expiresStr != "" {
		expiresHours, err := strconv.Atoi(expiresStr)
		if err != nil {
			h.logger.Warnf("Invalid expires parameter: %v", err)
			response.Error(c, response.ErrBadRequestResponse, "Invalid expires parameter")
			return
		}
		expires = time.Duration(expiresHours) * time.Hour
	}

	url, err := h.fileService.GetFileURL(c, key, expires)
	if err != nil {
		h.logger.Errorf("Failed to get file URL: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	result := map[string]interface{}{
		"key": key,
		"url": url,
	}

	if expires > 0 {
		result["expires_at"] = time.Now().Add(expires).UTC().Format(time.RFC3339)
	}

	h.logger.Infof("File URL generated successfully for key: %s", key)
	response.Success(c, result, "File URL generated successfully")
}