package system

import (
	"os"
)

var MockLogFilePath = "configs/dummy_system.log"

// GetSystemLogs reads the content of the mocked system log file.
// In a real implementation, this would interact with journald or a log file.
func GetSystemLogs() (string, error) {
	content, err := os.ReadFile(MockLogFilePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
