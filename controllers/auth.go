package controllers

import (
	"errors"
	"net/http"

	"github.com/CyberTea0X/goauth/src/backend/models/token"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (p *PublicController) Auth(c *gin.Context) {
	accessToken, err := token.ExtractToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrToMap(ErrNoTokenSpecified{}))
		return
	}

	_, err = token.AccessFromString(accessToken, p.AccessTokenCfg.Secret)

	if errors.Is(err, jwt.ErrTokenExpired) {
		c.JSON(http.StatusUnauthorized, ErrToMap(ErrTokenExpired{}))
		return
	}

	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrToMap(ErrInvalidAccessToken{}))
		return
	}

	c.Status(http.StatusOK)
}
