package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/CyberTea0X/goauth/src/backend/models/token"
	"github.com/gin-gonic/gin"
)

// Structure describing the json fields that should be in the refresh request
type RefreshInput struct {
	RefreshToken string `json:"token" binding:"required"`
}

func (p *PublicController) Refresh(c *gin.Context) {
	var input RefreshInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	old_token, err := token.RefreshFromString(input.RefreshToken, p.RefreshTokenCfg.Secret)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := old_token.Exists(p.DB)
	if err != nil {
		log.Println("Error while checking if refresh token exists: ", err.Error())
	}

	if exists == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token does not exist"})
		return
	}

	err = token.DeleteOldToken(p.DB, old_token.UserID, old_token.DeviceID)

	if err != nil {
		log.Println("Error deleting old tokens: ", err.Error())
	}

	expiresAt := time.Now().Add(time.Hour * time.Duration(p.RefreshTokenCfg.Lifespan))
	refreshClaims := token.NewRefresh(old_token.DeviceID, old_token.UserID, expiresAt)
	refreshToken, err := refreshClaims.TokenString(p.RefreshTokenCfg.Secret)

	if err != nil {
		log.Println("Error generating refresh token from old token: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	expiresAt = time.Now().Add(time.Minute * time.Duration(p.AccessTokenCfg.Lifespan))
	accessClaims := token.NewAccess(old_token.UserID, expiresAt)
	accessToken, err := accessClaims.TokenString(p.AccessTokenCfg.Secret)

	if err != nil {
		log.Println("Error generating access token from old refresh token: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if _, err := refreshClaims.InsertToDb(p.DB); err != nil {
		log.Println("Error inserting refresh token identifier to the database: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"expires_at":   expiresAt.Unix(),
	})

}
