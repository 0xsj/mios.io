package repository

import (
	"context"

	db "github.com/0xsj/gin-sqlc/db/sqlc"
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

type CreateContentItemParams struct {
	UserID      uuid.UUID
	ContentID   string
	ContentType string
	Title       *string
	Href        *string
	URL         *string
	MediaType   *string
	DesktopX    *int32
	DesktopY    *int32
	DesktopStyle *string
	MobileX     *int32
	MobileY     *int32
	MobileStyle *string
	HAlign      *string
	VAlign      *string
	ContentData pgtype.JSONB
	Overrides   pgtype.JSONB
	IsActive    bool
}

type UpdateContentItemParams struct {
	ItemID      uuid.UUID
	Title       *string
	Href        *string
	URL         *string
	MediaType   *string
	DesktopStyle *string
	MobileStyle *string
	HAlign      *string
	VAlign      *string
	ContentData *pgtype.JSONB
	Overrides   *pgtype.JSONB
	IsActive    *bool
}

type UpdatePositionParams struct {
	ItemID      uuid.UUID
	DesktopX    *int32
	DesktopY    *int32
	MobileX     *int32
	MobileY     *int32
}

type SQLContentRepository struct {
	db *db.Queries
}

func NewContentRepository(db *db.Queries) ContentRepository {
	return &SQLContentRepository{
		db: db,
	}
}

func (r *SQLContentRepository) CreateContentItem(ctx context.Context, params CreateContentItemParams)(*db.ContentItem, error){}

func (r *SQLContentRepository) GetContentItem(ctx context.Context, itemID uuid.UUID) (*db.ContentItem, error) {}

func (r *SQLContentRepository) GetUserContentItems(ctx context.Context, userID uuid.UUID)  ([]*db.ContentItem, error) {}


func (r *SQLContentRepository) UpdateContentItem(ctx context.Context,params UpdateContentItemParams) error {}


func (r *SQLContentRepository) UpdateContentItemPosition(ctx context.Context, params UpdatePositionParams) error {
	sqlcParams := db.UpdateContentItemPositionParams{
		ItemID: params.ItemID,
		DesktopX: params.DesktopX,
		DesktopY: params.DesktopY,
		MobileX: params.MobileX,
		MobileY: params.MobileY,
	}

	err := r.db.UpdateContentItemPosition(ctx, sqlcParams)
	if err != nil {
		return ErrDatabase
	}

	return nil
}

func (r *SQLContentRepository) DeleteContentItem(ctx context.Context, itemID uuid.UUID) error{
	err := r.db.DeleteContentItem(ctx, itemID)
	if err != nil {
		return ErrDatabase
	}

	return nil
}