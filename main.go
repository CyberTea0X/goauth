package main

import (
	"fmt"

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
	public.GET("login", c.Login)
	public.GET("refresh", c.Refresh)
	public.GET("auth", c.Auth)
	public.GET("guest", c.Guest)

	fmt.Println("Auth server starting on port", port)

	router.Run(":" + port)
}
