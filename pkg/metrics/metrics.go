// pkg/metrics/metrics.go
package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "gin_sqlc"
	subsystem = "api"
)

// Metrics holds all the Prometheus metrics
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

// NewMetrics creates and registers all metrics
func NewMetrics() *Metrics {
	m := &Metrics{
		// HTTP metrics
		HTTPRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "http_requests_total",
				Help:      "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status_code", "status_class"},
		),

		HTTPRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "http_request_duration_seconds",
				Help:      "HTTP request duration in seconds",
				Buckets:   []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
			},
			[]string{"method", "endpoint", "status_class"},
		),

		HTTPRequestSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "http_request_size_bytes",
				Help:      "HTTP request size in bytes",
				Buckets:   prometheus.ExponentialBuckets(100, 10, 8),
			},
			[]string{"method", "endpoint"},
		),

		HTTPResponseSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "http_response_size_bytes",
				Help:      "HTTP response size in bytes",
				Buckets:   prometheus.ExponentialBuckets(100, 10, 8),
			},
			[]string{"method", "endpoint", "status_class"},
		),

		HTTPActiveConnections: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "http_active_connections",
				Help:      "Number of active HTTP connections",
			},
		),

		// Database metrics
		DBConnectionsActive: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "database",
				Name:      "connections_active",
				Help:      "Number of active database connections",
			},
		),

		DBConnectionsIdle: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "database",
				Name:      "connections_idle",
				Help:      "Number of idle database connections",
			},
		),

		DBQueriesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "database",
				Name:      "queries_total",
				Help:      "Total number of database queries",
			},
			[]string{"operation", "table", "status"},
		),

		DBQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "database",
				Name:      "query_duration_seconds",
				Help:      "Database query duration in seconds",
				Buckets:   []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2, 5},
			},
			[]string{"operation", "table"},
		),

		DBTransactionsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "database",
				Name:      "transactions_total",
				Help:      "Total number of database transactions",
			},
			[]string{"status"},
		),

		// Cache metrics
		CacheOperationsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "cache",
				Name:      "operations_total",
				Help:      "Total number of cache operations",
			},
			[]string{"operation", "result"},
		),

		CacheHitRatio: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "cache",
				Name:      "hit_ratio",
				Help:      "Cache hit ratio",
			},
			[]string{"cache_type"},
		),

		CacheDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "cache",
				Name:      "operation_duration_seconds",
				Help:      "Cache operation duration in seconds",
				Buckets:   []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1},
			},
			[]string{"operation"},
		),

		// Rate limiting metrics
		RateLimitHitsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "rate_limit",
				Name:      "hits_total",
				Help:      "Total number of rate limit hits",
			},
			[]string{"endpoint", "result"},
		),

		RateLimitAllowed: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "rate_limit",
				Name:      "allowed_total",
				Help:      "Total number of allowed requests",
			},
			[]string{"endpoint", "key_type"},
		),

		RateLimitRemaining: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "rate_limit",
				Name:      "remaining",
				Help:      "Remaining requests in rate limit window",
			},
			[]string{"endpoint", "key"},
		),

		// Business metrics
		UsersTotal: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "business",
				Name:      "users_total",
				Help:      "Total number of users",
			},
		),

		ContentItemsTotal: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "business",
				Name:      "content_items_total",
				Help:      "Total number of content items",
			},
		),

		AnalyticsEventsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "business",
				Name:      "analytics_events_total",
				Help:      "Total number of analytics events",
			},
			[]string{"event_type"},
		),

		EmailsSentTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "business",
				Name:      "emails_sent_total",
				Help:      "Total number of emails sent",
			},
			[]string{"template", "status"},
		),

		// Error metrics
		ErrorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "errors_total",
				Help:      "Total number of errors",
			},
			[]string{"type", "component", "severity"},
		),

		PanicTotal: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "panics_total",
				Help:      "Total number of panics",
			},
		),

		// Performance metrics
		GoRoutinesActive: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "runtime",
				Name:      "goroutines_active",
				Help:      "Number of active goroutines",
			},
		),

		MemoryUsage: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "runtime",
				Name:      "memory_usage_bytes",
				Help:      "Memory usage in bytes",
			},
			[]string{"type"},
		),

		GCDuration: promauto.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "runtime",
				Name:      "gc_duration_seconds",
				Help:      "Garbage collection duration in seconds",
				Buckets:   []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1},
			},
		),
	}

	return m
}

// Helper methods for common operations
func (m *Metrics) RecordHTTPRequest(method, endpoint string, statusCode int, duration time.Duration, requestSize, responseSize int64) {
	statusClass := getStatusClass(statusCode)
	statusCodeStr := strconv.Itoa(statusCode)

	m.HTTPRequestsTotal.WithLabelValues(method, endpoint, statusCodeStr, statusClass).Inc()
	m.HTTPRequestDuration.WithLabelValues(method, endpoint, statusClass).Observe(duration.Seconds())
	m.HTTPRequestSize.WithLabelValues(method, endpoint).Observe(float64(requestSize))
	m.HTTPResponseSize.WithLabelValues(method, endpoint, statusClass).Observe(float64(responseSize))
}

func (m *Metrics) RecordDBQuery(operation, table string, duration time.Duration, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}

	m.DBQueriesTotal.WithLabelValues(operation, table, status).Inc()
	m.DBQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}

func (m *Metrics) RecordCacheOperation(operation, result string, duration time.Duration) {
	m.CacheOperationsTotal.WithLabelValues(operation, result).Inc()
	m.CacheDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

func (m *Metrics) RecordRateLimit(endpoint string, allowed bool, remaining int, keyType string) {
	result := "blocked"
	if allowed {
		result = "allowed"
		m.RateLimitAllowed.WithLabelValues(endpoint, keyType).Inc()
	}

	m.RateLimitHitsTotal.WithLabelValues(endpoint, result).Inc()
}

func (m *Metrics) RecordError(errorType, component, severity string) {
	m.ErrorsTotal.WithLabelValues(errorType, component, severity).Inc()
}

func (m *Metrics) RecordAnalyticsEvent(eventType string) {
	m.AnalyticsEventsTotal.WithLabelValues(eventType).Inc()
}

func (m *Metrics) RecordEmailSent(template, status string) {
	m.EmailsSentTotal.WithLabelValues(template, status).Inc()
}

// Helper functions
func getStatusClass(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return "2xx"
	case statusCode >= 300 && statusCode < 400:
		return "3xx"
	case statusCode >= 400 && statusCode < 500:
		return "4xx"
	case statusCode >= 500:
		return "5xx"
	default:
		return "unknown"
	}
}

// normalizeEndpoint normalizes endpoints for metrics (remove IDs, etc.)
func NormalizeEndpoint(path string) string {
	// This is a simple implementation - you might want to use a more sophisticated
	// approach like replacing UUIDs and IDs with placeholders
	// For now, return the path as-is, but you could implement pattern matching
	return path
}