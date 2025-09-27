package system

import (
	"k2ray/internal/utils"
)

// TrafficMetrics represents traffic data.
type TrafficMetrics struct {
	Uplink   int64 `json:"uplink"`
	Downlink int64 `json:"downlink"`
}

// ConnectionMetrics represents connection data.
type ConnectionMetrics struct {
	Active   int `json:"active"`
	Total    int `json:"total"`
	Failures int `json:"failures"`
}

// PerformanceMetrics represents system performance data.
type PerformanceMetrics struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
}

// GetTrafficMetrics generates mock traffic metrics.
func GetTrafficMetrics() *TrafficMetrics {
	return &TrafficMetrics{
		Uplink:   utils.SecureIntn(100000),
		Downlink: utils.SecureIntn(1000000),
	}
}

// GetConnectionMetrics generates mock connection metrics.
func GetConnectionMetrics() *ConnectionMetrics {
	return &ConnectionMetrics{
		Active:   int(utils.SecureIntn(100)),
		Total:    int(utils.SecureIntn(1000)),
		Failures: int(utils.SecureIntn(10)),
	}
}

// GetPerformanceMetrics generates mock performance metrics.
func GetPerformanceMetrics() *PerformanceMetrics {
	return &PerformanceMetrics{
		CPUUsage:    utils.SecureFloat64() * 100,
		MemoryUsage: utils.SecureFloat64() * 100,
	}
}