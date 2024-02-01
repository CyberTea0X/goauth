package controllers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/CyberTea0X/goauth/src/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const loginPath = "/api/login"

func setupLoginTest(t *testing.T) (*gin.Engine, *models.ClientMock, *PublicController, *httptest.ResponseRecorder, string) {
	client := models.NewClientMock()
	router, controller := SetupTestRouter(t, client)
	w := httptest.NewRecorder()
	u, err := url.Parse(loginPath)
	if err != nil {
		t.Fatal(err)
	}
	q := u.Query()
	q.Add("username", "test")
	q.Add("password", "PASSWORD")
	q.Add("email", "test@example.com")
	q.Add("device_id", "123")
	u.RawQuery = q.Encode()
	return router, client, controller, w, u.String()
}

func TestLoginSucceed(t *testing.T) {
	// Fakelogin function fails the test if Login failed
	_, controller, _ := FakeLogin(t)
	defer models.TruncateDatabase(controller.DB)
}

func TestLoginServiceError(t *testing.T) {
	const errMsg = "example"
	const errStatus = http.StatusUnauthorized
	router, client, controller, w, address := setupLoginTest(t)
	client.Engine.GET(controller.LoginServiceURL.Path, func(c *gin.Context) {
		c.JSON(errStatus, gin.H{"error": errMsg})
	})
	r, _ := http.NewRequest("GET", address, nil)
	router.ServeHTTP(w, r)
	res := w.Result()
	defer res.Body.Close()
	err := models.ErrFromResponse(res)
	if err == nil {
		t.Fatal("Failed to get error from response")
	}
	assert.Equal(t, errStatus, res.StatusCode)
	assert.Equal(t, errMsg, err.Error())
}
