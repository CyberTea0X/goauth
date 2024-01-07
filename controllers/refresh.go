package controllers

import (
	"errors"
	"net/http"

	"github.com/CyberTea0X/delta_art/src/backend/models"
	"github.com/CyberTea0X/delta_art/src/backend/utils/token"
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
    old_token, err := token.RefreshParse(input.RefreshToken)

    if err != nil {
		c.JSON(http.StatusBadRequest, err)
        return 
    }

    user_id, err := token.ExtractUint(old_token, "user_id")

    if err != nil {
        c.JSON(http.StatusBadRequest, err)
        return 
    }

    device_id, err := token.ExtractUint(old_token, "device_id")

    if err != nil {
        c.JSON(http.StatusBadRequest, err)
        return 
    }

    if models.RefreshTokenExists(p.DB, user_id, device_id) == false {
        c.JSON(http.StatusBadRequest, errors.New("token does not exist"))
        return
    }

    models.DeleteOldTokens(p.DB, user_id, device_id)

    refresh_token, err := token.GenerateRefresh(user_id, device_id)

    if err != nil {
		c.JSON(http.StatusBadRequest, err)
        return 
    }

    access_token, expires_at, err := token.GenerateAccessToken(user_id)

    refresh_model := models.RefreshToken{}
    refresh_model.Token = refresh_token
    refresh_model.UserID = user_id
    refresh_model.DeviceID = device_id

    if _, err := refresh_model.SaveToken(p.DB); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    }

    if err != nil {
		c.JSON(http.StatusBadRequest, err)
        return 
    }
    c.JSON(http.StatusOK, gin.H{
        "accessToken": access_token,
        "refreshToken": refresh_token,
        "expires_at": expires_at,
    })

}
