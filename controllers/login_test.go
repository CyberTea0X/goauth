package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CyberTea0X/goauth/src/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupLoginTest(t *testing.T) (*gin.Engine, *models.ClientMock, *PublicController, *httptest.ResponseRecorder, []byte) {
	client := models.NewClientMock()
	router, controller := SetupTestRouter(t, client)
	input := LoginInput{
		Username: "test",
		Password: "PASSWORD",
		Email:    "EMAIL",
		DeviceId: 1,
	}
	w := httptest.NewRecorder()
	inputJson, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	return router, client, controller, w, inputJson
}

func TestLoginSucceed(t *testing.T) {
	// Fakelogin function fails the test if Login failed
	_, controller, _ := FakeLogin(t)
	defer models.TruncateDatabase(controller.DB)
}

func TestLoginServiceError(t *testing.T) {
	const errMsg = "example"
	const errStatus = http.StatusUnauthorized
	router, client, controller, w, inputJson := setupLoginTest(t)
	client.Engine.GET(controller.LoginServiceURL.Path, func(c *gin.Context) {
		c.JSON(errStatus, gin.H{"error": errMsg})
	})
	r, _ := http.NewRequest("GET", "/api/login", bytes.NewBuffer(inputJson))
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
