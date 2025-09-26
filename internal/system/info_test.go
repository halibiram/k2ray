package system_test

import (
	"k2ray/internal/system"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSystemInfo(t *testing.T) {
	info, err := system.GetSystemInfo()

	assert.NoError(t, err)
	assert.NotNil(t, info)

	// Check some mocked values
	assert.Equal(t, "keenetic-k2ray", info.Hostname)
	assert.Equal(t, "Keenetic Giga (KN-1011)", info.KeeneticModel)
	assert.Equal(t, "4.1.1", info.FirmwareVersion)

	// Check some runtime values
	assert.NotEmpty(t, info.OS)
	assert.NotZero(t, info.CPUCores)

	// Check some randomized values
	assert.GreaterOrEqual(t, info.CPUUsage, 0.0)
	assert.LessOrEqual(t, info.CPUUsage, 100.0)
	assert.NotEmpty(t, info.Uptime)
}
