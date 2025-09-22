package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"k2ray/internal/api"
	"k2ray/internal/config"
	"k2ray/internal/db"
	"k2ray/internal/system"
	"k2ray/internal/utils"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var testRouter *gin.Engine

type TestUser struct {
	ID       int64
	Username string
	Password string
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	tmpfile, err := os.CreateTemp("", "test_handlers_*.db")
	if err != nil {
		log.Fatalf("Failed to create temp db file: %v", err)
	}
	dbPath := tmpfile.Name()
	tmpfile.Close()

	config.LoadConfig("")
	config.AppConfig.DatabaseURL = dbPath
	config.AppConfig.JWTSecret = "a-very-secure-test-secret"

	db.InitDB()
	db.RunMigrations()

	createTestUser("user1", "password123")
	createTestUser("user2", "password456")

	// Create dummy log file for log tests
	// The test binary runs from within the package dir, so we need to create the parent dir first.
	err = os.MkdirAll("configs", 0755)
	if err != nil {
		log.Fatalf("Failed to create configs dir for dummy log: %v", err)
	}
	err = os.WriteFile(system.MockLogFilePath, []byte("[2025-09-22 14:20:01] K2Ray[123]: Service starting..."), 0644)
	if err != nil {
		log.Fatalf("Failed to create dummy log file: %v", err)
	}

	testRouter = gin.Default()
	api.SetupRouter(testRouter)

	code := m.Run()

	db.DB.Close()
	os.Remove(dbPath)
	os.Exit(code)
}

func createTestUser(username, password string) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Fatalf("Failed to hash test user password: %v", err)
	}
	insertSQL := `INSERT INTO users (username, password_hash) VALUES (?, ?)`
	_, err = db.DB.Exec(insertSQL, username, hashedPassword)
	if err != nil {
		log.Fatalf("Failed to create test user '%s': %v", username, err)
	}
}

func loginAs(t *testing.T, username, password string) (accessToken, refreshToken string) {
	loginPayload := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(loginPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if !assert.Equal(t, http.StatusOK, w.Code, "Login helper failed for user "+username) {
		t.FailNow()
	}
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	accessToken = response["access_token"]
	refreshToken = response["refresh_token"]
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
	return
}

func TestAuthEndpoints(t *testing.T) {
	t.Run("Login and Middleware", func(t *testing.T) {
		accessToken, _ := loginAs(t, "user1", "password123")
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "user1", response["username"])
	})

	t.Run("Logout and Revocation", func(t *testing.T) {
		accessToken, _ := loginAs(t, "user1", "password123")
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		req2, _ := http.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
		req2.Header.Set("Authorization", "Bearer "+accessToken)
		w2 := httptest.NewRecorder()
		testRouter.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusUnauthorized, w2.Code)
	})

	t.Run("Refresh and Rotation", func(t *testing.T) {
		_, refreshToken := loginAs(t, "user1", "password123")
		refreshPayload := fmt.Sprintf(`{"refresh_token": "%s"}`, refreshToken)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewBufferString(refreshPayload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		req2, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewBufferString(refreshPayload))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		testRouter.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusUnauthorized, w2.Code)
	})
}

