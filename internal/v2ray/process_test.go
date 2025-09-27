package v2ray_test

import (
	"k2ray/internal/api"
	"k2ray/internal/config"
	"k2ray/internal/db"
	"k2ray/internal/utils"
	"k2ray/internal/v2ray"
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	tmpfile, err := os.CreateTemp("", "test_v2ray_process_*.db")
	if err != nil {
		log.Fatalf("Failed to create temp db file: %v", err)
	}
	dbPath := tmpfile.Name()
	tmpfile.Close()

	config.LoadConfig("")
	config.AppConfig.DatabaseURL = dbPath

	db.InitDB()

	// Setup router for any potential handler calls if needed (good practice)
	testRouter := gin.Default()
	api.SetupRouter(testRouter, false) // Disable rate limiting for tests

	code := m.Run()

	db.DB.Close()
	os.Remove(dbPath)
	os.Exit(code)
}

func createTestUserAndConfig(t *testing.T) (userID, configID int64) {
	// Create User
	hashedPassword, _ := utils.HashPassword("password")
	res, err := db.DB.Exec(`INSERT INTO users (username, password_hash) VALUES (?, ?)`, "testuser", hashedPassword)
	assert.NoError(t, err)
	userID, _ = res.LastInsertId()

	// Create Config
	configData := `{"v": "2", "add": "test.com", "port": 443}`
	res, err = db.DB.Exec(`INSERT INTO configurations (user_id, name, protocol, config_data) VALUES (?, ?, ?, ?)`, userID, "test-config", "vmess", configData)
	assert.NoError(t, err)
	configID, _ = res.LastInsertId()
	return
}

func TestV2RayProcessManager(t *testing.T) {
	// Cleanup any previous test artifacts
	os.Remove(v2ray.V2RayConfigPath)

	// 1. Initial status should be stopped
	isRunning, _ := v2ray.Status()
	assert.False(t, isRunning, "Initial status should be stopped")

	// 2. Starting should fail if no active config is set
	err := v2ray.Start()
	assert.Error(t, err, "Start should fail without an active config")

	// 3. Set an active config
	_, configID := createTestUserAndConfig(t)
	_, err = db.DB.Exec(`INSERT INTO settings (key, value) VALUES (?, ?)`, v2ray.ActiveConfigKey, configID)
	assert.NoError(t, err)

	// 4. Start the service
	err = v2ray.Start()
	assert.NoError(t, err, "Start should succeed with an active config")

	// 5. Verify status and config file
	isRunning, pid := v2ray.Status()
	assert.True(t, isRunning, "Status should be running after start")
	assert.NotZero(t, pid, "PID should be non-zero after start")
	_, err = os.Stat(v2ray.V2RayConfigPath)
	assert.NoError(t, err, "Config file should be created")

	// 6. Stop the service
	err = v2ray.Stop()
	assert.NoError(t, err, "Stop should succeed")

	// 7. Verify final status
	isRunning, _ = v2ray.Status()
	assert.False(t, isRunning, "Status should be stopped after stop")

	// Cleanup
	os.Remove(v2ray.V2RayConfigPath)
}
