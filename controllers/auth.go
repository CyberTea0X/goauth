package controllers

import (
	"net/http"
	"time"

	"github.com/CyberTea0X/goauth/src/backend/models/token"
	"github.com/gin-gonic/gin"
)

func (p *PublicController) Auth(c *gin.Context) {
	accessToken, err := token.ExtractToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no token specified"})
		return
	}
	accessClaims, err := token.AccessFromString(accessToken, p.AccessTokenCfg.Secret)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
		return
	}

	if accessClaims.ExpiresAt.Time.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
		return
	}

	c.Status(http.StatusOK)
}
