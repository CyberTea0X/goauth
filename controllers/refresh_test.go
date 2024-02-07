package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/CyberTea0X/goauth/models"
	"github.com/CyberTea0X/goauth/models/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const refreshPath = "/api/refresh"

// Generates refresh token for testing purposes
func generateTestRefresh(t *testing.T, refresh string, router *gin.Engine) *RefreshOutput {
	w := httptest.NewRecorder()
	url := url.URL{
		Host:   "test",
		Path:   refreshPath,
		Scheme: "http",
	}
	q := url.Query()
	q.Add("token", refresh)
	url.RawQuery = q.Encode()
	req, _ := http.NewRequest("GET", url.String(), nil)
	router.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal("Failed to read refresh responce body")
	}
	output := new(RefreshOutput)
	if err := json.Unmarshal(rawBody, output); err != nil {
		t.Fatal("Failed to parse refresh token output from body", err.Error())
	}
	return output
}

func TestRefreshSucceds(t *testing.T) {
	l, p, router := FakeLogin(t)
	defer models.TruncateDatabase(p.DB)
	refreshed := generateTestRefresh(t, l.RefreshToken, router)
	r := generateTestRefresh(t, refreshed.RefreshToken, router)
	if r.AccessToken == "" || r.RefreshToken == "" || r.ExpiresAt == 0 {
		t.Fatal("Failed to refresh from refreshed login output: ", r)
	}
}

// if controller or token is nil, token not added to query
func testRequestRefresh(t *testing.T, router *gin.Engine, p *PublicController, address string, token *token.RefreshToken) *http.Response {
	w := httptest.NewRecorder()
	if token != nil && p != nil {
		tokenString, err := token.TokenString(p.RefreshTokenCfg.Secret)
		if err != nil {
			t.Fatal("Failed to generate refresh token for testing")
		}
		u, err := url.Parse(address)
		if err != nil {
			t.Fatal("Failed to parse url", err.Error())
		}
		q := u.Query()
		q.Add("token", tokenString)
		u.RawQuery = q.Encode()
		address = u.String()
	}
	req, _ := http.NewRequest("GET", address, nil)
	router.ServeHTTP(w, req)
	res := w.Result()
	return res
}

func TestRefreshNoToken(t *testing.T) {
	p := SetupTestController(t, http.DefaultClient)
	router := SetupTestRouter(t, p)
	res := testRequestRefresh(t, router, nil, refreshPath, nil)
	defer res.Body.Close()
	err := models.ErrFromResponse(res)
	if err == nil {
		t.Fatal("Failed to get error from response")
	}
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, models.ErrNoTokenSpecified.Error(), err.Error())
}

func TestRefreshExpired(t *testing.T) {
	p := SetupTestController(t, http.DefaultClient)
	router := SetupTestRouter(t, p)
	claims := token.NewRefresh(123, 123, 123, []string{"test"}, time.Now().Add(-time.Hour))
	res := testRequestRefresh(t, router, p, refreshPath, claims)
	defer res.Body.Close()
	err := models.ErrFromResponse(res)
	if err == nil {
		t.Fatal("Failed to get error from response")
	}
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assert.Equal(t, models.ErrTokenExpired.Error(), err.Error())
}

func TestRefreshNotExists(t *testing.T) {
	p := SetupTestController(t, http.DefaultClient)
	router := SetupTestRouter(t, p)
	claims := token.NewRefresh(123, 123, 123, []string{"test"}, time.Now().Add(time.Hour))
	res := testRequestRefresh(t, router, p, refreshPath, claims)
	defer res.Body.Close()
	err := models.ErrFromResponse(res)
	if err == nil {
		t.Fatal("Failed to get error from response")
	}
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, models.ErrTokenNotExists.Error(), err.Error())
}
