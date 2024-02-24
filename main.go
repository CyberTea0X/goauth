package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/CyberTea0X/goauth/controllers"
	"github.com/CyberTea0X/goauth/models"

	// "github.com/CyberTea0X/goauth/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {
	debug := os.Getenv("DEBUG")
	if debug != "" && debug != "0" {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		fmt.Println("Debug log enabled")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	config, err := models.ParseConfig("config.toml")

	if err != nil {
		log.Fatal(err)
	}

	db, err := models.SetupDatabase(&config.Database)

	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{Timeout: time.Second * time.Duration(config.App.TimeoutSeconds)}

	controller := controllers.NewController(config.Tokens, config.Services, client, db)

	port := config.App.Port
	fmt.Println("Auth server starting on port", port)
	controllers.SetupRouter(controller).Run(":" + port)
}
