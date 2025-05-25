// pkg/metrics/system.go
package metrics

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

// SystemMetricsCollector collects runtime and system metrics
type SystemMetricsCollector struct {
	metrics *Metrics
}

func NewSystemMetricsCollector(metrics *Metrics) *SystemMetricsCollector {
	return &SystemMetricsCollector{
		metrics: metrics,
	}
}

// StartCollection starts collecting system metrics periodically
func (s *SystemMetricsCollector) StartCollection(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Collect initial metrics
	s.collectRuntimeMetrics()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.collectRuntimeMetrics()
		}
	}
}

func (s *SystemMetricsCollector) collectRuntimeMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Goroutines
	s.metrics.GoRoutinesActive.Set(float64(runtime.NumGoroutine()))

	// Memory metrics
	s.metrics.MemoryUsage.WithLabelValues("alloc").Set(float64(m.Alloc))
	s.metrics.MemoryUsage.WithLabelValues("total_alloc").Set(float64(m.TotalAlloc))
	s.metrics.MemoryUsage.WithLabelValues("sys").Set(float64(m.Sys))
	s.metrics.MemoryUsage.WithLabelValues("heap_alloc").Set(float64(m.HeapAlloc))
	s.metrics.MemoryUsage.WithLabelValues("heap_sys").Set(float64(m.HeapSys))
	s.metrics.MemoryUsage.WithLabelValues("heap_idle").Set(float64(m.HeapIdle))
	s.metrics.MemoryUsage.WithLabelValues("heap_inuse").Set(float64(m.HeapInuse))
	s.metrics.MemoryUsage.WithLabelValues("stack_inuse").Set(float64(m.StackInuse))
	s.metrics.MemoryUsage.WithLabelValues("stack_sys").Set(float64(m.StackSys))

	// GC metrics - record pause time if there was a recent GC
	if m.NumGC > 0 {
		// Get the most recent GC pause time
		lastPause := m.PauseNs[(m.NumGC+255)%256]
		s.metrics.GCDuration.Observe(float64(lastPause) / 1e9) // Convert to seconds
	}
}

// HealthChecker provides health check metrics
type HealthChecker struct {
	metrics *Metrics
	checks  map[string]HealthCheck
}

type HealthCheck func(ctx context.Context) error

func NewHealthChecker(metrics *Metrics) *HealthChecker {
	return &HealthChecker{
		metrics: metrics,
		checks:  make(map[string]HealthCheck),
	}
}

func (h *HealthChecker) RegisterCheck(name string, check HealthCheck) {
	h.checks[name] = check
}

func (h *HealthChecker) RunChecks(ctx context.Context) map[string]error {
	results := make(map[string]error)
	
	for name, check := range h.checks {
		start := time.Now()
		err := check(ctx)
		duration := time.Since(start)

		fmt.Println(duration)
		
		status := "healthy"
		if err != nil {
			status = "unhealthy"
			results[name] = err
		}
		
		// Record health check metrics
		h.metrics.ErrorsTotal.WithLabelValues("health_check", name, status).Inc()
	}
	
	return results
}

// StartHealthChecks runs health checks periodically
func (h *HealthChecker) StartHealthChecks(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			h.RunChecks(ctx)
		}
	}
}