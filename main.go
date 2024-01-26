package main

import (
	"fmt"

	"github.com/CyberTea0X/goauth/src/backend/controllers"

	// "github.com/CyberTea0X/goauth/src/backend/middlewares"
	"github.com/gin-gonic/gin"
)

// @title           Golang Authentication microservice
// @version         0.01
// @description     Simple authentication service
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {

	c := controllers.Setup()

	port := "8080"
	router := gin.Default()

	public := router.Group("api")
	public.GET("health_check", c.HealthCheck)
	public.GET("login", c.Login)
	public.GET("refresh", c.Refresh)
	public.GET("auth", c.Auth)
	public.GET("guest", c.Guest)

	fmt.Println("Auth server starting on port", port)

	router.Run(":" + port)
}
