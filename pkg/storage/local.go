// pkg/storage/local.go
package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/0xsj/mios.io/log"
)

// LocalStorage implements Storage interface for local file system
type LocalStorage struct {
	basePath  string
	baseURL   string
	logger    log.Logger
}

// NewLocalStorage creates a new local storage instance
func NewLocalStorage(basePath, baseURL string, logger log.Logger) *LocalStorage {
	return &LocalStorage{
		basePath: basePath,
		baseURL:  baseURL,
		logger:   logger,
	}
}

func (l *LocalStorage) Upload(ctx context.Context, key string, reader io.Reader, opts UploadOptions) (*UploadResult, error) {
	fullPath := filepath.Join(l.basePath, key)
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file
	file, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy data with size tracking
	written, err := io.Copy(file, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	// Check size limit
	if opts.MaxSize > 0 && written > opts.MaxSize {
		os.Remove(fullPath)
		return nil, fmt.Errorf("file size %d exceeds maximum %d", written, opts.MaxSize)
	}

	url := fmt.Sprintf("%s/%s", l.baseURL, key)

	return &UploadResult{
		Key:         key,
		Size:        written,
		ContentType: opts.ContentType,
		URL:         url,
		ETag:        fmt.Sprintf("%d-%d", written, time.Now().Unix()),
	}, nil
}

func (l *LocalStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	fullPath := filepath.Join(l.basePath, key)
	return os.Open(fullPath)
}

func (l *LocalStorage) Delete(ctx context.Context, key string) error {
	fullPath := filepath.Join(l.basePath, key)
	return os.Remove(fullPath)
}

func (l *LocalStorage) GetURL(ctx context.Context, key string, opts GetURLOptions) (string, error) {
	if opts.CDNDomain != "" {
		return fmt.Sprintf("%s/%s", opts.CDNDomain, key), nil
	}
	return fmt.Sprintf("%s/%s", l.baseURL, key), nil
}

func (l *LocalStorage) GetPresignedUploadURL(ctx context.Context, key string, opts PresignedUploadOptions) (*PresignedUploadResult, error) {
	// Local storage doesn't support presigned URLs
	return nil, fmt.Errorf("presigned uploads not supported for local storage")
}