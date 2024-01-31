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

func TestGuestSucceeds(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	client := models.NewClientMock()
	router, controller, err := SetupTestRouter(client)
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
	defer models.TruncateDatabase(controller.DB)

	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	jsonInput, _ := json.Marshal(guestInput)
	req, _ := http.NewRequest("GET", "/api/guest", bytes.NewReader(jsonInput))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGuestFails(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	client := models.NewClientMock()
	router, controller, err := SetupTestRouter(client)
	client.Engine.POST(controller.GuestServiceURL.Path, func(c *gin.Context) {
		c.Status(http.StatusBadRequest)
	})
	guestInput := GuestInput{
		DeviceId: 1,
	}
	defer models.TruncateDatabase(controller.DB)

	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	jsonInput, _ := json.Marshal(guestInput)
	req, _ := http.NewRequest("GET", "/api/guest", bytes.NewReader(jsonInput))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
