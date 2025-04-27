package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/0xsj/gin-sqlc/api"
	db "github.com/0xsj/gin-sqlc/db/sqlc"
	"github.com/0xsj/gin-sqlc/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
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

func (s *contentService) CreateContentItem(ctx context.Context, input CreateContentItemInput) (*ContentItemDTO, error){
	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return nil, api.ErrInvalidInput
	}

	_, err = s.userRepo.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, api.ErrNotFound
		}
		return nil, api.ErrInternalServer
	}

	var contentData, overrides pgtype.JSONB
	if len(input.ContentData) > 0 {
		contentData.Status = pgtype.Present
		contentData.Bytes, err = json.Marshal(input.ContentData)
		if err != nil {
			return nil, api.ErrInvalidInput
		}
	} else {
		contentData.Status = pgtype.Null
	}

	if len(input.Overrides) > 0 {
		overrides.Status = pgtype.Present
		overrides.Bytes , err = json.Marshal(input.Overrides)
		if err != nil {
			return nil, api.ErrInvalidInput
		}
	} else {
		overrides.Status = pgtype.Null
	}

	params := repository.CreateContentItemParams{
		UserID:       userID,
		ContentID:    input.ContentID,
		ContentType:  input.ContentType,
		Title:        input.Title,
		Href:         input.Href,
		URL:          input.URL,
		MediaType:    input.MediaType,
		DesktopX:     input.DesktopX,
		DesktopY:     input.DesktopY,
		DesktopStyle: input.DesktopStyle,
		MobileX:      input.MobileX,
		MobileY:      input.MobileY,
		MobileStyle:  input.MobileStyle,
		HAlign:       input.HAlign,
		VAlign:       input.VAlign,
		ContentData:  contentData,
		Overrides:    overrides,
		IsActive:     true, 
	}	
	
	contentItem, err := s.contentRepo.CreateContentItem(ctx, params)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateKey) {
			return nil, api.ErrDuplicateEntry
		}
		return nil, api.ErrInternalServer
	}
	return mapContentItemToDTO(contentItem), nil
}

func (s *contentService) GetContentItem(ctx context.Context, itemIDStr string) (*ContentItemDTO, error) {
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		return nil, api.ErrInvalidInput
	}

	contentItem, err := s.contentRepo.GetContentItem(ctx, itemID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, api.ErrNotFound
		}
		return nil, api.ErrInternalServer
	}
	return mapContentItemToDTO(contentItem), nil
}

func (s *contentService) GetUserContentItems(ctx context.Context, userIDStr string)([]*ContentItemDTO, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, api.ErrInvalidInput
	}

	_, err = s.userRepo.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, api.ErrNotFound
		}
		return nil, api.ErrInternalServer
	}

	contentItems, err := s.contentRepo.GetUserContentItems(ctx, userID)
	if err != nil {
		return nil, api.ErrInternalServer
	}

	dtos := make([]*ContentItemDTO, len(contentItems))
	for i, item := range contentItems {
		dtos[i] = mapContentItemToDTO(item)
	}

	return dtos, nil

}

func (s *contentService) UpdateContentItem(ctx context.Context, itemIDStr string, input UpdateContentItemInput) (*ContentItemDTO, error) {
    itemID, err := uuid.Parse(itemIDStr)
    if err != nil {
        return nil, api.ErrInvalidInput
    }

    _, err = s.contentRepo.GetContentItem(ctx, itemID)
    if err != nil {
        if errors.Is(err, repository.ErrRecordNotFound) {
            return nil, api.ErrNotFound
        }
        return nil, api.ErrInternalServer
    }

    var contentData, overrides *pgtype.JSONB
    
    if len(input.ContentData) > 0 {
		var cData pgtype.JSONB
		cData.Status = pgtype.Present
		cData.Bytes, err = json.Marshal(input.ContentData)
		if err != nil {
			return nil, api.ErrInvalidInput
		}
		contentData = &cData
	}
	
    if len(input.Overrides) > 0 {
		var oData pgtype.JSONB
		oData.Status = pgtype.Present
		oData.Bytes, err = json.Marshal(input.Overrides)
		if err != nil {
			return nil, api.ErrInvalidInput
		}
		overrides = &oData
	}

    params := repository.UpdateContentItemParams{
        ItemID:       itemID,
        Title:        input.Title,
        Href:         input.Href,
        URL:          input.URL,
        MediaType:    input.MediaType,
        DesktopStyle: input.DesktopStyle,
        MobileStyle:  input.MobileStyle,
        HAlign:       input.HAlign,
        VAlign:       input.VAlign,
        ContentData:  contentData,
        Overrides:    overrides,
        IsActive:     input.IsActive,
    }

    err = s.contentRepo.UpdateContentItem(ctx, params)
    if err != nil {
        return nil, api.ErrInternalServer
    }

    updatedItem, err := s.contentRepo.GetContentItem(ctx, itemID)
    if err != nil {
        return nil, api.ErrInternalServer
    }

    return mapContentItemToDTO(updatedItem), nil
}

