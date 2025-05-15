// service/link_metadata_service.go
package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	db "github.com/0xsj/gin-sqlc/db/sqlc"
	"github.com/0xsj/gin-sqlc/log"
	"github.com/0xsj/gin-sqlc/pkg/errors"
	"github.com/0xsj/gin-sqlc/repository"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type LinkMetadataService interface {
	GetMetadata(ctx context.Context, urlString string) (*LinkMetadataDTO, error)
	FetchAndStoreMetadata(ctx context.Context, urlString string) (*LinkMetadataDTO, error)
	IsKnownPlatform(domain string) bool
	GetPlatformInfo(domain string) *PlatformInfo
	ListKnownPlatforms(ctx context.Context) ([]*PlatformInfo, error)
}

type PlatformInfo struct {
	Domain      string `json:"domain"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Color       string `json:"color"`
	Icon        string `json:"icon"`
	URLTemplate string `json:"url_template"`
}

type LinkMetadataDTO struct {
	ID           string `json:"id"`
	URL          string `json:"url"`
	Domain       string `json:"domain"`
	Title        string `json:"title,omitempty"`
	Description  string `json:"description,omitempty"`
	FaviconURL   string `json:"favicon_url,omitempty"`
	ImageURL     string `json:"image_url,omitempty"`
	PlatformName string `json:"platform_name,omitempty"`
	PlatformType string `json:"platform_type,omitempty"`
	PlatformColor string `json:"platform_color,omitempty"`
	IsVerified   bool   `json:"is_verified"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

// PlatformRegistry defines known platforms and their metadata
var PlatformRegistry = map[string]PlatformInfo{
	"instagram.com": {
		Domain:      "instagram.com",
		Name:        "Instagram",
		Type:        "social",
		Color:       "#E1306C",
		Icon:        "/assets/icons/instagram.svg",
		URLTemplate: "https://instagram.com/{username}",
	},
	"twitter.com": {
		Domain:      "twitter.com",
		Name:        "Twitter",
		Type:        "social",
		Color:       "#1DA1F2",
		Icon:        "/assets/icons/twitter.svg", 
		URLTemplate: "https://twitter.com/{username}",
	},
	"x.com": {
		Domain:      "x.com",
		Name:        "X",
		Type:        "social",
		Color:       "#000000",
		Icon:        "/assets/icons/x.svg",
		URLTemplate: "https://x.com/{username}",
	},
	"github.com": {
		Domain:      "github.com",
		Name:        "GitHub",
		Type:        "dev",
		Color:       "#333333",
		Icon:        "/assets/icons/github.svg",
		URLTemplate: "https://github.com/{username}",
	},
	"linkedin.com": {
		Domain:      "linkedin.com",
		Name:        "LinkedIn",
		Type:        "professional",
		Color:       "#0A66C2",
		Icon:        "/assets/icons/linkedin.svg",
		URLTemplate: "https://linkedin.com/in/{username}",
	},
	"youtube.com": {
		Domain:      "youtube.com",
		Name:        "YouTube",
		Type:        "video",
		Color:       "#FF0000",
		Icon:        "/assets/icons/youtube.svg",
		URLTemplate: "https://youtube.com/{channel}",
	},
	"tiktok.com": {
		Domain:      "tiktok.com",
		Name:        "TikTok",
		Type:        "video",
		Color:       "#000000",
		Icon:        "/assets/icons/tiktok.svg",
		URLTemplate: "https://tiktok.com/@{username}",
	},
	"facebook.com": {
		Domain:      "facebook.com",
		Name:        "Facebook",
		Type:        "social",
		Color:       "#1877F2",
		Icon:        "/assets/icons/facebook.svg",
		URLTemplate: "https://facebook.com/{username}",
	},
	"spotify.com": {
		Domain:      "spotify.com",
		Name:        "Spotify",
		Type:        "music",
		Color:       "#1DB954",
		Icon:        "/assets/icons/spotify.svg",
		URLTemplate: "https://open.spotify.com/user/{username}",
	},
	"twitch.tv": {
		Domain:      "twitch.tv",
		Name:        "Twitch",
		Type:        "streaming",
		Color:       "#9146FF",
		Icon:        "/assets/icons/twitch.svg",
		URLTemplate: "https://twitch.tv/{username}",
	},
}

type linkMetadataService struct {
	repo   repository.LinkMetadataRepository
	logger log.Logger
	client *http.Client
}

func NewLinkMetadataService(repo repository.LinkMetadataRepository, logger log.Logger) LinkMetadataService {
	return &linkMetadataService{
		repo:   repo,
		logger: logger,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *linkMetadataService) GetMetadata(ctx context.Context, urlString string) (*LinkMetadataDTO, error) {
	s.logger.Debugf("Getting metadata for URL: %s", urlString)
	
	// Clean and normalize the URL
	normalizedURL, err := normalizeURL(urlString)
	if err != nil {
		s.logger.Warnf("Invalid URL format: %v", err)
		return nil, errors.NewValidationError("Invalid URL format", err)
	}

	// Check if we already have metadata for this URL
	metadata, err := s.repo.GetLinkMetadataByURL(ctx, normalizedURL)
	
	// If not found or error other than NotFound, fetch fresh metadata
	if err != nil {
		if errors.IsNotFound(err) {
			s.logger.Debugf("Metadata not found for URL: %s, fetching fresh data", normalizedURL)
			return s.FetchAndStoreMetadata(ctx, normalizedURL)
		}
		s.logger.Errorf("Error retrieving metadata: %v", err)
		return nil, errors.Wrap(err, "Failed to retrieve metadata")
	}
	
	// If metadata is older than a week, refresh it asynchronously
	if metadata.UpdatedAt != nil && time.Since(*metadata.UpdatedAt) > 7*24*time.Hour {
		s.logger.Debugf("Metadata for URL %s is older than a week, refreshing asynchronously", normalizedURL)
		go func() {
			bgCtx := context.Background()
			_, err := s.FetchAndStoreMetadata(bgCtx, normalizedURL)
			if err != nil {
				s.logger.Warnf("Failed to refresh metadata for URL %s: %v", normalizedURL, err)
			}
		}()
	}
	
	return mapLinkMetadataToDTO(metadata), nil
}

func (s *linkMetadataService) FetchAndStoreMetadata(ctx context.Context, urlString string) (*LinkMetadataDTO, error) {
	s.logger.Infof("Fetching metadata for URL: %s", urlString)
	
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		s.logger.Warnf("Failed to parse URL: %v", err)
		return nil, errors.NewValidationError("Invalid URL format", err)
	}
	
	domain := parsedURL.Hostname()
	
	// Check if it's a known platform
	var platformName, platformType, platformColor *string
	if platform, found := PlatformRegistry[domain]; found {
		s.logger.Debugf("URL %s matches known platform: %s", urlString, platform.Name)
		platformName = &platform.Name
		platformType = &platform.Type
		platformColor = &platform.Color
	}
	
	// Fetch page content
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlString, nil)
	if err != nil {
		s.logger.Warnf("Failed to create request: %v", err)
		return nil, errors.NewExternalServiceError("Failed to create request", err)
	}
	
	req.Header.Set("User-Agent", "Link Metadata Service 1.0")
	
	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Warnf("Failed to fetch URL: %v", err)
		
		// Store minimal information if we can't fetch
		createParams := repository.CreateLinkMetadataParams{
			Domain:       domain,
			URL:          urlString,
			PlatformName: platformName,
			PlatformType: platformType,
			PlatformColor: platformColor,
			IsVerified:   nil,
		}
		
		metadata, err := s.repo.CreateLinkMetadata(ctx, createParams)
		if err != nil {
			s.logger.Errorf("Failed to store minimal metadata: %v", err)
			return nil, errors.Wrap(err, "Failed to store metadata")
		}
		
		return mapLinkMetadataToDTO(metadata), nil
	}
	defer resp.Body.Close()
	
	// Parse HTML to extract metadata
	doc, err := html.Parse(resp.Body)
	if err != nil {
		s.logger.Warnf("Failed to parse HTML: %v", err)
		return nil, errors.NewExternalServiceError("Failed to parse page content", err)
	}
	
	// Extract metadata
	var (
		title       *string
		description *string
		faviconURL  *string
		imageURL    *string
	)
	
	// Extract metadata from HTML
	metadata := extractMetadata(doc, urlString)
	
	if metadata.Title != "" {
		title = &metadata.Title
	}
	
	if metadata.Description != "" {
		description = &metadata.Description
	}
	
	if metadata.FaviconURL != "" {
		faviconURL = &metadata.FaviconURL
	} else {
		// Try default favicon location
		defaultFavicon := fmt.Sprintf("%s://%s/favicon.ico", parsedURL.Scheme, domain)
		faviconURL = &defaultFavicon
	}
	
	if metadata.ImageURL != "" {
		imageURL = &metadata.ImageURL
	}
	
	// Check if we already have this URL in the database
	existingMetadata, err := s.repo.GetLinkMetadataByURL(ctx, urlString)
	if err == nil && existingMetadata != nil {
		// Update existing metadata
		updateParams := repository.UpdateLinkMetadataParams{
			URL:          urlString,
			Title:        title,
			Description:  description,
			FaviconURL:   faviconURL,
			ImageURL:     imageURL,
			PlatformName: platformName,
			PlatformType: platformType,
			PlatformColor: platformColor,
			IsVerified:   nil,
		}
		
		updatedMetadata, err := s.repo.UpdateLinkMetadata(ctx, updateParams)
		if err != nil {
			s.logger.Errorf("Failed to update metadata: %v", err)
			return nil, errors.Wrap(err, "Failed to update metadata")
		}
		
		s.logger.Infof("Updated metadata for URL: %s", urlString)
		return mapLinkMetadataToDTO(updatedMetadata), nil
	}
	
	// Create new metadata
	createParams := repository.CreateLinkMetadataParams{
		Domain:       domain,
		URL:          urlString,
		Title:        title,
		Description:  description,
		FaviconURL:   faviconURL,
		ImageURL:     imageURL,
		PlatformName: platformName,
		PlatformType: platformType,
		PlatformColor: platformColor,
		IsVerified:   nil,
	}
	
	newMetadata, err := s.repo.CreateLinkMetadata(ctx, createParams)
	if err != nil {
		s.logger.Errorf("Failed to store metadata: %v", err)
		return nil, errors.Wrap(err, "Failed to store metadata")
	}
	
	s.logger.Infof("Created metadata for URL: %s", urlString)
	return mapLinkMetadataToDTO(newMetadata), nil
}

