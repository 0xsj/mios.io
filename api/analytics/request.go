package analytics

type RecordClickRequest struct {
	ItemID    string `json:"item_id" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	Referrer  string `json:"referrer"`
}

type RecordPageViewRequest struct {
	ProfileID string `json:"profile_id" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	Referrer  string `json:"referrer"`
}

type TimeRangeRequest struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	Limit     int    `json:"limit"`
}
