package metirc

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

type AppMetricsExporter struct {
	registry            *prometheus.Registry     // Registry untuk semua metrik
	httpRequestsTotal   *prometheus.CounterVec   // HTTP metrics
	httpRequestDuration *prometheus.HistogramVec // HTTP metrics
	businessEvents      *prometheus.CounterVec   // Business metrics
	// System metrics
	memoryUsage     prometheus.Gauge
	goroutinesCount prometheus.Gauge
	uptime          prometheus.Counter
	buildInfo       *prometheus.GaugeVec // Instance metadata
	startTime       time.Time            // Waktu mulai aplikasi untuk perhitungan uptime
}

func NewAppMetricsExporter() *AppMetricsExporter {
	// Buat registry baru untuk metrik kustom aplikasi
	registry := prometheus.NewRegistry()

	// Daftarkan collector default untuk metrics Go runtime (GC, goroutines, dll)
	registry.MustRegister(collectors.NewGoCollector())
	// Daftarkan process collector untuk metrics tingkat OS (CPU, memory, file descriptors)
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	// Buat eksporter dengan semua metrik
	exporter := &AppMetricsExporter{
		registry: registry,

		// HTTP metrics
		httpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "app",
				Subsystem: "http",
				Name:      "requests_total",
				Help:      "Total count of HTTP requests by status, method, and endpoint",
			},
			[]string{"status", "method", "endpoint"},
		),
		httpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "app",
				Subsystem: "http",
				Name:      "request_duration_seconds",
				Help:      "Duration of HTTP requests in seconds",
				Buckets:   []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2.5, 5, 10},
			},
			[]string{"status", "method", "endpoint"},
		),

		// Business metrics
		businessEvents: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "app",
				Subsystem: "business",
				Name:      "events_total",
				Help:      "Total count of business events by type and user",
			},
			[]string{"event_type", "user_id"},
		),

		// System metrics
		memoryUsage: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "app",
				Subsystem: "system",
				Name:      "memory_bytes",
				Help:      "Current memory usage in bytes",
			},
		),
		goroutinesCount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "app",
				Subsystem: "system",
				Name:      "goroutines",
				Help:      "Current number of goroutines",
			},
		),
		uptime: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: "app",
				Subsystem: "system",
				Name:      "uptime_seconds",
				Help:      "The uptime of the application in seconds",
			},
		),

		// Build info
		buildInfo: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "app",
				Name:      "build_info",
				Help:      "Build information about the application",
			},
			[]string{"version", "go_version", "commit_hash"},
		),

		startTime: time.Now(),
	}

	// Register semua metrik ke registry
	registry.MustRegister(
		exporter.httpRequestsTotal,
		exporter.httpRequestDuration,
		exporter.businessEvents,
		exporter.memoryUsage,
		exporter.goroutinesCount,
		exporter.uptime,
		exporter.buildInfo,
	)

	// Set build info (sebagai contoh)
	exporter.buildInfo.WithLabelValues("1.0.0", runtime.Version(), "abc123").Set(1)

	// Mulai goroutine untuk memperbarui metrik sistem secara periodik
	go exporter.collectSystemMetrics()

	return exporter
}

// collectSystemMetrics mengumpulkan metrik sistem secara periodik
func (e *AppMetricsExporter) collectSystemMetrics() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Update metrik uptime
		e.uptime.Add(15) // 15 detik sejak tick terakhir

		// Update metrik lainnya...
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		e.memoryUsage.Set(float64(memStats.Alloc))
		e.goroutinesCount.Set(float64(runtime.NumGoroutine()))
	}
}

// MetricsHandler mengembalikan HTTP handler untuk endpoint /metrics
func (e *AppMetricsExporter) MetricsHandler() http.Handler {
	return promhttp.HandlerFor(e.registry, promhttp.HandlerOpts{})
}

// ObserveHTTPRequest mencatat metrik untuk request HTTP
func (e *AppMetricsExporter) ObserveHTTPRequest(status int, method, endpoint string, duration time.Duration) {
	statusStr := strconv.Itoa(status)
	e.httpRequestsTotal.WithLabelValues(statusStr, method, endpoint).Inc()
	e.httpRequestDuration.WithLabelValues(statusStr, method, endpoint).Observe(duration.Seconds())
}

// RecordBusinessEvent mencatat event bisnis
func (e *AppMetricsExporter) RecordBusinessEvent(eventType, userID string) {
	e.businessEvents.WithLabelValues(eventType, userID).Inc()
}

// FiberMetricMiddleware menyediakan middleware Gin untuk merekam metrik HTTP
func (e *AppMetricsExporter) FiberMetricMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Proses request
		err := c.Next()

		// Setelah request selesai
		duration := time.Since(start)
		status := c.Response().StatusCode()
		method := c.Method()
		path := c.Route().Path // atau c.Path() kalau tidak pakai route pattern

		if path == "" {
			path = "unknown"
		}

		e.ObserveHTTPRequest(status, method, path, duration)

		return err
	}
}
