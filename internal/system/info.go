package system

import (
	"k2ray/internal/utils"
	"runtime"
	"time"
)

// SystemInfo holds various pieces of system information.
type SystemInfo struct {
	Hostname        string  `json:"hostname"`
	OS              string  `json:"os"`
	Kernel          string  `json:"kernel"`
	CPU             string  `json:"cpu"`
	CPUCores        int     `json:"cpu_cores"`
	CPUUsage        float64 `json:"cpu_usage"`
	MemoryTotalMB   uint64  `json:"memory_total_mb"`
	MemoryUsedMB    uint64  `json:"memory_used_mb"`
	MemoryUsage     float64 `json:"memory_usage"`
	Uptime          string  `json:"uptime"`
	KeeneticModel   string  `json:"keenetic_model"`
	FirmwareVersion string  `json:"firmware_version"`
}

// GetSystemInfo gathers and returns system information.
// NOTE: This implementation uses mocked data for demonstration in a sandboxed environment.
func GetSystemInfo() (*SystemInfo, error) {
	// Mocked data
	usedMem := uint64(128) + utils.SecureUint64n(64)
	totalMem := uint64(256)

	info := &SystemInfo{
		Hostname:        "keenetic-k2ray",
		OS:              runtime.GOOS,
		Kernel:          "5.4.0-k2ray", // Mocked kernel
		CPU:             "MIPS 74Kc V5.0", // Mocked CPU for Keenetic
		CPUCores:        runtime.NumCPU(),
		CPUUsage:        utils.SecureFloat64() * 100,
		MemoryTotalMB:   totalMem,
		MemoryUsedMB:    usedMem,
		MemoryUsage:     (float64(usedMem) / float64(totalMem)) * 100,
		Uptime:          (time.Duration(utils.SecureIntn(3600*24*7)) * time.Second).String(), // Mocked uptime
		KeeneticModel:   "Keenetic Giga (KN-1011)",
		FirmwareVersion: "4.1.1",
	}

	return info, nil
}
