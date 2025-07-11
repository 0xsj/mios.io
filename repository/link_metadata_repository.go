// repository/link_metadata_repository.go
package repository

import (
	"context"
	"time"

	db "github.com/0xsj/mios.io/db/sqlc"
	"github.com/0xsj/mios.io/log"
	"github.com/0xsj/mios.io/pkg/errors"
	"github.com/google/uuid"
)

type LinkMetadataRepository interface {
	CreateLinkMetadata(ctx context.Context, params CreateLinkMetadataParams) (*db.LinkMetadatum, error)
	GetLinkMetadataByURL(ctx context.Context, url string) (*db.LinkMetadatum, error)
	GetLinkMetadataByDomain(ctx context.Context, domain string) ([]*db.LinkMetadatum, error)
	UpdateLinkMetadata(ctx context.Context, params UpdateLinkMetadataParams) (*db.LinkMetadatum, error)
	DeleteLinkMetadata(ctx context.Context, id uuid.UUID) error
}

type CreateLinkMetadataParams struct {
	Domain        string
	URL           string
	Title         *string
	Description   *string
	FaviconURL    *string
	ImageURL      *string
	PlatformName  *string
	PlatformType  *string
	PlatformColor *string
	IsVerified    *bool
}

type UpdateLinkMetadataParams struct {
	URL           string
	Title         *string
	Description   *string
	FaviconURL    *string
	ImageURL      *string
	PlatformName  *string
	PlatformType  *string
	PlatformColor *string
	IsVerified    *bool
}

type SQLCLinkMetadataRepository struct {
	db     *db.Queries
	logger log.Logger
}

func NewLinkMetadataRepository(db *db.Queries, logger log.Logger) LinkMetadataRepository {
	return &SQLCLinkMetadataRepository{
		db:     db,
		logger: logger,
	}
}

func (r *SQLCLinkMetadataRepository) CreateLinkMetadata(ctx context.Context, params CreateLinkMetadataParams) (*db.LinkMetadatum, error) {
	r.logger.Infof("Creating link metadata for URL: %s", params.URL)

	sqlcParams := db.CreateLinkMetadataParams{
		Domain:        params.Domain,
		Url:           params.URL,
		Title:         params.Title,
		Description:   params.Description,
		FaviconUrl:    params.FaviconURL,
		ImageUrl:      params.ImageURL,
		PlatformName:  params.PlatformName,
		PlatformType:  params.PlatformType,
		PlatformColor: params.PlatformColor,
		IsVerified:    params.IsVerified,
	}

	start := time.Now()
	metadata, err := r.db.CreateLinkMetadata(ctx, sqlcParams)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "link metadata")
		appErr.Log(r.logger)
		return nil, appErr
	}

	r.logger.Infof("Link metadata created successfully for URL: %s in %v", params.URL, duration)
	return metadata, nil
}

func (r *SQLCLinkMetadataRepository) GetLinkMetadataByURL(ctx context.Context, url string) (*db.LinkMetadatum, error) {
	r.logger.Debugf("Getting link metadata for URL: %s", url)

	start := time.Now()
	metadata, err := r.db.GetLinkMetadataByURL(ctx, url)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "link metadata")
		appErr.Log(r.logger)
		return nil, appErr
	}

	r.logger.Debugf("Link metadata retrieved successfully for URL: %s in %v", url, duration)
	return metadata, nil
}

func (r *SQLCLinkMetadataRepository) GetLinkMetadataByDomain(ctx context.Context, domain string) ([]*db.LinkMetadatum, error) {
	r.logger.Debugf("Getting link metadata for domain: %s", domain)

	start := time.Now()
	metadataList, err := r.db.GetLinkMetadataByDomain(ctx, domain)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "link metadata")
		appErr.Log(r.logger)
		return nil, appErr
	}

	r.logger.Debugf("Retrieved %d link metadata entries for domain: %s in %v", len(metadataList), domain, duration)
	return metadataList, nil
}

func (r *SQLCLinkMetadataRepository) UpdateLinkMetadata(ctx context.Context, params UpdateLinkMetadataParams) (*db.LinkMetadatum, error) {
	r.logger.Infof("Updating link metadata for URL: %s", params.URL)

	sqlcParams := db.UpdateLinkMetadataParams{
		Url:           params.URL,
		Title:         params.Title,
		Description:   params.Description,
		FaviconUrl:    params.FaviconURL,
		ImageUrl:      params.ImageURL,
		PlatformName:  params.PlatformName,
		PlatformType:  params.PlatformType,
		PlatformColor: params.PlatformColor,
		IsVerified:    params.IsVerified,
	}

	start := time.Now()
	updatedMetadata, err := r.db.UpdateLinkMetadata(ctx, sqlcParams)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "link metadata update")
		appErr.Log(r.logger)
		return nil, appErr
	}

	r.logger.Infof("Link metadata updated successfully for URL: %s in %v", params.URL, duration)
	return updatedMetadata, nil
}

func (r *SQLCLinkMetadataRepository) DeleteLinkMetadata(ctx context.Context, id uuid.UUID) error {
	r.logger.Infof("Deleting link metadata with ID: %s", id)

	start := time.Now()
	err := r.db.DeleteLinkMetadata(ctx, id)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "link metadata deletion")
		appErr.Log(r.logger)
		return appErr
	}

	r.logger.Infof("Link metadata deleted successfully with ID: %s in %v", id, duration)
	return nil
}
