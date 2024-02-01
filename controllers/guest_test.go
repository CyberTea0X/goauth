package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CyberTea0X/goauth/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func guestTestSetup(t *testing.T) (*models.ClientMock, *gin.Engine, *PublicController, *httptest.ResponseRecorder) {
	gin.SetMode(gin.ReleaseMode)
	client := models.NewClientMock()
	p := SetupTestController(t, client)
	router := SetupTestRouter(t, p)
	w := httptest.NewRecorder()
	return client, router, p, w
}

func TestGuestSucceeds(t *testing.T) {
	client, router, p, w := guestTestSetup(t)
	defer models.TruncateDatabase(p.DB)
	guest := models.Guest{
		FullName: "Test",
		Id:       1,
	}
	client.Engine.POST(p.GuestServiceURL.Path, func(c *gin.Context) {
		c.JSON(http.StatusOK, &guest)
	})
	guestInput := GuestInput{
		FullName: guest.FullName,
		DeviceId: 1,
	}
	jsonInput, _ := json.Marshal(guestInput)
	req, _ := http.NewRequest("POST", "/api/guest", bytes.NewReader(jsonInput))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGuestInvalidJSON(t *testing.T) {
	client, router, p, w := guestTestSetup(t)
	defer models.TruncateDatabase(p.DB)

	client.Engine.POST(p.GuestServiceURL.Path, func(c *gin.Context) {
		c.Status(http.StatusBadRequest)
	})

	guestInput := GuestInput{
		DeviceId: 1,
	}

	jsonInput, _ := json.Marshal(guestInput)
	req, _ := http.NewRequest("POST", "/api/guest", bytes.NewReader(jsonInput))
	router.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, w.Code)
	err := models.ErrFromResponse(res)
	if err == nil {
		t.Fatal("Failed to get error from response")
	}
	assert.Equal(t, models.ErrInvalidJson.Error(), err.Error())
}

func TestGuestServiceError(t *testing.T) {
	client, router, p, w := guestTestSetup(t)
	defer models.TruncateDatabase(p.DB)
	const errMsg = "Unauthorized"
	const status = http.StatusUnauthorized
	client.Engine.POST(p.GuestServiceURL.Path, func(c *gin.Context) {
		c.JSON(status, gin.H{"error": errMsg})
	})
	guestInput := GuestInput{
		FullName: "Guest",
		DeviceId: 1,
	}

	jsonInput, _ := json.Marshal(guestInput)
	req, _ := http.NewRequest("POST", "/api/guest", bytes.NewReader(jsonInput))
	router.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, status, res.StatusCode)
	err := models.ErrFromResponse(res)
	if err == nil {
		t.Fatal("Failed to get error from response")
	}
	assert.Equal(t, errMsg, err.Error())
}
