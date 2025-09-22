package v2ray

import (
	"database/sql"
	"errors"
	"fmt"
	"k2ray/internal/db"
	"log"
	"os"
	"sync"
)

const (
	ActiveConfigKey   = "active_config_id"
	V2RayConfigPath   = "/tmp/k2ray_config.json"
	V2RayExecutable   = "/usr/bin/v2ray" // Assumed path
)

// ManagerState holds the current state of the V2Ray process manager.
// NOTE: This is a mocked implementation for development purposes, as the v2ray-core
// executable is not available in the sandboxed environment. The logic here simulates
// the behavior of a real process manager.
type ManagerState struct {
	mu        sync.RWMutex
	isRunning bool
	pid       int // In a real implementation, this would store the process ID.
}

// manager is a singleton instance of the ManagerState.
var manager = &ManagerState{}

// Start fetches the active config, writes it to a file, and mocks starting the V2Ray process.
func Start() error {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	if manager.isRunning {
		return errors.New("V2Ray process is already running")
	}

	// 1. Get active config ID from DB
	var configID int64
	err := db.DB.QueryRow("SELECT value FROM system_settings WHERE key = ?", ActiveConfigKey).Scan(&configID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("no active V2Ray configuration is set")
		}
		return fmt.Errorf("could not get active config: %w", err)
	}

	// 2. Get config data from DB
	var configData string
	err = db.DB.QueryRow("SELECT config_data FROM v2ray_configs WHERE id = ?", configID).Scan(&configData)
	if err != nil {
		return fmt.Errorf("could not retrieve config data for ID %d: %w", configID, err)
	}

	// 3. Write config to file
	err = os.WriteFile(V2RayConfigPath, []byte(configData), 0644)
	if err != nil {
		return fmt.Errorf("could not write V2Ray config file: %w", err)
	}

	// 4. Mock starting the process
	log.Printf("MOCK: Would run command: %s -config %s", V2RayExecutable, V2RayConfigPath)
	manager.isRunning = true
	manager.pid = 12345 // Mock PID

	log.Println("Mock V2Ray process started successfully.")
	return nil
}

// Stop mocks stopping the V2Ray process.
func Stop() error {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	if !manager.isRunning {
		return errors.New("V2Ray process is not running")
	}

	// Mock killing the process
	log.Printf("MOCK: Would kill process with PID: %d", manager.pid)
	manager.isRunning = false
	manager.pid = 0

	log.Println("Mock V2Ray process stopped successfully.")
	return nil
}

// Status returns the current (mocked) status of the V2Ray process.
func Status() (isRunning bool, pid int) {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	return manager.isRunning, manager.pid
}
