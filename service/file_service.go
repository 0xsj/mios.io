// service/file_service.go
package service

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/0xsj/mios.io/log"
	"github.com/0xsj/mios.io/pkg/errors"
	"github.com/0xsj/mios.io/pkg/storage"
	"github.com/google/uuid"
)

type FileService interface {
	UploadFile(ctx context.Context, input UploadFileInput) (*FileUploadResult, error)
	UploadUserAvatar(ctx context.Context, userID string, input UploadFileInput) (*FileUploadResult, error)
	UploadContentMedia(ctx context.Context, userID string, input UploadFileInput) (*FileUploadResult, error)
	GetPresignedUploadURL(ctx context.Context, input PresignedUploadInput) (*PresignedUploadResult, error)
	DeleteFile(ctx context.Context, key string) error
	GetFileURL(ctx context.Context, key string, expires time.Duration) (string, error)
}

type fileService struct {
	storage storage.Storage
	logger  log.Logger
	config  FileServiceConfig
}

type FileServiceConfig struct {
	MaxFileSize       int64  // in bytes
	MaxAvatarSize     int64  // in bytes
	AllowedImageTypes []string
	AllowedVideoTypes []string
	AllowedFileTypes  []string
	CDNDomain         string
}

type UploadFileInput struct {
	File        io.Reader
	Filename    string
	ContentType string
	Category    string // "avatar", "content", "general"
	UserID      string
}

type FileUploadResult struct {
	Key         string    `json:"key"`
	URL         string    `json:"url"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	Filename    string    `json:"filename"`
	UploadedAt  time.Time `json:"uploaded_at"`
}

type PresignedUploadInput struct {
	Filename    string `json:"filename" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
	Category    string `json:"category" binding:"required"`
	UserID      string `json:"user_id" binding:"required"`
}

type PresignedUploadResult struct {
	UploadURL string            `json:"upload_url"`
	Key       string            `json:"key"`
	Fields    map[string]string `json:"fields,omitempty"`
	ExpiresAt time.Time         `json:"expires_at"`
}

func NewFileService(storage storage.Storage, config FileServiceConfig, logger log.Logger) FileService {
	return &fileService{
		storage: storage,
		logger:  logger,
		config:  config,
	}
}

func (s *fileService) UploadFile(ctx context.Context, input UploadFileInput) (*FileUploadResult, error) {
	// Validate content type
	if err := s.validateContentType(input.ContentType, input.Category); err != nil {
		return nil, err
	}

	// Generate unique key
	key := s.generateFileKey(input.UserID, input.Category, input.Filename)

	// Determine max size based on category
	maxSize := s.config.MaxFileSize
	if input.Category == "avatar" {
		maxSize = s.config.MaxAvatarSize
	}

	// Upload to storage
	uploadOpts := storage.UploadOptions{
		ContentType: input.ContentType,
		MaxSize:     maxSize,
		ACL:         "public-read",
		Metadata: map[string]string{
			"user-id":     input.UserID,
			"category":    input.Category,
			"filename":    input.Filename,
			"uploaded-at": time.Now().UTC().Format(time.RFC3339),
		},
	}

	result, err := s.storage.Upload(ctx, key, input.File, uploadOpts)
	if err != nil {
		s.logger.Errorf("Failed to upload file: %v", err)
		return nil, errors.Wrap(err, "Failed to upload file")
	}

	s.logger.Infof("File uploaded successfully: %s", key)

	return &FileUploadResult{
		Key:         result.Key,
		URL:         result.URL,
		Size:        result.Size,
		ContentType: result.ContentType,
		Filename:    input.Filename,
		UploadedAt:  time.Now().UTC(),
	}, nil
}

func (s *fileService) UploadUserAvatar(ctx context.Context, userID string, input UploadFileInput) (*FileUploadResult, error) {
	input.Category = "avatar"
	input.UserID = userID

	// Validate that it's an image
	if !s.isImageType(input.ContentType) {
		return nil, errors.NewValidationError("Avatar must be an image", nil)
	}

	return s.UploadFile(ctx, input)
}

