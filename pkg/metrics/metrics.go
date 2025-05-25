package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "gin_sqlc"
	subsystem = "api"
)

type Metrics struct {
	// HTTP metrics
	HTTPRequestsTotal     *prometheus.CounterVec
	HTTPRequestDuration   *prometheus.HistogramVec
	HTTPRequestSize       *prometheus.HistogramVec
	HTTPResponseSize      *prometheus.HistogramVec
	HTTPActiveConnections prometheus.Gauge

	// Database metrics
	DBConnectionsActive   prometheus.Gauge
	DBConnectionsIdle     prometheus.Gauge
	DBQueriesTotal        *prometheus.CounterVec
	DBQueryDuration       *prometheus.HistogramVec
	DBTransactionsTotal   *prometheus.CounterVec

	// Cache metrics
	CacheOperationsTotal *prometheus.CounterVec
	CacheHitRatio        *prometheus.GaugeVec
	CacheDuration        *prometheus.HistogramVec

	// Rate limiting metrics
	RateLimitHitsTotal     *prometheus.CounterVec
	RateLimitAllowed       *prometheus.CounterVec
	RateLimitRemaining     *prometheus.GaugeVec

	// Business metrics
	UsersTotal            prometheus.Gauge
	ContentItemsTotal     prometheus.Gauge
	AnalyticsEventsTotal  *prometheus.CounterVec
	EmailsSentTotal       *prometheus.CounterVec

	// Error metrics
	ErrorsTotal           *prometheus.CounterVec
	PanicTotal            prometheus.Counter

	// Performance metrics
	GoRoutinesActive      prometheus.Gauge
	MemoryUsage          *prometheus.GaugeVec
	GCDuration           prometheus.Histogram
}