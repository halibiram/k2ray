package system_test

import (
	"k2ray/internal/system"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSystemLogs(t *testing.T) {
	// Create a temporary dummy log file
	content := "line 1\nline 2"
	tmpfile, err := os.CreateTemp("", "testlog_*.log")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString(content)
	assert.NoError(t, err)
	tmpfile.Close()

	// Temporarily point the MockLogFilePath to our test file
	originalPath := system.MockLogFilePath
	system.MockLogFilePath = tmpfile.Name()
	defer func() { system.MockLogFilePath = originalPath }() // Restore original path

	// Test reading the logs
	logs, err := system.GetSystemLogs()
	assert.NoError(t, err)
	assert.Equal(t, content, logs)
}

func TestGetSystemLogs_FileNotExists(t *testing.T) {
	// Point to a non-existent file
	originalPath := system.MockLogFilePath
	system.MockLogFilePath = "non-existent-log-file.log"
	defer func() { system.MockLogFilePath = originalPath }()

	// Test that it returns an error
	_, err := system.GetSystemLogs()
	assert.Error(t, err)
}
