package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(NewHealthHandler().HealthHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Status code is not OK")

	expected := `{"status": "ok"}`
	assert.Equal(t, expected, rr.Body.String(), "Response body differs")
}
