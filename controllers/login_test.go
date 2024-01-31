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

func TestLoginSucceed(t *testing.T) {
	client := models.NewClientMock()
	router, controller := SetupTestRouter(t, client)
	input := LoginInput{
		Username: "test",
		Password: "PASSWORD",
		Email:    "EMAIL",
		DeviceId: 1,
	}

	client.Engine.GET(controller.LoginServiceURL.Path, func(c *gin.Context) {
		c.JSON(http.StatusOK, models.User{Id: 1})
	})

	w := httptest.NewRecorder()
	inputJson, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	r, _ := http.NewRequest("GET", "/api/login", bytes.NewBuffer(inputJson))
	router.ServeHTTP(w, r)
	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
