package main

import (
	"fmt"
	"github.com/CyberTea0X/delta_art/src/backend/controllers"
	"github.com/CyberTea0X/delta_art/src/backend/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {
    c := controllers.Setup()

	port := "8080"
	router := gin.Default()

	public := router.Group("api")
	public.GET("health_check", c.HealthCheck)
	public.POST("register", c.Register)
	public.POST("login", c.Login)
	public.POST("refresh", c.Refresh)

	protected := router.Group("api")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/user", c.CurrentUser)
	fmt.Println("delta_art api server starting on port ", port)

	router.Run(":" + port)
}
