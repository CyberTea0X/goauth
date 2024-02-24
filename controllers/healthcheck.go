package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// HealthCheck checks if the server is healthy
func (p *PublicController) HealthCheck(c *gin.Context) {
	// TODO: implement https://pkg.go.dev/github.com/heptiolabs/healthcheck
	c.Status(http.StatusOK)
}
