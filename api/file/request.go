// api/file/types.go
package file

// PresignedUploadRequest represents a request for a presigned upload URL
type PresignedUploadRequest struct {
	Filename    string `json:"filename" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
	Category    string `json:"category" binding:"required"` // "avatar", "content", "general"
}

// FileUploadResponse represents the response after a successful file upload
type FileUploadResponse struct {
	Key         string `json:"key"`
	URL         string `json:"url"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
	Filename    string `json:"filename"`
	UploadedAt  string `json:"uploaded_at"`
}

// PresignedUploadResponse represents the response for a presigned upload URL
type PresignedUploadResponse struct {
	UploadURL string            `json:"upload_url"`
	Key       string            `json:"key"`
	Fields    map[string]string `json:"fields,omitempty"`
	ExpiresAt string            `json:"expires_at"`
}