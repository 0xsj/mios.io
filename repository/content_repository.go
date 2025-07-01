package repository

import (
	"context"
	"time"

	db "github.com/0xsj/mios.io/db/sqlc"
	"github.com/0xsj/mios.io/log"
	"github.com/0xsj/mios.io/pkg/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

type ContentRepository interface {
	CreateContentItem(ctx context.Context, params CreateContentItemParams) (*db.ContentItem, error)
	GetContentItem(ctx context.Context, itemID uuid.UUID) (*db.ContentItem, error)
	GetUserContentItems(ctx context.Context, userID uuid.UUID) ([]*db.ContentItem, error)
	UpdateContentItem(ctx context.Context, params UpdateContentItemParams) error
	UpdateContentItemPosition(ctx context.Context, params UpdatePositionParams) error
	DeleteContentItem(ctx context.Context, itemID uuid.UUID) error
}

// CreateContentItemParams matches the service input types
type CreateContentItemParams struct {
	UserID       uuid.UUID
	ContentID    string
	ContentType  string
	Title        *string
	Href         *string
	URL          *string
	MediaType    *string
	DesktopX     *int32
	DesktopY     *int32
	DesktopStyle *string
	MobileX      *int32
	MobileY      *int32
	MobileStyle  *string
	HAlign       *string
	VAlign       *string
	ContentData  pgtype.JSONB
	Overrides    pgtype.JSONB
	IsActive     bool
}

// UpdateContentItemParams matches the service input types
type UpdateContentItemParams struct {
	ItemID       uuid.UUID
	Title        *string
	Href         *string
	URL          *string
	MediaType    *string
	DesktopStyle *string
	MobileStyle  *string
	HAlign       *string
	VAlign       *string
	ContentData  *pgtype.JSONB
	Overrides    *pgtype.JSONB
	IsActive     *bool
}

// UpdatePositionParams matches the service input types
type UpdatePositionParams struct {
	ItemID   uuid.UUID
	DesktopX *int32
	DesktopY *int32
	MobileX  *int32
	MobileY  *int32
}

type SQLContentRepository struct {
	db     *db.Queries
	logger log.Logger
}

func NewContentRepository(db *db.Queries, logger log.Logger) ContentRepository {
	return &SQLContentRepository{
		db:     db,
		logger: logger,
	}
}

func (r *SQLContentRepository) CreateContentItem(ctx context.Context, params CreateContentItemParams) (*db.ContentItem, error) {
	r.logger.Infof("Creating content item with type: %s for user ID: %s", params.ContentType, params.UserID)

	// We directly pass the pointers since the types now match
	sqlcParams := db.CreateContentItemParams{
		UserID:       params.UserID,
		ContentID:    params.ContentID,
		ContentType:  params.ContentType,
		Title:        params.Title,
		Href:         params.Href,
		Url:          params.URL,
		MediaType:    params.MediaType,
		DesktopX:     params.DesktopX,
		DesktopY:     params.DesktopY,
		DesktopStyle: params.DesktopStyle,
		MobileX:      params.MobileX,
		MobileY:      params.MobileY,
		MobileStyle:  params.MobileStyle,
		Halign:       params.HAlign,
		Valign:       params.VAlign,
		ContentData:  params.ContentData,
		Overrides:    params.Overrides,
		IsActive:     &params.IsActive,
	}

	start := time.Now()
	item, err := r.db.CreateContentItem(ctx, sqlcParams)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "content item")
		appErr.Log(r.logger)
		return nil, appErr
	}

	r.logger.Infof("Content item created successfully with ID: %s in %v", item.ItemID, duration)
	return item, nil
}

func (r *SQLContentRepository) GetContentItem(ctx context.Context, itemID uuid.UUID) (*db.ContentItem, error) {
	r.logger.Debugf("Getting content item with ID: %s", itemID)

	start := time.Now()
	item, err := r.db.GetContentItem(ctx, itemID)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "content item")
		appErr.Log(r.logger)
		return nil, appErr
	}

	r.logger.Debugf("Content item retrieved successfully with ID: %s in %v", itemID, duration)
	return item, nil
}

func (r *SQLContentRepository) GetUserContentItems(ctx context.Context, userID uuid.UUID) ([]*db.ContentItem, error) {
	r.logger.Debugf("Getting content items for user ID: %s", userID)

	start := time.Now()
	items, err := r.db.GetUserContentItems(ctx, userID)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "content items")
		appErr.Log(r.logger)
		return nil, appErr
	}

	r.logger.Debugf("Retrieved %d content items for user ID: %s in %v", len(items), userID, duration)
	return items, nil
}

func (r *SQLContentRepository) UpdateContentItem(ctx context.Context, params UpdateContentItemParams) error {
	r.logger.Infof("Updating content item with ID: %s", params.ItemID)

	// Initialize ContentData and Overrides as null by default
	contentData := pgtype.JSONB{Status: pgtype.Null}
	overrides := pgtype.JSONB{Status: pgtype.Null}

	// If we have ContentData, use it
	if params.ContentData != nil {
		contentData = *params.ContentData
	}

	// If we have Overrides, use it
	if params.Overrides != nil {
		overrides = *params.Overrides
	}

	sqlcParams := db.UpdateContentItemParams{
		ItemID:       params.ItemID,
		Title:        params.Title,
		Href:         params.Href,
		Url:          params.URL,
		MediaType:    params.MediaType,
		DesktopStyle: params.DesktopStyle,
		MobileStyle:  params.MobileStyle,
		Halign:       params.HAlign,
		Valign:       params.VAlign,
		ContentData:  contentData,
		Overrides:    overrides,
		IsActive:     params.IsActive,
	}

	start := time.Now()
	err := r.db.UpdateContentItem(ctx, sqlcParams)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "content item update")
		appErr.Log(r.logger)
		return appErr
	}

	r.logger.Infof("Content item updated successfully with ID: %s in %v", params.ItemID, duration)
	return nil
}

func (r *SQLContentRepository) UpdateContentItemPosition(ctx context.Context, params UpdatePositionParams) error {
	r.logger.Infof("Updating position for content item with ID: %s", params.ItemID)

	sqlcParams := db.UpdateContentItemPositionParams{
		ItemID:   params.ItemID,
		DesktopX: params.DesktopX,
		DesktopY: params.DesktopY,
		MobileX:  params.MobileX,
		MobileY:  params.MobileY,
	}

	start := time.Now()
	err := r.db.UpdateContentItemPosition(ctx, sqlcParams)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "content item position update")
		appErr.Log(r.logger)
		return appErr
	}

	r.logger.Infof("Position updated successfully for content item ID: %s in %v", params.ItemID, duration)
	return nil
}

func (r *SQLContentRepository) DeleteContentItem(ctx context.Context, itemID uuid.UUID) error {
	r.logger.Infof("Deleting content item with ID: %s", itemID)

	start := time.Now()
	err := r.db.DeleteContentItem(ctx, itemID)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.HandleDBError(err, "content item deletion")
		appErr.Log(r.logger)
		return appErr
	}

	r.logger.Infof("Content item deleted successfully with ID: %s in %v", itemID, duration)
	return nil
}
