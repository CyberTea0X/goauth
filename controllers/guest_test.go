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

func guestTestSetup(t *testing.T) (*models.ClientMock, *gin.Engine, *PublicController, *httptest.ResponseRecorder) {
	gin.SetMode(gin.ReleaseMode)
	client := models.NewClientMock()
	router, controller, err := SetupTestRouter(client)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	return client, router, controller, w
}

func TestGuestSucceeds(t *testing.T) {
	client, router, controller, w := guestTestSetup(t)
	defer models.TruncateDatabase(controller.DB)
	guest := models.Guest{
		FullName: "Test",
		Id:       1,
	}
	client.Engine.POST(controller.GuestServiceURL.Path, func(c *gin.Context) {
		c.JSON(http.StatusOK, &guest)
	})
	guestInput := GuestInput{
		FullName: guest.FullName,
		DeviceId: 1,
	}
	jsonInput, _ := json.Marshal(guestInput)
	req, _ := http.NewRequest("GET", "/api/guest", bytes.NewReader(jsonInput))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGuestInvalidJSON(t *testing.T) {
	client, router, controller, w := guestTestSetup(t)
	defer models.TruncateDatabase(controller.DB)

	client.Engine.POST(controller.GuestServiceURL.Path, func(c *gin.Context) {
		c.Status(http.StatusBadRequest)
	})

	guestInput := GuestInput{
		DeviceId: 1,
	}

	jsonInput, _ := json.Marshal(guestInput)
	req, _ := http.NewRequest("GET", "/api/guest", bytes.NewReader(jsonInput))
	router.ServeHTTP(w, req)
	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, w.Code)
	err := models.ErrFromResponse(res)
	if err == nil {
		t.Fatal("Failed to get error from response")
	}
	assert.Equal(t, models.ErrInvalidJson.Error(), err.Error())
}

func TestGuestServiceError(t *testing.T) {
	client, router, controller, w := guestTestSetup(t)
	defer models.TruncateDatabase(controller.DB)
	client.Engine.POST(controller.GuestServiceURL.Path, func(c *gin.Context) {
		c.Status(http.StatusUnauthorized)
	})
	guestInput := GuestInput{
		DeviceId: 1,
	}

	jsonInput, _ := json.Marshal(guestInput)
	req, _ := http.NewRequest("GET", "/api/guest", bytes.NewReader(jsonInput))
	router.ServeHTTP(w, req)
	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, w.Code)
	err := models.ErrFromResponse(res)
	if err == nil {
		t.Fatal("Failed to get error from response")
	}
	assert.Equal(t, models.ErrInvalidJson.Error(), err.Error())
}
