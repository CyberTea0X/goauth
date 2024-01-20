package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/CyberTea0X/goauth/src/backend/models"
	"github.com/CyberTea0X/goauth/src/backend/models/token"
	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml/v2"

	// Import mysql driver for sq.Open to work
	_ "github.com/go-sql-driver/mysql"
)

type PublicController struct {
	DB              *sql.DB
	AccessTokenCfg  models.AccessTokenCfg
	RefreshTokenCfg models.RefreshTokenCfg
}

func Setup() *PublicController {

	debug := os.Getenv("DEBUG")
	if debug != "" && debug != "0" {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		fmt.Println("Debug log enabled")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	pCtrl := new(PublicController)

	configFile, err := os.ReadFile("config.toml")

	if err != nil {
		log.Fatal("Config file not found: ", err.Error())
	}

	tomlConfig := new(models.TomlConfig)

	err = toml.Unmarshal(configFile, tomlConfig)

	pCtrl.AccessTokenCfg = tomlConfig.TokensCfg.AccessTokenCfg
	pCtrl.RefreshTokenCfg = tomlConfig.TokensCfg.RefreshTokenCfg

	if err != nil {
		log.Fatal("Error parsing config file: ", err.Error())
	}

	db, err := sql.Open(tomlConfig.Database.Driver, tomlConfig.Database.GetUrl())

	if err != nil {
		log.Fatal("Error connecting to the database: ", err.Error())
	}
	if err = token.CreateRefreshTable(db); err != nil {
		log.Fatal("Failed to create refresh token table: ", err.Error())
	}
	if err = models.CreateGuestTable(db); err != nil {
		log.Fatal("Failed to create guests table: ", err.Error())
	}

	pCtrl.DB = db
	return pCtrl
}
