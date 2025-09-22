package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestSystemStatus verifies that the SystemStatus handler returns the correct
// status code and JSON payload.
func TestSystemStatus(t *testing.T) {
	// Set Gin to test mode to reduce verbose output
	gin.SetMode(gin.TestMode)

	// Create a response recorder to capture the handler's output
	w := httptest.NewRecorder()

	// Create a new test Gin context
	c, _ := gin.CreateTestContext(w)

	// Call the handler function directly
	SystemStatus(c)

	// Assert that the HTTP status code is 200 OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert that the response body is the expected JSON
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)

	assert.NoError(t, err, "Response body should be valid JSON")
	assert.Equal(t, "ok", response["status"], "Response JSON should contain status: ok")
}
