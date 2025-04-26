package repository

import (
	db "github.com/0xsj/gin-sqlc/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

type ContentRepository interface {}

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

