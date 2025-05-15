// api/link_metadata/request.go
package link_metadata

// Request types

// FetchLinkMetadataRequest represents the payload for fetching link metadata
type FetchLinkMetadataRequest struct {
	URL string `json:"url" binding:"required"`
}