func (s *contentService) UpdateContentItemPosition(ctx context.Context, itemIDStr string, input UpdatePositionInput) (*ContentItemDTO, error) {
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		return nil, api.ErrInvalidInput
	}

	_, err = s.contentRepo.GetContentItem(ctx, itemID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, api.ErrNotFound
		}
		return nil, api.ErrInternalServer
	}

	params := repository.UpdatePositionParams{
		ItemID:    itemID,
		DesktopX:  input.DesktopX,
		DesktopY:  input.DesktopY,
		MobileX:   input.MobileX,
		MobileY:   input.MobileY,
	}

	err = s.contentRepo.UpdateContentItemPosition(ctx, params)
	if err != nil {
		return nil, api.ErrInternalServer
	}

	updatedItem, err := s.contentRepo.GetContentItem(ctx, itemID)
	if err != nil {
		return nil, api.ErrInternalServer
	}

	return mapContentItemToDTO(updatedItem), nil
}

func (s *contentService) DeleteContentItem(ctx context.Context, itemIDStr string) error {
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		return api.ErrInvalidInput
	}

	_, err = s.contentRepo.GetContentItem(ctx, itemID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return api.ErrNotFound
		}
		return api.ErrInternalServer
	}

	err = s.contentRepo.DeleteContentItem(ctx, itemID)
	if err != nil {
		return api.ErrInternalServer
	}

	return nil
}

func mapContentItemToDTO(item *db.ContentItem) *ContentItemDTO {
	dto := &ContentItemDTO{
		ID:          item.ItemID.String(),
		UserID:      item.UserID.String(),
		ContentID:   item.ContentID,
		ContentType: item.ContentType,
		IsActive:    item.IsActive != nil && *item.IsActive,
		Position: PositionDTO{
			Desktop: struct {
				X int32 `json:"x"`
				Y int32 `json:"y"`
			}{
				X: 0,
				Y: 0,
			},
			Mobile: struct {
				X int32 `json:"x"`
				Y int32 `json:"y"`
			}{
				X: 0,
				Y: 0,
			},
		},
		Style: StyleDTO{},
	}

	if item.Title != nil {
		dto.Title = *item.Title
	}

	if item.Href != nil {
		dto.Href = *item.Href
	}

	if item.Url != nil {
		dto.URL = *item.Url
	}

	if item.MediaType != nil {
		dto.MediaType = *item.MediaType
	}

	if item.DesktopX != nil {
		dto.Position.Desktop.X = *item.DesktopX
	}

	if item.DesktopY != nil {
		dto.Position.Desktop.Y = *item.DesktopY
	}

	if item.DesktopStyle != nil {
		dto.Style.Desktop = *item.DesktopStyle
	}

	if item.MobileX != nil {
		dto.Position.Mobile.X = *item.MobileX
	}

	if item.MobileY != nil {
		dto.Position.Mobile.Y = *item.MobileY
	}

	if item.MobileStyle != nil {
		dto.Style.Mobile = *item.MobileStyle
	}

	if item.Halign != nil || item.Valign != nil {
		dto.HAlign = make(map[string]string)
		dto.VAlign = make(map[string]string)

		if item.Halign != nil {
			dto.HAlign["default"] = *item.Halign
		}

		if item.Valign != nil {
			dto.VAlign["default"] = *item.Valign
		}
	}

	if item.ContentData.Status == pgtype.Present {
		var contentData map[string]interface{}
		if err := json.Unmarshal(item.ContentData.Bytes, &contentData); err == nil {
			dto.ContentData = contentData
		}
	}

	if item.Overrides.Status == pgtype.Present {
		var overrides map[string]interface{}
		if err := json.Unmarshal(item.Overrides.Bytes, &overrides); err == nil {
			dto.Overrides = overrides
		}
	}

	if item.CreatedAt != nil {
		dto.CreatedAt = item.CreatedAt.Format(time.RFC3339)
	}

	if item.UpdatedAt != nil {
		dto.UpdatedAt = item.UpdatedAt.Format(time.RFC3339)
	}

	return dto
}