func (s *linkMetadataService) IsKnownPlatform(domain string) bool {
	_, found := PlatformRegistry[domain]
	return found
}

func (s *linkMetadataService) GetPlatformInfo(domain string) *PlatformInfo {
	if platform, found := PlatformRegistry[domain]; found {
		return &platform
	}
	return nil
}

func (s *linkMetadataService) ListKnownPlatforms(ctx context.Context) ([]*PlatformInfo, error) {
	platforms := make([]*PlatformInfo, 0, len(PlatformRegistry))
	
	for _, platform := range PlatformRegistry {
		platformCopy := platform // Create a copy to avoid reference issues
		platforms = append(platforms, &platformCopy)
	}
	
	return platforms, nil
}

// Helper functions
type HTMLMetadata struct {
	Title       string
	Description string
	FaviconURL  string
	ImageURL    string
}

func extractMetadata(n *html.Node, baseURL string) HTMLMetadata {
	var metadata HTMLMetadata
	
	// Extract data from head tags
	var extractFromHead func(*html.Node)
	extractFromHead = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.DataAtom {
			case atom.Title:
				metadata.Title = extractTextContent(n)
			case atom.Meta:
				// Check for description and OpenGraph tags
				var property, content string
				for _, attr := range n.Attr {
					switch attr.Key {
					case "name":
						property = attr.Val
					case "property":
						property = attr.Val
					case "content":
						content = attr.Val
					}
				}
				
				// Handle different meta tags
				switch strings.ToLower(property) {
				case "description":
					metadata.Description = content
				case "og:title":
					if metadata.Title == "" {
						metadata.Title = content
					}
				case "og:description":
					if metadata.Description == "" {
						metadata.Description = content
					}
				case "og:image":
					metadata.ImageURL = resolveURL(baseURL, content)
				case "twitter:image":
					if metadata.ImageURL == "" {
						metadata.ImageURL = resolveURL(baseURL, content)
					}
				}
			case atom.Link:
				// Check for favicon
				var rel, href string
				for _, attr := range n.Attr {
					switch attr.Key {
					case "rel":
						rel = attr.Val
					case "href":
						href = attr.Val
					}
				}
				
				if rel == "icon" || rel == "shortcut icon" {
					metadata.FaviconURL = resolveURL(baseURL, href)
				}
			}
		}
		
		// Process child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractFromHead(c)
		}
	}
	
	// Find the head tag and extract metadata
	var findHead func(*html.Node) bool
	findHead = func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.DataAtom == atom.Head {
			extractFromHead(n)
			return true
		}
		
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if findHead(c) {
				return true
			}
		}
		
		return false
	}
	
	findHead(n)
	return metadata
}

