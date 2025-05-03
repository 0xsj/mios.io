package repository

import (
	"context"

	db "github.com/0xsj/gin-sqlc/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
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

type UpdatePositionParams struct {
	ItemID   uuid.UUID
	DesktopX *int32
	DesktopY *int32
	MobileX  *int32
	MobileY  *int32
}

type SQLContentRepository struct {
	db *db.Queries
}

func NewContentRepository(db *db.Queries) ContentRepository {
	return &SQLContentRepository{
		db: db,
	}
}

func (r *SQLContentRepository) CreateContentItem(ctx context.Context, params CreateContentItemParams) (*db.ContentItem, error) {
	var title, href, url, mediaType, desktopStyle, mobileStyle, halign, valign *string
	var desktopX, desktopY, mobileX, mobileY *int32

	if params.Title != nil {
		title = params.Title
	}

	if params.Href != nil {
		href = params.Href
	}

	if params.URL != nil {
		url = params.URL
	}

	if params.MediaType != nil {
		mediaType = params.MediaType
	}

	if params.DesktopX != nil {
		desktopX = params.DesktopX
	}

	if params.DesktopY != nil {
		desktopY = params.DesktopY
	}

	if params.MobileX != nil {
		mobileX = params.MobileX
	}

	if params.MobileY != nil {
		mobileY = params.MobileY
	}

	if params.MobileStyle != nil {
		mobileStyle = params.MobileStyle
	}

	if params.HAlign != nil {
		halign = params.HAlign
	}

	if params.VAlign != nil {
		valign = params.VAlign
	}

	isActive := &params.IsActive

	sqlcParams := db.CreateContentItemParams{
		UserID:       params.UserID,
		ContentID:    params.ContentType,
		ContentType:  params.ContentType,
		Title:        title,
		Href:         href,
		Url:          url,
		MediaType:    mediaType,
		DesktopX:     desktopX,
		DesktopY:     desktopY,
		DesktopStyle: desktopStyle,
		MobileX:      mobileX,
		MobileY:      mobileY,
		MobileStyle:  mobileStyle,
		Halign:       halign,
		Valign:       valign,
		ContentData:  params.ContentData,
		Overrides:    params.Overrides,
		IsActive:     isActive,
	}

	item, err := r.db.CreateContentItem(ctx, sqlcParams)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok {
			if pgErr.Code == "23505" {
				return nil, ErrDuplicateKey
			}

			if pgErr.Code == "23503" {
				return nil, ErrForeignKeyViolation
			}
		}

		return nil, ErrDatabase
	}

	return item, nil
}

func (r *SQLContentRepository) GetContentItem(ctx context.Context, itemID uuid.UUID) (*db.ContentItem, error) {
	item, err := r.db.GetContentItem(ctx, itemID)
	if err != nil {
		return nil, ErrRecordNotFound
	}

	return item, nil
}

func (r *SQLContentRepository) GetUserContentItems(ctx context.Context, userID uuid.UUID) ([]*db.ContentItem, error) {
	items, err := r.db.GetUserContentItems(ctx, userID)
	if err != nil {
		return nil, ErrDatabase
	}

	return items, nil
}

func (r *SQLContentRepository) UpdateContentItem(ctx context.Context, params UpdateContentItemParams) error {
	// Convert our params to SQLC params
	var title, href, url, mediaType, desktopStyle, mobileStyle, halign, valign *string
	var isActive *bool

	if params.Title != nil {
		title = params.Title
	}
	if params.Href != nil {
		href = params.Href
	}
	if params.URL != nil {
		url = params.URL
	}
	if params.MediaType != nil {
		mediaType = params.MediaType
	}
	if params.DesktopStyle != nil {
		desktopStyle = params.DesktopStyle
	}
	if params.MobileStyle != nil {
		mobileStyle = params.MobileStyle
	}
	if params.HAlign != nil {
		halign = params.HAlign
	}
	if params.VAlign != nil {
		valign = params.VAlign
	}
	if params.IsActive != nil {
		isActive = params.IsActive
	}

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
		Title:        title,
		Href:         href,
		Url:          url,
		MediaType:    mediaType,
		DesktopStyle: desktopStyle,
		MobileStyle:  mobileStyle,
		Halign:       halign,
		Valign:       valign,
		ContentData:  contentData,
		Overrides:    overrides,
		IsActive:     isActive,
	}

	err := r.db.UpdateContentItem(ctx, sqlcParams)
	if err != nil {
		return ErrDatabase
	}

	return nil
}

func (r *SQLContentRepository) UpdateContentItemPosition(ctx context.Context, params UpdatePositionParams) error {
	sqlcParams := db.UpdateContentItemPositionParams{
		ItemID:   params.ItemID,
		DesktopX: params.DesktopX,
		DesktopY: params.DesktopY,
		MobileX:  params.MobileX,
		MobileY:  params.MobileY,
	}

	err := r.db.UpdateContentItemPosition(ctx, sqlcParams)
	if err != nil {
		return ErrDatabase
	}

	return nil
}

func (r *SQLContentRepository) DeleteContentItem(ctx context.Context, itemID uuid.UUID) error {
	err := r.db.DeleteContentItem(ctx, itemID)
	if err != nil {
		return ErrDatabase
	}

	return nil
}
