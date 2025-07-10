// pkg/storage/interface.go
package storage

import (
	"context"
	"io"
	"time"
)

// Storage defines the interface for file storage operations
type Storage interface {
	Upload(ctx context.Context, key string, reader io.Reader, opts UploadOptions) (*UploadResult, error)
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	GetURL(ctx context.Context, key string, opts GetURLOptions) (string, error)
	GetPresignedUploadURL(ctx context.Context, key string, opts PresignedUploadOptions) (*PresignedUploadResult, error)
}

// UploadOptions contains options for upload operations
type UploadOptions struct {
	ContentType string
	ACL         string
	Metadata    map[string]string
	MaxSize     int64 // Maximum file size in bytes
}

// UploadResult contains the result of an upload operation
type UploadResult struct {
	Key         string
	ETag        string
	Size        int64
	ContentType string
	URL         string
}

// GetURLOptions contains options for getting file URLs
type GetURLOptions struct {
	Expires   time.Duration
	CDNDomain string
}

// PresignedUploadOptions contains options for generating presigned upload URLs
type PresignedUploadOptions struct {
	ContentType string
	MaxSize     int64
	Expires     time.Duration
	ACL         string
}

// PresignedUploadResult contains the result of a presigned upload URL generation
type PresignedUploadResult struct {
	UploadURL string
	Key       string
	Fields    map[string]string
	Expires   time.Time
}