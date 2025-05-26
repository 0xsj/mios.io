// pkg/metrics/database.go
package metrics

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// DatabaseMetricsCollector collects database pool metrics
type DatabaseMetricsCollector struct {
	pool    *pgxpool.Pool
	metrics *Metrics
}

func NewDatabaseMetricsCollector(pool *pgxpool.Pool, metrics *Metrics) *DatabaseMetricsCollector {
	return &DatabaseMetricsCollector{
		pool:    pool,
		metrics: metrics,
	}
}

// StartCollection starts collecting database metrics periodically
func (d *DatabaseMetricsCollector) StartCollection(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Collect initial metrics
	d.collectPoolMetrics()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			d.collectPoolMetrics()
		}
	}
}

func (d *DatabaseMetricsCollector) collectPoolMetrics() {
	stats := d.pool.Stat()
	
	// Connection pool metrics
	d.metrics.DBConnectionsActive.Set(float64(stats.AcquiredConns()))
	d.metrics.DBConnectionsIdle.Set(float64(stats.IdleConns()))
}

