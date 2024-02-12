package controllers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/CyberTea0X/goauth/models"
	"github.com/CyberTea0X/goauth/models/token"
	"github.com/gin-gonic/gin"
)

type GuestInput struct {
	FullName string `json:"full_name" binding:"required"`
	DeviceId uint   `json:"device_id" binding:"required"`
}

type GuestOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

const GUEST_ROLE = "guest"

// Guest registers new guest and returns tokens
func (p *PublicController) Guest(c *gin.Context) {

	var input GuestInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrToMap(models.ErrInvalidJson))
		return
	}

	guest, err := models.RegisterGuest(input.FullName, p.GuestServiceURL, p.Client)

	if err != nil {
		targetErr := new(models.ExternalServiceError)
		if errors.As(err, &targetErr) {
			c.JSON(targetErr.Status, models.ErrToMap(targetErr))
		} else {
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	expiresAt := time.Now().Add(p.RefreshTokenCfg.Lifespan())
	refreshClaims := token.NewRefresh(-1, input.DeviceId, guest.Id, []string{GUEST_ROLE}, expiresAt)

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

	expiresAt = time.Now().Add(p.AccessTokenCfg.Lifespan())
	accessClaims := token.NewAccess(guest.Id, []string{GUEST_ROLE}, expiresAt)
	accessToken, err := accessClaims.TokenString(p.AccessTokenCfg.Secret)

	if err != nil {
		log.Println("Error generating access token: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, GuestOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt.Unix(),
	})

}