func extractTextContent(n *html.Node) string {
	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			text += c.Data
		} else if c.Type == html.ElementNode {
			text += extractTextContent(c)
		}
	}
	return strings.TrimSpace(text)
}

func resolveURL(base, ref string) string {
	baseURL, err := url.Parse(base)
	if err != nil {
		return ref
	}
	
	refURL, err := url.Parse(ref)
	if err != nil {
		return ref
	}
	
	resolvedURL := baseURL.ResolveReference(refURL)
	return resolvedURL.String()
}

func normalizeURL(urlString string) (string, error) {
	// Add scheme if missing
	if !strings.HasPrefix(urlString, "http://") && !strings.HasPrefix(urlString, "https://") {
		urlString = "https://" + urlString
	}
	
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}
	
	// Clean up URL by removing unnecessary parts
	parsedURL.RawQuery = ""
	parsedURL.Fragment = ""
	
	return parsedURL.String(), nil
}

func mapLinkMetadataToDTO(metadata *db.LinkMetadatum) *LinkMetadataDTO {
	dto := &LinkMetadataDTO{
		ID:           metadata.MetadataID.String(),
		URL:          metadata.Url,
		Domain:       metadata.Domain,
		IsVerified:   metadata.IsVerified != nil && *metadata.IsVerified,
	}
	
	if metadata.Title != nil {
		dto.Title = *metadata.Title
	}
	
	if metadata.Description != nil {
		dto.Description = *metadata.Description
	}
	
	if metadata.FaviconUrl != nil {
		dto.FaviconURL = *metadata.FaviconUrl
	}
	
	if metadata.ImageUrl != nil {
		dto.ImageURL = *metadata.ImageUrl
	}
	
	if metadata.PlatformName != nil {
		dto.PlatformName = *metadata.PlatformName
	}
	
	if metadata.PlatformType != nil {
		dto.PlatformType = *metadata.PlatformType
	}
	
	if metadata.PlatformColor != nil {
		dto.PlatformColor = *metadata.PlatformColor
	}
	
	if metadata.CreatedAt != nil {
		dto.CreatedAt = metadata.CreatedAt.Format(time.RFC3339)
	}
	
	if metadata.UpdatedAt != nil {
		dto.UpdatedAt = metadata.UpdatedAt.Format(time.RFC3339)
	}
	
	return dto
}