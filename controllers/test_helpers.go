package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CyberTea0X/goauth/src/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetupTestRouter(t *testing.T, client models.HTTPClient) (*gin.Engine, *PublicController) {
	gin.SetMode(gin.ReleaseMode)
	config, err := models.ParseConfig("../config_test.toml")

	if err != nil {
		t.Fatal(err)
	}

	db, err := models.SetupDatabase(&config.Database)

	if err != nil {
		t.Fatal(err)
	}

	controller := NewController(config.Tokens, config.Services, client, db)

	return SetupRouter(controller), controller
}

// Creates fake data both in the database and the returned structure that can be used in tests.
//
// If encounters error fails the test.
// Database should be cleaned up manually after calling this function.
// Controller and router can be used in further tests
func FakeLogin(t *testing.T) (*LoginOutput, *PublicController, *gin.Engine) {
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

	client.Engine.GET(controller.LoginServiceURL.Path, func(c *gin.Context) {
		c.JSON(http.StatusOK, models.User{Id: 1})
	})

	r, _ := http.NewRequest("GET", "/api/login", bytes.NewBuffer(inputJson))
	router.ServeHTTP(w, r)
	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	bodyRaw, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	output := new(LoginOutput)
	err = json.Unmarshal(bodyRaw, output)
	if err != nil {
		t.Fatal(err)
	}
	return output, controller, router
}
