package controllers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/CyberTea0X/goauth/models"
	"github.com/CyberTea0X/goauth/models/token"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type RefreshOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

func (p *PublicController) Refresh(c *gin.Context) {
	inputToken := c.Query("token")
	if inputToken == "" {
		c.JSON(http.StatusBadRequest, models.ErrToMap(models.ErrNoTokenSpecified))
		return
	}

	refreshClaims, err := token.RefreshFromString(inputToken, p.RefreshTokenCfg.Secret)

	if errors.Is(err, jwt.ErrTokenExpired) {
		c.JSON(http.StatusUnauthorized, models.ErrToMap(models.ErrTokenExpired))
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrToMap(models.ErrInvalidToken))
		return
	}

	exists, err := refreshClaims.Exists(p.DB)
	if err != nil {
		log.Println("Error while checking if refresh token exists: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if exists == false {
		c.JSON(http.StatusBadRequest, models.ErrToMap(models.ErrTokenNotExists))
		return
	}

	expiresAt := time.Now().Add(time.Hour * time.Duration(p.RefreshTokenCfg.LifespanHour))
	refreshClaims.ExpiresAt = jwt.NewNumericDate(expiresAt)
	expiresUnix := refreshClaims.ExpiresAt.Unix()

	_, err = refreshClaims.Update(p.DB, uint64(expiresUnix))

	if err != nil {
		log.Println("Error updating refresh token identifier in the database: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	refreshToken, err := refreshClaims.TokenString(p.RefreshTokenCfg.Secret)

	if err != nil {
		log.Println("Error generating refresh token from old refresh token: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	expiresAt = time.Now().Add(time.Minute * time.Duration(p.AccessTokenCfg.LifespanMinute))
	accessClaims := token.NewAccess(refreshClaims.UserID, refreshClaims.Roles, expiresAt)
	accessToken, err := accessClaims.TokenString(p.AccessTokenCfg.Secret)

	if err != nil {
		log.Println("Error generating access token from old refresh token: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, RefreshOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresUnix,
	})

}
