package controllers

import (
	"testing"

	"github.com/CyberTea0X/goauth/src/backend/models"
	"github.com/gin-gonic/gin"
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
