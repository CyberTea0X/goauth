package main

import (
	"fmt"
	"log"
	"os"

	"github.com/CyberTea0X/goauth/src/backend/controllers"

	// "github.com/CyberTea0X/goauth/src/backend/middlewares"
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

	c := controllers.Setup()

	port := "8080"
	router := gin.Default()

	public := router.Group("api")
	public.GET("health_check", c.HealthCheck)
	public.GET("login", c.Login)
	public.GET("refresh", c.Refresh)
	public.GET("auth", c.Auth)

	fmt.Println("Auth server starting on port", port)

	router.Run(":" + port)
}
