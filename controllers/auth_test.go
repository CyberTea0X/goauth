package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/CyberTea0X/goauth/src/backend/models"
	"github.com/CyberTea0X/goauth/src/backend/models/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func authTestSetup(t *testing.T) (*gin.Engine, *PublicController, *http.Request, *httptest.ResponseRecorder) {
	gin.SetMode(gin.ReleaseMode)
	router, controller, err := SetupTestRouter(nil)

	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/auth", nil)
	return router, controller, req, w
}

func TestAuthSucceed(t *testing.T) {
	router, controller, req, w := authTestSetup(t)

	accessClaims := token.NewAccess(123, "guest", time.Now().Add(time.Hour))
	accessToken, err := accessClaims.TokenString(controller.AccessTokenCfg.Secret)

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", accessToken)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthNoToken(t *testing.T) {
	router, _, req, w := authTestSetup(t)

	router.ServeHTTP(w, req)
	res := w.Result()
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	expected, _ := json.Marshal(models.ErrToMap(models.ErrNoTokenSpecified))
	body, _ := io.ReadAll(res.Body)
	assert.JSONEq(t, string(expected), string(body))
}

func TestAuthInvalidToken(t *testing.T) {
	router, controller, req, w := authTestSetup(t)
	accessClaims := token.NewAccess(123, "guest", time.Now().Add(time.Hour))
	accessToken, err := accessClaims.TokenString(controller.AccessTokenCfg.Secret + "123")

	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", accessToken)

	router.ServeHTTP(w, req)
	res := w.Result()
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	expected, _ := json.Marshal(models.ErrToMap(models.ErrInvalidToken))
	body, _ := io.ReadAll(res.Body)
	assert.JSONEq(t, string(expected), string(body))
}

func TestAuthTokenExpired(t *testing.T) {
	router, controller, req, w := authTestSetup(t)
	accessClaims := token.NewAccess(123, "guest", time.Now().Add(-time.Hour))
	accessToken, err := accessClaims.TokenString(controller.AccessTokenCfg.Secret)

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", accessToken)

	router.ServeHTTP(w, req)
	res := w.Result()
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	expected, _ := json.Marshal(models.ErrToMap(models.ErrTokenExpired))
	body, _ := io.ReadAll(res.Body)
	assert.JSONEq(t, string(expected), string(body))
	w.Flush()
}
