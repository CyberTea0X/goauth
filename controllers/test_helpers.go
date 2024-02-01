package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/CyberTea0X/goauth/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetupTestRouter(t *testing.T, controller *PublicController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return SetupRouter(controller)
}

func SetupTestController(t *testing.T, client models.HTTPClient) *PublicController {
	config, err := models.ParseConfig("../config_test.toml")

	if err != nil {
		t.Fatal(err)
	}

	db, err := models.SetupDatabase(&config.Database)

	if err != nil {
		t.Fatal(err)
	}

	return NewController(config.Tokens, config.Services, client, db)
}

// Creates fake data both in the database and the returned structure that can be used in tests.
//
// If encounters error fails the test.
// Database should be cleaned up manually after calling this function.
// Controller and router can be used in further tests
func FakeLogin(t *testing.T) (*LoginOutput, *PublicController, *gin.Engine) {
	client := models.NewClientMock()
	controller := SetupTestController(t, client)
	router := SetupTestRouter(t, controller)
	u, err := url.Parse("/api/login")
	if err != nil {
		t.Fatal(err)
	}
	q := u.Query()
	q.Add("username", "test")
	q.Add("password", "PASSWORD")
	q.Add("email", "test@example.com")
	q.Add("device_id", "123")
	u.RawQuery = q.Encode()
	w := httptest.NewRecorder()

	client.Engine.GET(controller.LoginServiceURL.Path, func(c *gin.Context) {
		c.JSON(http.StatusOK, models.User{Id: 1})
	})

	r, _ := http.NewRequest("GET", u.String(), nil)
	router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	bodyRaw, err := io.ReadAll(w.Body)
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
