package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// @BasePath /api

// HealthCheck godoc
// @Summary checks if server is running
// @Schemes
// @Description to healthcheck
// @Tags healthcheck
// @Success 200 {string} Server is alive!
// @Router /health_check [get]
func (p *PublicController) HealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "Server is alive!")
}
