package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (p *PublicController) HealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "Server is alive!")
}
