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

//	@title			Golang Authentication microservice
//	@version		0.01
//	@description	Simple authentication service
//	@termsOfService	http://swagger.io/terms/

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
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

	client := &http.Client{Timeout: time.Second * 15}

	controller := controllers.NewController(config.Tokens, config.Services, client, db)

	port := "8080"
	fmt.Println("Auth server starting on port", port)
	controllers.SetupRouter(controller).Run(":" + port)
}
