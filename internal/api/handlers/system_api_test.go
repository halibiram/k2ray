package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Note: The setup for this test is now handled by the TestMain function
// in v2ray_test.go, which is part of the same package (handlers_test).
// It initializes a global 'testRouter'.

func TestSystemStatusAPI(t *testing.T) {
	// The 'testRouter' is a global variable initialized in v2ray_test.go's TestMain.
	assert.NotNil(t, testRouter, "testRouter should be initialized by TestMain")

	// Create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/system/status", nil)

	// Serve request
	testRouter.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	expectedResponse := map[string]string{"status": "ok"}
	assert.Equal(t, expectedResponse, response)
}