func TestV2rayConfigCRUD(t *testing.T) {
	accessToken, _ := loginAs(t, "user1", "password123")
	var createdConfig db.V2rayConfig

	// 1. Create
	configPayload := `{"name": "My Server", "protocol": "vmess", "config_data": {"v": "2", "add": "test.com", "port": 443}}`
	createReq, _ := http.NewRequest(http.MethodPost, "/api/v1/configs", bytes.NewBufferString(configPayload))
	createReq.Header.Set("Authorization", "Bearer "+accessToken)
	createW := httptest.NewRecorder()
	testRouter.ServeHTTP(createW, createReq)
	assert.Equal(t, http.StatusCreated, createW.Code)
	json.Unmarshal(createW.Body.Bytes(), &createdConfig)
	assert.Equal(t, "My Server", createdConfig.Name)

	// 2. Get
	getReq, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/configs/%d", createdConfig.ID), nil)
	getReq.Header.Set("Authorization", "Bearer "+accessToken)
	getW := httptest.NewRecorder()
	testRouter.ServeHTTP(getW, getReq)
	assert.Equal(t, http.StatusOK, getW.Code)

	// 3. List
	listReq, _ := http.NewRequest(http.MethodGet, "/api/v1/configs", nil)
	listReq.Header.Set("Authorization", "Bearer "+accessToken)
	listW := httptest.NewRecorder()
	testRouter.ServeHTTP(listW, listReq)
	assert.Equal(t, http.StatusOK, listW.Code)
	var configs []db.V2rayConfig
	json.Unmarshal(listW.Body.Bytes(), &configs)
	assert.NotEmpty(t, configs)

	// 4. Update
	updatePayload := `{"name": "My Updated Server"}`
	updateReq, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/configs/%d", createdConfig.ID), bytes.NewBufferString(updatePayload))
	updateReq.Header.Set("Authorization", "Bearer "+accessToken)
	updateW := httptest.NewRecorder()
	testRouter.ServeHTTP(updateW, updateReq)
	assert.Equal(t, http.StatusOK, updateW.Code)
	var updatedConfig db.V2rayConfig
	json.Unmarshal(updateW.Body.Bytes(), &updatedConfig)
	assert.Equal(t, "My Updated Server", updatedConfig.Name)

	// 5. Delete
	delReq, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/configs/%d", createdConfig.ID), nil)
	delReq.Header.Set("Authorization", "Bearer "+accessToken)
	delW := httptest.NewRecorder()
	testRouter.ServeHTTP(delW, delReq)
	assert.Equal(t, http.StatusNoContent, delW.Code)

	// 6. Verify Deletion
	getReq2, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/configs/%d", createdConfig.ID), nil)
	getReq2.Header.Set("Authorization", "Bearer "+accessToken)
	getW2 := httptest.NewRecorder()
	testRouter.ServeHTTP(getW2, getReq2)
	assert.Equal(t, http.StatusNotFound, getW2.Code)
}

func TestV2rayAccessControl(t *testing.T) {
	// User 1 creates a config
	user1Token, _ := loginAs(t, "user1", "password123")
	configPayload := `{"name": "User 1s Secret", "protocol": "vmess", "config_data": {"v": "2", "add": "user1.com", "port": 443}}`
	createReq, _ := http.NewRequest(http.MethodPost, "/api/v1/configs", bytes.NewBufferString(configPayload))
	createReq.Header.Set("Authorization", "Bearer "+user1Token)
	createW := httptest.NewRecorder()
	testRouter.ServeHTTP(createW, createReq)
	assert.Equal(t, http.StatusCreated, createW.Code)
	var user1Config db.V2rayConfig
	json.Unmarshal(createW.Body.Bytes(), &user1Config)

	// User 2 logs in
	user2Token, _ := loginAs(t, "user2", "password456")

	// User 2 tries to GET User 1's config
	getReq, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/configs/%d", user1Config.ID), nil)
	getReq.Header.Set("Authorization", "Bearer "+user2Token)
	getW := httptest.NewRecorder()
	testRouter.ServeHTTP(getW, getReq)
	assert.Equal(t, http.StatusNotFound, getW.Code)

	// User 2 tries to UPDATE User 1's config
	updatePayload := `{"name": "Hacked"}`
	updateReq, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/configs/%d", user1Config.ID), bytes.NewBufferString(updatePayload))
	updateReq.Header.Set("Authorization", "Bearer "+user2Token)
	updateW := httptest.NewRecorder()
	testRouter.ServeHTTP(updateW, updateReq)
	assert.Equal(t, http.StatusNotFound, updateW.Code)

	// User 2 tries to DELETE User 1's config
	delReq, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/configs/%d", user1Config.ID), nil)
	delReq.Header.Set("Authorization", "Bearer "+user2Token)
	delW := httptest.NewRecorder()
	testRouter.ServeHTTP(delW, delReq)
	assert.Equal(t, http.StatusNotFound, delW.Code)
}

