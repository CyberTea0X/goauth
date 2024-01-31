package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/CyberTea0X/goauth/src/backend/models"
	"github.com/CyberTea0X/goauth/src/backend/models/token"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Structure describing the json fields that should be in the refresh request
type RefreshInput struct {
	RefreshToken string `json:"token" binding:"required"`
}

func (p *PublicController) Refresh(c *gin.Context) {
	var input RefreshInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrToMap(models.ErrInvalidJson))
		return
	}

	refreshClaims, err := token.RefreshFromString(input.RefreshToken, p.RefreshTokenCfg.Secret)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if refreshClaims.ExpiresAt.Time.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token expired"})
		return
	}

	exists, err := refreshClaims.Exists(p.DB)
	if err != nil {
		log.Println("Error while checking if refresh token exists: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if exists == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token does not exist"})
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
	accessClaims := token.NewAccess(refreshClaims.UserID, refreshClaims.Role, expiresAt)
	accessToken, err := accessClaims.TokenString(p.AccessTokenCfg.Secret)

	if err != nil {
		log.Println("Error generating access token from old refresh token: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"expires_at":   expiresUnix,
	})

}
