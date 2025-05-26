// service/instrumented_analytics_service.go
package service

import (
	"context"

	"github.com/0xsj/gin-sqlc/log"
	"github.com/0xsj/gin-sqlc/pkg/metrics"
)

// InstrumentedAnalyticsService wraps AnalyticsService with metrics
type InstrumentedAnalyticsService struct {
	base    AnalyticsService
	metrics *metrics.Metrics
	logger  log.Logger
}

func NewInstrumentedAnalyticsService(base AnalyticsService, metrics *metrics.Metrics, logger log.Logger) AnalyticsService {
	return &InstrumentedAnalyticsService{
		base:    base,
		metrics: metrics,
		logger:  logger,
	}
}

func (s *InstrumentedAnalyticsService) RecordClick(ctx context.Context, input RecordClickInput) error {
	err := s.base.RecordClick(ctx, input)
	s.metrics.RecordAnalyticsEvent("click")
	
	if err != nil {
		s.metrics.RecordError("analytics_record_failure", "analytics_service", "error")
	}
	
	return err
}

func (s *InstrumentedAnalyticsService) RecordPageView(ctx context.Context, input RecordPageViewInput) error {
	err := s.base.RecordPageView(ctx, input)
	s.metrics.RecordAnalyticsEvent("page_view")
	
	if err != nil {
		s.metrics.RecordError("analytics_record_failure", "analytics_service", "error")
	}
	
	return err
}

func (s *InstrumentedAnalyticsService) GetContentItemAnalytics(ctx context.Context, itemID string, page, pageSize int) (*ContentItemAnalyticsDTO, error) {
	result, err := s.base.GetContentItemAnalytics(ctx, itemID, page, pageSize)
	
	if err != nil {
		s.metrics.RecordError("analytics_fetch_failure", "analytics_service", "warning")
	}
	
	return result, err
}

func (s *InstrumentedAnalyticsService) GetUserAnalytics(ctx context.Context, userID string, page, pageSize int) (*UserAnalyticsDTO, error) {
	result, err := s.base.GetUserAnalytics(ctx, userID, page, pageSize)
	
	if err != nil {
		s.metrics.RecordError("analytics_fetch_failure", "analytics_service", "warning")
	}
	
	return result, err
}

func (s *InstrumentedAnalyticsService) GetUserAnalyticsByTimeRange(ctx context.Context, userID string, input TimeRangeInput) (*TimeRangeAnalyticsDTO, error) {
	result, err := s.base.GetUserAnalyticsByTimeRange(ctx, userID, input)
	
	if err != nil {
		s.metrics.RecordError("analytics_fetch_failure", "analytics_service", "warning")
	}
	
	return result, err
}

func (s *InstrumentedAnalyticsService) GetItemAnalyticsByTimeRange(ctx context.Context, itemID string, input TimeRangeInput) (*ItemTimeRangeAnalyticsDTO, error) {
	result, err := s.base.GetItemAnalyticsByTimeRange(ctx, itemID, input)
	
	if err != nil {
		s.metrics.RecordError("analytics_fetch_failure", "analytics_service", "warning")
	}
	
	return result, err
}

func (s *InstrumentedAnalyticsService) GetProfilePageViewsByTimeRange(ctx context.Context, userID string, input TimeRangeInput) (*PageViewAnalyticsDTO, error) {
	result, err := s.base.GetProfilePageViewsByTimeRange(ctx, userID, input)
	
	if err != nil {
		s.metrics.RecordError("analytics_fetch_failure", "analytics_service", "warning")
	}
	
	return result, err
}

func (s *InstrumentedAnalyticsService) GetProfileDashboard(ctx context.Context, userID string, days int) (*ProfileDashboardDTO, error) {
	result, err := s.base.GetProfileDashboard(ctx, userID, days)
	
	if err != nil {
		s.metrics.RecordError("analytics_fetch_failure", "analytics_service", "warning")
	}
	
	return result, err
}

func (s *InstrumentedAnalyticsService) GetReferrerAnalytics(ctx context.Context, userID string, input TimeRangeInput) (*ReferrerAnalyticsDTO, error) {
	result, err := s.base.GetReferrerAnalytics(ctx, userID, input)
	
	if err != nil {
		s.metrics.RecordError("analytics_fetch_failure", "analytics_service", "warning")
	}
	
	return result, err
}