func TestSystemEndpoints(t *testing.T) {
	accessToken, _ := loginAs(t, "user1", "password123")

	t.Run("Get System Info", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/system/info", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var info system.SystemInfo
		err := json.Unmarshal(w.Body.Bytes(), &info)
		assert.NoError(t, err)
		assert.Equal(t, "keenetic-k2ray", info.Hostname)
	})

	t.Run("Get System Logs", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/system/logs", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "[2025-09-22 14:20:01] K2Ray[123]: Service starting...")
	})
}

func TestV2RayProcessEndpoints(t *testing.T) {
	accessToken, _ := loginAs(t, "user1", "password123")

	// 1. Check initial status
	statusReq, _ := http.NewRequest(http.MethodGet, "/api/v1/v2ray/status", nil)
	statusReq.Header.Set("Authorization", "Bearer "+accessToken)
	statusW := httptest.NewRecorder()
	testRouter.ServeHTTP(statusW, statusReq)
	assert.Equal(t, http.StatusOK, statusW.Code)
	var statusResponse map[string]any
	json.Unmarshal(statusW.Body.Bytes(), &statusResponse)
	assert.Equal(t, "stopped", statusResponse["status"])

	// 2. Create a config to use
	configPayload := `{"name": "My Active Server", "protocol": "vmess", "config_data": {"v": "2", "add": "active.com", "port": 443}}`
	createReq, _ := http.NewRequest(http.MethodPost, "/api/v1/configs", bytes.NewBufferString(configPayload))
	createReq.Header.Set("Authorization", "Bearer "+accessToken)
	createW := httptest.NewRecorder()
	testRouter.ServeHTTP(createW, createReq)
	assert.Equal(t, http.StatusCreated, createW.Code)
	var createdConfig db.V2rayConfig
	json.Unmarshal(createW.Body.Bytes(), &createdConfig)

	// 3. Set it as active
	activePayload := fmt.Sprintf(`{"config_id": %d}`, createdConfig.ID)
	activeReq, _ := http.NewRequest(http.MethodPost, "/api/v1/system/active-config", bytes.NewBufferString(activePayload))
	activeReq.Header.Set("Authorization", "Bearer "+accessToken)
	activeW := httptest.NewRecorder()
	testRouter.ServeHTTP(activeW, activeReq)
	assert.Equal(t, http.StatusOK, activeW.Code)

	// 4. Start V2Ray
	startReq, _ := http.NewRequest(http.MethodPost, "/api/v1/v2ray/start", nil)
	startReq.Header.Set("Authorization", "Bearer "+accessToken)
	startW := httptest.NewRecorder()
	testRouter.ServeHTTP(startW, startReq)
	assert.Equal(t, http.StatusOK, startW.Code)

	// 5. Check status is now running
	statusReq2, _ := http.NewRequest(http.MethodGet, "/api/v1/v2ray/status", nil)
	statusReq2.Header.Set("Authorization", "Bearer "+accessToken)
	statusW2 := httptest.NewRecorder()
	testRouter.ServeHTTP(statusW2, statusReq2)
	assert.Equal(t, http.StatusOK, statusW2.Code)
	json.Unmarshal(statusW2.Body.Bytes(), &statusResponse)
	assert.Equal(t, "running", statusResponse["status"])

	// 6. Stop V2Ray
	stopReq, _ := http.NewRequest(http.MethodPost, "/api/v1/v2ray/stop", nil)
	stopReq.Header.Set("Authorization", "Bearer "+accessToken)
	stopW := httptest.NewRecorder()
	testRouter.ServeHTTP(stopW, stopReq)
	assert.Equal(t, http.StatusOK, stopW.Code)

	// 7. Check status is now stopped
	statusReq3, _ := http.NewRequest(http.MethodGet, "/api/v1/v2ray/status", nil)
	statusReq3.Header.Set("Authorization", "Bearer "+accessToken)
	statusW3 := httptest.NewRecorder()
	testRouter.ServeHTTP(statusW3, statusReq3)
	assert.Equal(t, http.StatusOK, statusW3.Code)
	json.Unmarshal(statusW3.Body.Bytes(), &statusResponse)
	assert.Equal(t, "stopped", statusResponse["status"])
}
