package content

import "github.com/google/uuid"


type CreateContentItemRequest struct {
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

type ContentItemResponse struct {
	ID          uuid.UUID              `json:"id"`
	UserID      uuid.UUID              `json:"user_id"`
	ContentID   string                 `json:"content_id"`
	ContentType string                 `json:"content_type"`
	Title       string                 `json:"title,omitempty"`
	Href        string                 `json:"href,omitempty"`
	URL         string                 `json:"url,omitempty"`
	MediaType   string                 `json:"media_type,omitempty"`
	Position    PositionResponse       `json:"position"`
	Style       StyleResponse          `json:"style"`
	HAlign      string                 `json:"halign,omitempty"`
	VAlign      string                 `json:"valign,omitempty"`
	ContentData map[string]interface{} `json:"content_data,omitempty"`
	Overrides   map[string]interface{} `json:"overrides,omitempty"`
	IsActive    bool                   `json:"is_active"`
	CreatedAt   string                 `json:"created_at,omitempty"`
	UpdatedAt   string                 `json:"updated_at,omitempty"`
}

type PositionResponse struct {
	Desktop struct {
		X int32 `json:"x"`
		Y int32 `json:"y"`
	} `json:"desktop"`
	Mobile struct {
		X int32 `json:"x"`
		Y int32 `json:"y"`
	} `json:"mobile"`
}

type StyleResponse struct {
	Desktop string `json:"desktop,omitempty"`
	Mobile  string `json:"mobile,omitempty"`
}

type UpdateContentItemRequest struct {
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

type UpdatePositionRequest struct {
	DesktopX *int32 `json:"desktop_x"`
	DesktopY *int32 `json:"desktop_y"`
	MobileX  *int32 `json:"mobile_x"`
	MobileY  *int32 `json:"mobile_y"`
}