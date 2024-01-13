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
	c := controllers.Setup()

	port := "8080"
	router := gin.Default()

	public := router.Group("api")
	public.GET("health_check", c.HealthCheck)
	public.POST("login", c.Login)
	public.POST("refresh", c.Refresh)

	// protected := router.Group("api")
	// protected.Use(middlewares.JwtAuthMiddleware())
	// protected.GET("/user", c.CurrentUser)
	fmt.Println("goauth api server starting on port ", port)
	debug := os.Getenv("DEBUG")
	if debug != "" {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		fmt.Println("Debug log enabled")
	}

	router.Run(":" + port)
}
