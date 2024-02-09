package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	p := SetupTestController(t, http.DefaultClient)
	router := SetupTestRouter(t, p)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest("GET", "/v1/health_check", nil))
	assert.Equal(t, http.StatusOK, w.Code)
}
