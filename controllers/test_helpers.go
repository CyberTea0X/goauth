package controllers

import (
	"github.com/CyberTea0X/goauth/src/backend/models"
	"github.com/gin-gonic/gin"
)

func SetupTestRouter(client models.HTTPClient) (*gin.Engine, *PublicController, error) {
	config, err := models.ParseConfig("../config_test.toml")

	if err != nil {
		return nil, nil, err
	}

	db, err := models.SetupDatabase(&config.Database)

	if err != nil {
		return nil, nil, err
	}

	controller := NewController(config.Tokens, config.Services, client, db)

	return SetupRouter(controller), controller, nil
}