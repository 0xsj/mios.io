// service/cached_analytics_service.go
package service

import (
	"context"
	"fmt"

	"github.com/0xsj/mios.io/log"
	"github.com/0xsj/mios.io/pkg/cache"
)

// CachedAnalyticsService wraps the regular analytics service with caching
type CachedAnalyticsService struct {
	baseService AnalyticsService
	cache       cache.CacheService
	keyBuilder  *cache.CacheKeyBuilder
	logger      log.Logger
}

func NewCachedAnalyticsService(
	baseService AnalyticsService,
	cacheService cache.CacheService,
	logger log.Logger,
) AnalyticsService {
	return &CachedAnalyticsService{
		baseService: baseService,
		cache:       cacheService,
		keyBuilder:  cache.NewCacheKeyBuilder(),
		logger:      logger,
	}
}

func (s *CachedAnalyticsService) RecordClick(ctx context.Context, input RecordClickInput) error {
	// Recording operations should invalidate related cache
	err := s.baseService.RecordClick(ctx, input)
	if err != nil {
		return err
	}

	// Invalidate related caches asynchronously
	go s.invalidateUserAnalyticsCache(context.Background(), input.UserID)
	
	return nil
}

func (s *CachedAnalyticsService) RecordPageView(ctx context.Context, input RecordPageViewInput) error {
	// Recording operations should invalidate related cache
	err := s.baseService.RecordPageView(ctx, input)
	if err != nil {
		return err
	}

	// Invalidate related caches asynchronously
	go s.invalidateUserAnalyticsCache(context.Background(), input.ProfileID)
	
	return nil
}

func (s *CachedAnalyticsService) GetContentItemAnalytics(ctx context.Context, itemID string, page, pageSize int) (*ContentItemAnalyticsDTO, error) {
	// Simple operations with pagination are not cached due to complexity
	return s.baseService.GetContentItemAnalytics(ctx, itemID, page, pageSize)
}

func (s *CachedAnalyticsService) GetUserAnalytics(ctx context.Context, userID string, page, pageSize int) (*UserAnalyticsDTO, error) {
	// Simple operations with pagination are not cached due to complexity
	return s.baseService.GetUserAnalytics(ctx, userID, page, pageSize)
}

func (s *CachedAnalyticsService) GetUserAnalyticsByTimeRange(ctx context.Context, userID string, input TimeRangeInput) (*TimeRangeAnalyticsDTO, error) {
	cacheKey := s.keyBuilder.TimeRangeAnalytics(userID, input.StartDate, input.EndDate)
	
	var result TimeRangeAnalyticsDTO
	err := s.cache.GetOrSet(ctx, cacheKey, &result, cache.GetAnalyticsTTL(), func() (interface{}, error) {
		s.logger.Debugf("Cache miss for time range analytics, fetching from database")
		return s.baseService.GetUserAnalyticsByTimeRange(ctx, userID, input)
	})
	
	if err != nil {
		s.logger.Errorf("Failed to get cached time range analytics: %v", err)
		// Fallback to direct service call
		return s.baseService.GetUserAnalyticsByTimeRange(ctx, userID, input)
	}
	
	return &result, nil
}

func (s *CachedAnalyticsService) GetItemAnalyticsByTimeRange(ctx context.Context, itemID string, input TimeRangeInput) (*ItemTimeRangeAnalyticsDTO, error) {
	timeRange := fmt.Sprintf("%s:%s:%d", input.StartDate, input.EndDate, input.Limit)
	cacheKey := s.keyBuilder.ContentItemAnalytics(itemID, timeRange)
	
	var result ItemTimeRangeAnalyticsDTO
	err := s.cache.GetOrSet(ctx, cacheKey, &result, cache.GetAnalyticsTTL(), func() (interface{}, error) {
		s.logger.Debugf("Cache miss for item analytics by time range, fetching from database")
		return s.baseService.GetItemAnalyticsByTimeRange(ctx, itemID, input)
	})
	
	if err != nil {
		s.logger.Errorf("Failed to get cached item analytics: %v", err)
		// Fallback to direct service call
		return s.baseService.GetItemAnalyticsByTimeRange(ctx, itemID, input)
	}
	
	return &result, nil
}

func (s *CachedAnalyticsService) GetProfilePageViewsByTimeRange(ctx context.Context, userID string, input TimeRangeInput) (*PageViewAnalyticsDTO, error) {
	cacheKey := s.keyBuilder.PageViewAnalytics(userID, input.StartDate, input.EndDate, input.Limit)
	
	var result PageViewAnalyticsDTO
	err := s.cache.GetOrSet(ctx, cacheKey, &result, cache.GetAnalyticsTTL(), func() (interface{}, error) {
		s.logger.Debugf("Cache miss for page views by time range, fetching from database")
		return s.baseService.GetProfilePageViewsByTimeRange(ctx, userID, input)
	})
	
	if err != nil {
		s.logger.Errorf("Failed to get cached page view analytics: %v", err)
		// Fallback to direct service call
		return s.baseService.GetProfilePageViewsByTimeRange(ctx, userID, input)
	}
	
	return &result, nil
}

func (s *CachedAnalyticsService) GetProfileDashboard(ctx context.Context, userID string, days int) (*ProfileDashboardDTO, error) {
	cacheKey := s.keyBuilder.ProfileDashboard(userID, days)
	
	var result ProfileDashboardDTO
	err := s.cache.GetOrSet(ctx, cacheKey, &result, cache.GetDashboardTTL(), func() (interface{}, error) {
		s.logger.Debugf("Cache miss for profile dashboard, fetching from database")
		return s.baseService.GetProfileDashboard(ctx, userID, days)
	})
	
	if err != nil {
		s.logger.Errorf("Failed to get cached profile dashboard: %v", err)
		// Fallback to direct service call
		return s.baseService.GetProfileDashboard(ctx, userID, days)
	}
	
	return &result, nil
}

func (s *CachedAnalyticsService) GetReferrerAnalytics(ctx context.Context, userID string, input TimeRangeInput) (*ReferrerAnalyticsDTO, error) {
	cacheKey := s.keyBuilder.ReferrerAnalytics(userID, input.StartDate, input.EndDate, input.Limit)
	
	var result ReferrerAnalyticsDTO
	err := s.cache.GetOrSet(ctx, cacheKey, &result, cache.GetAnalyticsTTL(), func() (interface{}, error) {
		s.logger.Debugf("Cache miss for referrer analytics, fetching from database")
		return s.baseService.GetReferrerAnalytics(ctx, userID, input)
	})
	
	if err != nil {
		s.logger.Errorf("Failed to get cached referrer analytics: %v", err)
		// Fallback to direct service call
		return s.baseService.GetReferrerAnalytics(ctx, userID, input)
	}
	
	return &result, nil
}

func (s *CachedAnalyticsService) invalidateUserAnalyticsCache(ctx context.Context, userID string) {
	// Invalidate all user-related analytics caches
	patterns := []string{
		s.keyBuilder.UserAnalyticsPattern(userID),
		s.keyBuilder.UserDashboardPattern(userID),
	}
	
	for _, pattern := range patterns {
		if err := s.cache.DeletePattern(ctx, pattern); err != nil {
			s.logger.Warnf("Failed to invalidate cache pattern %s: %v", pattern, err)
		} else {
			s.logger.Debugf("Invalidated cache pattern: %s", pattern)
		}
	}
}