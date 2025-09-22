package handlers_test

import (
	"bytes"
	"encoding/json"
	"k2ray/internal/api"
	"k2ray/internal/config"
	"k2ray/internal/db"
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

// TestMain sets up the test environment once for the entire package.
func TestMain(m *testing.M) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Create a temporary database file for the test suite
	tmpfile, err := os.CreateTemp("", "test_auth_*.db")
	if err != nil {
		log.Fatalf("Failed to create temp db file: %v", err)
	}
	dbPath := tmpfile.Name()
	tmpfile.Close()

	// Load config and override with test-specific values
	config.LoadConfig("") // Load .env or environment, which we will override
	config.AppConfig.DatabaseURL = dbPath
	config.AppConfig.JWTSecret = "a-secure-test-secret"

	// Initialize DB, run migrations, and seed with a test user
	db.InitDB()
	db.RunMigrations()
	createTestUser("testuser", "password123")

	// Setup the router that all tests will use
	testRouter = gin.Default()
	api.SetupRouter(testRouter)

	// Run the tests
	code := m.Run()

	// Teardown: close DB connection and remove the temp file
	db.DB.Close()
	os.Remove(dbPath)

	os.Exit(code)
}

// createTestUser is a helper to seed the test database.
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

// TestLoginHandler contains all sub-tests for the login endpoint.
func TestLoginHandler(t *testing.T) {

	t.Run("Successful Login", func(t *testing.T) {
		loginPayload := `{"username": "testuser", "password": "password123"}`
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(loginPayload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "access_token", "Response should contain an access token")
		assert.Contains(t, response, "refresh_token", "Response should contain a refresh token")
	})

	t.Run("Invalid Password", func(t *testing.T) {
		loginPayload := `{"username": "testuser", "password": "wrongpassword"}`
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(loginPayload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Non-existent User", func(t *testing.T) {
		loginPayload := `{"username": "nosuchuser", "password": "password123"}`
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(loginPayload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Bad Request Payload", func(t *testing.T) {
		loginPayload := `{"username": "testuser"}` // Missing password field
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(loginPayload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
