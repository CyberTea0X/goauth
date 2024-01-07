package controllers

import (
	"github.com/CyberTea0X/delta_art/src/backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Structure describing the json fields that should be in the login request
type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password" binding:"required"`
	Email string `json:"email"`
	DeviceId uint `json:"device_id" binding:"required"`
}

// Function that is responsible for user authorization.
// In response to a successful authorization request, returns 
// Access Token and Refresh Token, as well as the time of death of the Access Token
func (p *PublicController) Login(c *gin.Context) {

	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	u.Username = input.Username
	u.Password = input.Password
    u.Email = input.Email

	tokens, err := models.Login(p.DB, u, input.DeviceId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

    c.JSON(http.StatusOK, gin.H{
        "accessToken": tokens.AccessToken,
        "refreshToken": tokens.RefreshToken,
        "expires_at": tokens.AccessTokenExpires,
    })

}
