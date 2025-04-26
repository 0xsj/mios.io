package service

import (
	"context"

	"github.com/0xsj/gin-sqlc/repository"
)

type ContentService interface {
	CreateContentItem(ctx context.Context, input CreateContentItemInput) (*ContentItemDTO, error)
	GetContentItem(ctx context.Context, itemID string) (*ContentItemDTO, error)
	GetUserContentItems(ctx context.Context, userID string) ([]*ContentItemDTO, error)
	UpdateContentItem(ctx context.Context, itemID string, input UpdateContentItemInput) (*ContentItemDTO, error)
	UpdateContentItemPosition(ctx context.Context, itemID string, input UpdatePositionInput) (*ContentItemDTO, error)
	DeleteContentItem(ctx context.Context, itemID string) error
}

type contentService struct {
	contentRepo repository.ContentRepository
	userRepo    repository.UserRepository
}

type CreateContentItemInput struct {
	UserID      string                 `json:"user_id" binding:"required"`
	ContentID   string                 `json:"content_id" binding:"required"`
	ContentType string                 `json:"content_type" binding:"required"`
	Title       *string                `json:"title"`
	Href        *string                `json:"href"`
	URL         *string                `json:"url"`
	MediaType   *string                `json:"media_type"`
	DesktopX    *int32                 `json:"desktop_x"`
	DesktopY    *int32                 `json:"desktop_y"`
	DesktopStyle *string               `json:"desktop_style"`
	MobileX     *int32                 `json:"mobile_x"`
	MobileY     *int32                 `json:"mobile_y"`
	MobileStyle *string                `json:"mobile_style"`
	HAlign      *string                `json:"halign"`
	VAlign      *string                `json:"valign"`
	ContentData map[string]interface{} `json:"content_data"`
	Overrides   map[string]interface{} `json:"overrides"`
}

type UpdateContentItemInput struct {
	Title       *string                `json:"title"`
	Href        *string                `json:"href"`
	URL         *string                `json:"url"`
	MediaType   *string                `json:"media_type"`
	DesktopStyle *string               `json:"desktop_style"`
	MobileStyle *string                `json:"mobile_style"`
	HAlign      *string                `json:"halign"`
	VAlign      *string                `json:"valign"`
	ContentData map[string]interface{} `json:"content_data"`
	Overrides   map[string]interface{} `json:"overrides"`
	IsActive    *bool                  `json:"is_active"`
}

type UpdatePositionInput struct {
	DesktopX *int32 `json:"desktop_x"`
	DesktopY *int32 `json:"desktop_y"`
	MobileX  *int32 `json:"mobile_x"`
	MobileY  *int32 `json:"mobile_y"`
}

type ContentItemDTO struct {
	ID          string                 `json:"id"`
	UserID      string                 `json:"user_id"`
	ContentID   string                 `json:"content_id"`
	ContentType string                 `json:"content_type"`
	Title       string                 `json:"title,omitempty"`
	Href        string                 `json:"href,omitempty"`
	URL         string                 `json:"url,omitempty"`
	MediaType   string                 `json:"media_type,omitempty"`
	Position    PositionDTO            `json:"position"`
	Style       StyleDTO               `json:"style"`
	HAlign      map[string]string      `json:"halign,omitempty"`
	VAlign      map[string]string      `json:"valign,omitempty"`
	ContentData map[string]interface{} `json:"content_data,omitempty"`
	Overrides   map[string]interface{} `json:"overrides,omitempty"`
	IsActive    bool                   `json:"is_active"`
	CreatedAt   string                 `json:"created_at,omitempty"`
	UpdatedAt   string                 `json:"updated_at,omitempty"`
}

type PositionDTO struct {
	Desktop struct {
		X int32 `json:"x"`
		Y int32 `json:"y"`
	} `json:"desktop"`
	Mobile struct {
		X int32 `json:"x"`
		Y int32 `json:"y"`
	} `json:"mobile"`
}

type StyleDTO struct {
	Desktop string `json:"desktop,omitempty"`
	Mobile  string `json:"mobile,omitempty"`
}

func NewContentService(contentRepo repository.ContentRepository, userRepo repository.UserRepository) ContentService {
	return &contentService{
		contentRepo: contentRepo,
		userRepo:    userRepo,
	}
}

func (s *contentService) CreateContentItem(){}

func (s *contentService) GetContentItem() {}

func (s *contentService) GetUserContentItems() {}

func (s *contentService) UpdateContentItem() {}

func (s *contentService) UpdateContentItemPosition() {}

func (s *contentService) DeleteContentItem() {}

func mapContentItemToDTO() {}