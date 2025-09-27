package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

var (
	// HTTPRequestsTotal is a counter for total HTTP requests.
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "k2ray_http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	// HTTPRequestDuration is a histogram for request durations.
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "k2ray_http_request_duration_seconds",
			Help:    "Histogram of HTTP request durations.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// UserLoginsTotal is a counter for successful user logins.
	UserLoginsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "k2ray_user_logins_total",
			Help: "Total number of successful user logins.",
		},
		[]string{"status"}, // "success" or "failure"
	)

	// AppInfo is a gauge to expose application version and other info.
	AppInfo = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "k2ray_app_info",
			Help: "Information about the k2ray application.",
		},
		[]string{"version", "go_version"},
	)

	// SystemCPUUsage is a gauge for current system-wide CPU usage percentage.
	SystemCPUUsage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "k2ray_system_cpu_usage_percent",
			Help: "Current system-wide CPU usage percentage.",
		},
	)

	// SystemMemoryUsage is a gauge for memory usage.
	SystemMemoryUsage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "k2ray_system_memory_usage_bytes",
			Help: "System memory usage in bytes.",
		},
		[]string{"type"}, // "total", "used", "free"
	)

	// SystemDiskUsage is a gauge for disk usage.
	SystemDiskUsage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "k2ray_system_disk_usage_bytes",
			Help: "System disk usage in bytes.",
		},
		[]string{"path", "type"}, // "total", "used", "free"
	)
)

// InitMetrics initializes application-wide metrics and starts the system metrics collector.
func InitMetrics(appVersion, goVersion string) {
	// Set static application info
	AppInfo.With(prometheus.Labels{
		"version":    appVersion,
		"go_version": goVersion,
	}).Set(1)

	// Start the background collector
	go startSystemMetricsCollector()
}

// startSystemMetricsCollector runs a loop to collect system metrics periodically.
func startSystemMetricsCollector() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		collectCPUUsage()
		collectMemoryUsage()
		collectDiskUsage()
	}
}

// collectCPUUsage gets the current CPU usage and updates the gauge.
func collectCPUUsage() {
	percentages, err := cpu.Percent(0, false)
	if err != nil {
		log.Error().Err(err).Msg("Failed to collect CPU usage")
		return
	}
	if len(percentages) > 0 {
		SystemCPUUsage.Set(percentages[0])
	}
}

// collectMemoryUsage gets the current memory usage and updates the gauges.
func collectMemoryUsage() {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Error().Err(err).Msg("Failed to collect memory usage")
		return
	}
	SystemMemoryUsage.WithLabelValues("total").Set(float64(vmStat.Total))
	SystemMemoryUsage.WithLabelValues("used").Set(float64(vmStat.Used))
	SystemMemoryUsage.WithLabelValues("free").Set(float64(vmStat.Free))
}

// collectDiskUsage gets the disk usage for the root partition and updates the gauges.
func collectDiskUsage() {
	// Collecting for the root directory. This can be extended to other paths if needed.
	path := "/"
	usage, err := disk.Usage(path)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("Failed to collect disk usage")
		return
	}
	SystemDiskUsage.WithLabelValues(path, "total").Set(float64(usage.Total))
	SystemDiskUsage.WithLabelValues(path, "used").Set(float64(usage.Used))
	SystemDiskUsage.WithLabelValues(path, "free").Set(float64(usage.Free))
}