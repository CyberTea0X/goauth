package controllers

import (
	"database/sql"
	"log"

	"github.com/CyberTea0X/goauth/src/backend/models"
	"github.com/caarlos0/env/v10"
	// Import mysql driver for sq.Open to work
	_ "github.com/go-sql-driver/mysql"
	// Autoloading of environment variables from .env file
	_ "github.com/joho/godotenv/autoload"
)

type PublicController struct {
	DB              *sql.DB
	AccessTokenCfg  models.AccessTokenCfg
	RefreshTokenCfg models.RefreshTokenCfg
}

func Setup() *PublicController {
	pCtrl := new(PublicController)
	cfg := new(models.DatabaseConfig)

	err := env.Parse(&pCtrl.AccessTokenCfg)
	if err != nil {
		log.Fatal("Error while parsing access token config from enviroment: ", err.Error())
	}

	err = env.Parse(&pCtrl.RefreshTokenCfg)
	if err != nil {
		log.Fatal("Error while parsing refresh token config from enviroment: ", err.Error())
	}

	err = env.Parse(cfg)

	if err != nil {
		log.Fatal("Error while parsing database config from enviroment: ", err.Error())
	}

	db, err := sql.Open(cfg.Dbdriver, cfg.GetUrl())

	if err != nil {
		log.Fatal("Error connecting to the database: ", err.Error())
	}

	pCtrl.DB = db
	return pCtrl
}
