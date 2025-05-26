// service/cached_link_metadata_service.go
package service

import (
	"context"
	"time"

	"github.com/0xsj/gin-sqlc/log"
	"github.com/0xsj/gin-sqlc/pkg/cache"
)

// CachedLinkMetadataService wraps the regular link metadata service with caching
type CachedLinkMetadataService struct {
	baseService LinkMetadataService
	cache       cache.CacheService
	keyBuilder  *cache.CacheKeyBuilder
	logger      log.Logger
}

func NewCachedLinkMetadataService(
	baseService LinkMetadataService,
	cacheService cache.CacheService,
	logger log.Logger,
) LinkMetadataService {
	return &CachedLinkMetadataService{
		baseService: baseService,
		cache:       cacheService,
		keyBuilder:  cache.NewCacheKeyBuilder(),
		logger:      logger,
	}
}

func (s *CachedLinkMetadataService) GetMetadata(ctx context.Context, urlString string) (*LinkMetadataDTO, error) {
	cacheKey := s.keyBuilder.LinkMetadata(urlString)
	
	var result LinkMetadataDTO
	err := s.cache.GetOrSet(ctx, cacheKey, &result, cache.GetMetadataTTL(), func() (interface{}, error) {
		s.logger.Debugf("Cache miss for link metadata, fetching from service")
		return s.baseService.GetMetadata(ctx, urlString)
	})
	
	if err != nil {
		s.logger.Errorf("Failed to get cached link metadata: %v", err)
		// Fallback to direct service call
		return s.baseService.GetMetadata(ctx, urlString)
	}
	
	return &result, nil
}

func (s *CachedLinkMetadataService) FetchAndStoreMetadata(ctx context.Context, urlString string) (*LinkMetadataDTO, error) {
	// Always fetch fresh data and update cache
	result, err := s.baseService.FetchAndStoreMetadata(ctx, urlString)
	if err != nil {
		return nil, err
	}
	
	// Update cache with fresh data
	cacheKey := s.keyBuilder.LinkMetadata(urlString)
	if cacheErr := s.cache.Set(ctx, cacheKey, result, cache.GetMetadataTTL()); cacheErr != nil {
		s.logger.Warnf("Failed to cache fresh metadata: %v", cacheErr)
	}
	
	return result, nil
}

func (s *CachedLinkMetadataService) IsKnownPlatform(domain string) bool {
	return s.baseService.IsKnownPlatform(domain)
}

func (s *CachedLinkMetadataService) GetPlatformInfo(domain string) *PlatformInfo {
	return s.baseService.GetPlatformInfo(domain)
}

func (s *CachedLinkMetadataService) ListKnownPlatforms(ctx context.Context) ([]*PlatformInfo, error) {
	// Platform list rarely changes, cache for longer
	cacheKey := "platforms:known"
	
	var result []*PlatformInfo
	err := s.cache.GetOrSet(ctx, cacheKey, &result, 24*time.Hour, func() (interface{}, error) {
		s.logger.Debugf("Cache miss for known platforms, fetching from service")
		return s.baseService.ListKnownPlatforms(ctx)
	})
	
	if err != nil {
		s.logger.Errorf("Failed to get cached platforms: %v", err)
		// Fallback to direct service call
		return s.baseService.ListKnownPlatforms(ctx)
	}
	
	return result, nil
}