func (s *fileService) UploadContentMedia(ctx context.Context, userID string, input UploadFileInput) (*FileUploadResult, error) {
	input.Category = "content"
	input.UserID = userID

	return s.UploadFile(ctx, input)
}

func (s *fileService) GetPresignedUploadURL(ctx context.Context, input PresignedUploadInput) (*PresignedUploadResult, error) {
	// Validate content type
	if err := s.validateContentType(input.ContentType, input.Category); err != nil {
		return nil, err
	}

	// Generate unique key
	key := s.generateFileKey(input.UserID, input.Category, input.Filename)

	// Determine max size based on category
	maxSize := s.config.MaxFileSize
	if input.Category == "avatar" {
		maxSize = s.config.MaxAvatarSize
	}

	// Get presigned upload URL
	opts := storage.PresignedUploadOptions{
		ContentType: input.ContentType,
		MaxSize:     maxSize,
		Expires:     15 * time.Minute,
		ACL:         "public-read",
	}

	result, err := s.storage.GetPresignedUploadURL(ctx, key, opts)
	if err != nil {
		s.logger.Errorf("Failed to generate presigned URL: %v", err)
		return nil, errors.Wrap(err, "Failed to generate upload URL")
	}

	return &PresignedUploadResult{
		UploadURL: result.UploadURL,
		Key:       result.Key,
		Fields:    result.Fields,
		ExpiresAt: result.Expires,
	}, nil
}

func (s *fileService) DeleteFile(ctx context.Context, key string) error {
	err := s.storage.Delete(ctx, key)
	if err != nil {
		s.logger.Errorf("Failed to delete file: %v", err)
		return errors.Wrap(err, "Failed to delete file")
	}

	s.logger.Infof("File deleted successfully: %s", key)
	return nil
}

func (s *fileService) GetFileURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	opts := storage.GetURLOptions{
		Expires:   expires,
		CDNDomain: s.config.CDNDomain,
	}

	url, err := s.storage.GetURL(ctx, key, opts)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get file URL")
	}

	return url, nil
}

// Helper methods

func (s *fileService) generateFileKey(userID, category, filename string) string {
	ext := filepath.Ext(filename)
	uniqueID := uuid.New().String()
	timestamp := time.Now().Format("2006/01/02")
	
	return fmt.Sprintf("%s/%s/%s/%s%s", category, userID, timestamp, uniqueID, ext)
}

func (s *fileService) validateContentType(contentType, category string) error {
	switch category {
	case "avatar":
		if !s.isImageType(contentType) {
			return errors.NewValidationError("Avatar must be an image", nil)
		}
		return s.validateImageType(contentType)
	case "content":
		if s.isImageType(contentType) {
			return s.validateImageType(contentType)
		}
		if s.isVideoType(contentType) {
			return s.validateVideoType(contentType)
		}
		return s.validateFileType(contentType)
	default:
		return s.validateFileType(contentType)
	}
}

func (s *fileService) isImageType(contentType string) bool {
	return strings.HasPrefix(contentType, "image/")
}

func (s *fileService) isVideoType(contentType string) bool {
	return strings.HasPrefix(contentType, "video/")
}

func (s *fileService) validateImageType(contentType string) error {
	for _, allowed := range s.config.AllowedImageTypes {
		if contentType == allowed {
			return nil
		}
	}
	return errors.NewValidationError(fmt.Sprintf("Image type %s is not allowed", contentType), nil)
}

func (s *fileService) validateVideoType(contentType string) error {
	for _, allowed := range s.config.AllowedVideoTypes {
		if contentType == allowed {
			return nil
		}
	}
	return errors.NewValidationError(fmt.Sprintf("Video type %s is not allowed", contentType), nil)
}

func (s *fileService) validateFileType(contentType string) error {
	for _, allowed := range s.config.AllowedFileTypes {
		if contentType == allowed {
			return nil
		}
	}
	return errors.NewValidationError(fmt.Sprintf("File type %s is not allowed", contentType), nil)
}