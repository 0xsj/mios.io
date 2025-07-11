package service

import (
	"context"
	"encoding/json"
	"time"

	db "github.com/0xsj/mios.io/db/sqlc"
	"github.com/0xsj/mios.io/log"
	"github.com/0xsj/mios.io/pkg/errors"
	"github.com/0xsj/mios.io/repository"
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
	logger      log.Logger
}

type CreateContentItemInput struct {
	UserID       string                 `json:"user_id" binding:"required"`
	ContentID    string                 `json:"content_id" binding:"required"`
	ContentType  string                 `json:"content_type" binding:"required"`
	Title        *string                `json:"title"`
	Href         *string                `json:"href"`
	URL          *string                `json:"url"`
	MediaType    *string                `json:"media_type"`
	DesktopX     *int32                 `json:"desktop_x"`
	DesktopY     *int32                 `json:"desktop_y"`
	DesktopStyle *string                `json:"desktop_style"`
	MobileX      *int32                 `json:"mobile_x"`
	MobileY      *int32                 `json:"mobile_y"`
	MobileStyle  *string                `json:"mobile_style"`
	HAlign       *string                `json:"halign"`
	VAlign       *string                `json:"valign"`
	ContentData  map[string]interface{} `json:"content_data"`
	Overrides    map[string]interface{} `json:"overrides"`
}

type UpdateContentItemInput struct {
	Title        *string                `json:"title"`
	Href         *string                `json:"href"`
	URL          *string                `json:"url"`
	MediaType    *string                `json:"media_type"`
	DesktopStyle *string                `json:"desktop_style"`
	MobileStyle  *string                `json:"mobile_style"`
	HAlign       *string                `json:"halign"`
	VAlign       *string                `json:"valign"`
	ContentData  map[string]interface{} `json:"content_data"`
	Overrides    map[string]interface{} `json:"overrides"`
	IsActive     *bool                  `json:"is_active"`
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

func NewContentService(contentRepo repository.ContentRepository, userRepo repository.UserRepository, logger log.Logger) ContentService {
	return &contentService{
		contentRepo: contentRepo,
		userRepo:    userRepo,
		logger:      logger,
	}
}

func (s *contentService) CreateContentItem(ctx context.Context, input CreateContentItemInput) (*ContentItemDTO, error) {
	s.logger.Infof("Creating content item for user ID: %s with content type: %s", input.UserID, input.ContentType)

	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		s.logger.Warnf("Invalid user ID format: %v", err)
		return nil, errors.NewBadRequestError("Invalid user ID format", err)
	}

	// Verify user exists
	_, err = s.userRepo.GetUser(ctx, userID)
	if err != nil {
		if errors.IsNotFound(err) {
			s.logger.Warnf("User not found with ID: %s", input.UserID)
			return nil, errors.NewNotFoundError("User not found", err)
		}
		s.logger.Errorf("Error retrieving user: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve user")
	}

	// Process JSON data
	var contentData, overrides pgtype.JSONB
	if len(input.ContentData) > 0 {
		contentData.Status = pgtype.Present
		contentData.Bytes, err = json.Marshal(input.ContentData)
		if err != nil {
			s.logger.Warnf("Failed to marshal content data: %v", err)
			return nil, errors.NewValidationError("Invalid content data format", err)
		}
	} else {
		contentData.Status = pgtype.Null
	}

	if len(input.Overrides) > 0 {
		overrides.Status = pgtype.Present
		overrides.Bytes, err = json.Marshal(input.Overrides)
		if err != nil {
			s.logger.Warnf("Failed to marshal overrides: %v", err)
			return nil, errors.NewValidationError("Invalid overrides format", err)
		}
	} else {
		overrides.Status = pgtype.Null
	}

	// Create content item
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
		if errors.IsConflict(err) {
			s.logger.Warnf("Content item already exists: %v", err)
			return nil, errors.NewConflictError("Content item already exists", err)
		}
		s.logger.Errorf("Failed to create content item: %v", err)
		return nil, errors.Wrap(err, "Failed to create content item")
	}

	s.logger.Infof("Content item created successfully with ID: %s", contentItem.ItemID)
	return mapContentItemToDTO(contentItem), nil
}

func (s *contentService) GetContentItem(ctx context.Context, itemIDStr string) (*ContentItemDTO, error) {
	s.logger.Debugf("Getting content item with ID: %s", itemIDStr)

	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		s.logger.Warnf("Invalid item ID format: %v", err)
		return nil, errors.NewBadRequestError("Invalid item ID format", err)
	}

	contentItem, err := s.contentRepo.GetContentItem(ctx, itemID)
	if err != nil {
		if errors.IsNotFound(err) {
			s.logger.Infof("Content item not found with ID: %s", itemIDStr)
			return nil, errors.NewNotFoundError("Content item not found", err)
		}
		s.logger.Errorf("Error retrieving content item: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve content item")
	}

	s.logger.Debugf("Content item retrieved successfully with ID: %s", itemIDStr)
	return mapContentItemToDTO(contentItem), nil
}

