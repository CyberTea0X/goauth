package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/CyberTea0X/goauth/src/backend/models"
	"github.com/CyberTea0X/goauth/src/backend/models/token"
	"github.com/gin-gonic/gin"
)

type GuestInput struct {
	FullName string `json:"full_name" binding:"required"`
	DeviceId uint   `json:"device_id" binding:"required"`
}

type GuestOutput struct {
	AccessToken  string `json:"access_token" example:"token"`
	RefreshToken string `json:"refresh_token" example:"token"`
	ExpiresAt    int64  `json:"expires_at" example:"244534234"`
	Role         string `json:"role" example:"root"`
}

const ROLE = "guest"

// Guest authorizes user as guest
//
//	@Summary		guest authorization
//	@Description	authorizes user as guest
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	GuestOutput
//	@Schemes
//	@Router	/guest [get]
func (p *PublicController) Guest(c *gin.Context) {

	var input GuestInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	guest := new(models.Guest)
	guest.FullName = input.FullName
	id, err := guest.InsertToDb(p.DB)
	if err != nil {
		log.Println("Error inserting guest into db", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	guest.Id = id

	expiresAt := time.Now().Add(time.Hour * time.Duration(p.RefreshTokenCfg.LifespanHour))
	refreshClaims := token.NewRefresh(-1, input.DeviceId, guest.Id, ROLE, expiresAt)

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
	accessClaims := token.NewAccess(guest.Id, ROLE, expiresAt)
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
		Role:         ROLE,
	})

}
