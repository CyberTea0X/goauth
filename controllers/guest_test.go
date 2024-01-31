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

func TestGuestSucceeds(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	client := new(models.ClientMock)
	guest := models.Guest{
		FullName: "Test",
	}
	jsonGuest, _ := json.Marshal(guest)
	client.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader(jsonGuest)),
	}
	guestInput := GuestInput{
		FullName: guest.FullName,
		DeviceId: 1,
	}
	router, controller, err := SetupTestRouter(client)
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
	client := new(models.ClientMock)
	guest := models.Guest{
		FullName: "Test",
	}
	jsonGuest, _ := json.Marshal(guest)
	client.Response = &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(bytes.NewReader(jsonGuest)),
	}
	guestInput := GuestInput{
		DeviceId: 1,
	}
	router, controller, err := SetupTestRouter(client)
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
