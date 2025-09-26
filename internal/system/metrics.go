package system

import (
	"math/rand"
	"time"
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
		Uplink:   rand.Int63n(100000),
		Downlink: rand.Int63n(1000000),
	}
}

// GetConnectionMetrics generates mock connection metrics.
func GetConnectionMetrics() *ConnectionMetrics {
	return &ConnectionMetrics{
		Active:   rand.Intn(100),
		Total:    rand.Intn(1000),
		Failures: rand.Intn(10),
	}
}

// GetPerformanceMetrics generates mock performance metrics.
func GetPerformanceMetrics() *PerformanceMetrics {
	return &PerformanceMetrics{
		CPUUsage:    rand.Float64() * 100,
		MemoryUsage: rand.Float64() * 100,
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}