func (s *contentService) GetUserContentItems(ctx context.Context, userIDStr string) ([]*ContentItemDTO, error) {
	s.logger.Debugf("Getting content items for user ID: %s", userIDStr)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		s.logger.Warnf("Invalid user ID format: %v", err)
		return nil, errors.NewBadRequestError("Invalid user ID format", err)
	}

	// Verify user exists
	_, err = s.userRepo.GetUser(ctx, userID)
	if err != nil {
		if errors.IsNotFound(err) {
			s.logger.Infof("User not found with ID: %s", userIDStr)
			return nil, errors.NewNotFoundError("User not found", err)
		}
		s.logger.Errorf("Error retrieving user: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve user")
	}

	contentItems, err := s.contentRepo.GetUserContentItems(ctx, userID)
	if err != nil {
		s.logger.Errorf("Failed to retrieve content items: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve content items")
	}

	dtos := make([]*ContentItemDTO, len(contentItems))
	for i, item := range contentItems {
		dtos[i] = mapContentItemToDTO(item)
	}

	s.logger.Debugf("Retrieved %d content items for user ID: %s", len(dtos), userIDStr)
	return dtos, nil
}

func (s *contentService) UpdateContentItem(ctx context.Context, itemIDStr string, input UpdateContentItemInput) (*ContentItemDTO, error) {
	s.logger.Infof("Updating content item with ID: %s", itemIDStr)

	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		s.logger.Warnf("Invalid item ID format: %v", err)
		return nil, errors.NewBadRequestError("Invalid item ID format", err)
	}

	// Verify content item exists
	_, err = s.contentRepo.GetContentItem(ctx, itemID)
	if err != nil {
		if errors.IsNotFound(err) {
			s.logger.Infof("Content item not found with ID: %s", itemIDStr)
			return nil, errors.NewNotFoundError("Content item not found", err)
		}
		s.logger.Errorf("Error retrieving content item: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve content item")
	}

	// Process JSON data
	var contentData, overrides *pgtype.JSONB

	if len(input.ContentData) > 0 {
		var cData pgtype.JSONB
		cData.Status = pgtype.Present
		cData.Bytes, err = json.Marshal(input.ContentData)
		if err != nil {
			s.logger.Warnf("Failed to marshal content data: %v", err)
			return nil, errors.NewValidationError("Invalid content data format", err)
		}
		contentData = &cData
	}

	if len(input.Overrides) > 0 {
		var oData pgtype.JSONB
		oData.Status = pgtype.Present
		oData.Bytes, err = json.Marshal(input.Overrides)
		if err != nil {
			s.logger.Warnf("Failed to marshal overrides: %v", err)
			return nil, errors.NewValidationError("Invalid overrides format", err)
		}
		overrides = &oData
	}

	// Update content item
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
		s.logger.Errorf("Failed to update content item: %v", err)
		return nil, errors.Wrap(err, "Failed to update content item")
	}

	// Retrieve updated item
	updatedItem, err := s.contentRepo.GetContentItem(ctx, itemID)
	if err != nil {
		s.logger.Errorf("Failed to retrieve updated content item: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve updated content item")
	}

	s.logger.Infof("Content item updated successfully with ID: %s", itemIDStr)
	return mapContentItemToDTO(updatedItem), nil
}

func (s *contentService) UpdateContentItemPosition(ctx context.Context, itemIDStr string, input UpdatePositionInput) (*ContentItemDTO, error) {
	s.logger.Infof("Updating position for content item with ID: %s", itemIDStr)

	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		s.logger.Warnf("Invalid item ID format: %v", err)
		return nil, errors.NewBadRequestError("Invalid item ID format", err)
	}

	// Verify content item exists
	_, err = s.contentRepo.GetContentItem(ctx, itemID)
	if err != nil {
		if errors.IsNotFound(err) {
			s.logger.Infof("Content item not found with ID: %s", itemIDStr)
			return nil, errors.NewNotFoundError("Content item not found", err)
		}
		s.logger.Errorf("Error retrieving content item: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve content item")
	}

	// Update position
	params := repository.UpdatePositionParams{
		ItemID:   itemID,
		DesktopX: input.DesktopX,
		DesktopY: input.DesktopY,
		MobileX:  input.MobileX,
		MobileY:  input.MobileY,
	}

	err = s.contentRepo.UpdateContentItemPosition(ctx, params)
	if err != nil {
		s.logger.Errorf("Failed to update content item position: %v", err)
		return nil, errors.Wrap(err, "Failed to update content item position")
	}

	// Retrieve updated item
	updatedItem, err := s.contentRepo.GetContentItem(ctx, itemID)
	if err != nil {
		s.logger.Errorf("Failed to retrieve updated content item: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve updated content item")
	}

	s.logger.Infof("Position updated successfully for content item ID: %s", itemIDStr)
	return mapContentItemToDTO(updatedItem), nil
}

func (s *contentService) DeleteContentItem(ctx context.Context, itemIDStr string) error {
	s.logger.Infof("Deleting content item with ID: %s", itemIDStr)

	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		s.logger.Warnf("Invalid item ID format: %v", err)
		return errors.NewBadRequestError("Invalid item ID format", err)
	}

	// Verify content item exists
	_, err = s.contentRepo.GetContentItem(ctx, itemID)
	if err != nil {
		if errors.IsNotFound(err) {
			s.logger.Infof("Content item not found with ID: %s", itemIDStr)
			return errors.NewNotFoundError("Content item not found", err)
		}
		s.logger.Errorf("Error retrieving content item: %v", err)
		return errors.Wrap(err, "Failed to retrieve content item")
	}

	// Delete content item
	err = s.contentRepo.DeleteContentItem(ctx, itemID)
	if err != nil {
		s.logger.Errorf("Failed to delete content item: %v", err)
		return errors.Wrap(err, "Failed to delete content item")
	}

	s.logger.Infof("Content item deleted successfully with ID: %s", itemIDStr)
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
