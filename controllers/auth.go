package controllers

import (
	"errors"
	"net/http"

	"github.com/CyberTea0X/goauth/models"
	"github.com/CyberTea0X/goauth/models/token"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (p *PublicController) Auth(c *gin.Context) {
	accessToken, err := token.ExtractToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrToMap(models.ErrNoTokenSpecified))
		return
	}

	_, err = token.AccessFromString(accessToken, p.AccessTokenCfg.Secret)

	if errors.Is(err, jwt.ErrTokenExpired) {
		c.JSON(http.StatusUnauthorized, models.ErrToMap(models.ErrTokenExpired))
		return
	}

	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrToMap(models.ErrInvalidToken))
		return
	}

	c.Status(http.StatusOK)
}
