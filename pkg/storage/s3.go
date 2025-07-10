// pkg/storage/s3.go
package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/0xsj/mios.io/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3Storage implements Storage interface for AWS S3
type S3Storage struct {
	client    *s3.Client
	uploader  *manager.Uploader
	bucket    string
	region    string
	cdnDomain string
	logger    log.Logger
}

// S3Config contains configuration for S3 storage
type S3Config struct {
	Region      string
	Bucket      string
	CDNDomain   string
	AccessKeyID string
	SecretKey   string
}

// NewS3Storage creates a new S3 storage instance
func NewS3Storage(cfg S3Config, logger log.Logger) (*S3Storage, error) {
	// Load AWS config
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Override credentials if provided
	if cfg.AccessKeyID != "" && cfg.SecretKey != "" {
		awsCfg.Credentials = aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     cfg.AccessKeyID,
				SecretAccessKey: cfg.SecretKey,
			}, nil
		})
	}

	client := s3.NewFromConfig(awsCfg)
	uploader := manager.NewUploader(client)

	return &S3Storage{
		client:    client,
		uploader:  uploader,
		bucket:    cfg.Bucket,
		region:    cfg.Region,
		cdnDomain: cfg.CDNDomain,
		logger:    logger,
	}, nil
}

func (s *S3Storage) Upload(ctx context.Context, key string, reader io.Reader, opts UploadOptions) (*UploadResult, error) {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(opts.ContentType),
	}

	// Set ACL if provided
	if opts.ACL != "" {
		input.ACL = types.ObjectCannedACL(opts.ACL)
	}

	// Set metadata if provided
	if len(opts.Metadata) > 0 {
		input.Metadata = opts.Metadata
	}

	result, err := s.uploader.Upload(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to S3: %w", err)
	}

	// Get object to retrieve size and ETag
	headOutput, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		s.logger.Warnf("Failed to get object metadata: %v", err)
	}

	var size int64
	var etag string
	if headOutput != nil {
		size = *headOutput.ContentLength
		etag = *headOutput.ETag
	}

	url := result.Location
	if s.cdnDomain != "" {
		url = fmt.Sprintf("%s/%s", s.cdnDomain, key)
	}

	return &UploadResult{
		Key:         key,
		ETag:        etag,
		Size:        size,
		ContentType: opts.ContentType,
		URL:         url,
	}, nil
}

func (s *S3Storage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download from S3: %w", err)
	}

	return result.Body, nil
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete from S3: %w", err)
	}

	return nil
}

func (s *S3Storage) GetURL(ctx context.Context, key string, opts GetURLOptions) (string, error) {
	if opts.CDNDomain != "" {
		return fmt.Sprintf("%s/%s", opts.CDNDomain, key), nil
	}

	if s.cdnDomain != "" {
		return fmt.Sprintf("%s/%s", s.cdnDomain, key), nil
	}

	// Generate S3 URL
	if opts.Expires > 0 {
		// Generate presigned URL for temporary access
		presigner := s3.NewPresignClient(s.client)
		request, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(key),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = opts.Expires
		})
		if err != nil {
			return "", fmt.Errorf("failed to generate presigned URL: %w", err)
		}
		return request.URL, nil
	}

	// Return public S3 URL
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, key), nil
}

// Fix in pkg/storage/s3.go - GetPresignedUploadURL method
func (s *S3Storage) GetPresignedUploadURL(ctx context.Context, key string, opts PresignedUploadOptions) (*PresignedUploadResult, error) {
	presigner := s3.NewPresignClient(s.client)

	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(opts.ContentType),
	}

	if opts.ACL != "" {
		input.ACL = types.ObjectCannedACL(opts.ACL)
	}

	request, err := presigner.PresignPutObject(ctx, input, func(presignOpts *s3.PresignOptions) {
		presignOpts.Expires = opts.Expires
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned upload URL: %w", err) // Fixed: return nil instead of ""
	}

	return &PresignedUploadResult{
		UploadURL: request.URL,
		Key:       key,
		Expires:   time.Now().Add(opts.Expires),
	}, nil
}