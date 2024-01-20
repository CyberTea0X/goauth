package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/CyberTea0X/goauth/src/backend/models"
	"github.com/CyberTea0X/goauth/src/backend/models/token"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Structure describing the json fields that should be in the login request
type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
	DeviceId uint   `json:"device_id" binding:"required"`
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

	var err error
	var udb = &models.User{}

	if input.Username != "" {
		udb, err = models.GetUserByUsername(p.DB, input.Username)
	} else {
		udb, err = models.GetUserByEmail(p.DB, input.Email)
	}

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	err = models.VerifyPassword(input.Password, udb.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	} else if err != nil {
		log.Println("Error verifying user password: ", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().Add(time.Hour * time.Duration(p.RefreshTokenCfg.LifespanHour))
	refreshClaims := token.NewRefresh(-1, input.DeviceId, udb.Id, udb.Role, expiresAt)

	refreshId, err := refreshClaims.InsertOrUpdate(p.DB)

	if err != nil {
		log.Println("Error inserting or updating on duplicate key refresh token in the db: ", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	refreshClaims.TokenID = refreshId

	refreshToken, err := refreshClaims.TokenString(p.RefreshTokenCfg.Secret)

	if err != nil {
		log.Println("Error generating refresh token: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	expiresAt = time.Now().Add(time.Minute * time.Duration(p.AccessTokenCfg.LifespanMinute))
	accessClaims := token.NewAccess(udb.Id, udb.Role, expiresAt)
	accessToken, err := accessClaims.TokenString(p.AccessTokenCfg.Secret)

	if err != nil {
		log.Println("Error generating access token: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"expires_at":   expiresAt.Unix(),
		"role":         udb.Role,
	})

}
