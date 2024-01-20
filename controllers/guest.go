package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/CyberTea0X/goauth/src/backend/models"
	"github.com/CyberTea0X/goauth/src/backend/models/token"
	"github.com/gin-gonic/gin"
)

// Structure describing the json fields that should be in the guest access request
type GuestInput struct {
	FullName string `json:"full_name"`
	DeviceId uint   `json:"device_id" binding:"required"`
}

const ROLE = "guest"

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
	refreshClaims := token.NewRefreshNoID(input.DeviceId, guest.Id, ROLE, expiresAt)

	refreshID, exists, err := refreshClaims.FindID(p.DB)
	if err != nil {
		log.Println("Error while trying to find refresh token ID: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	// Insert or update token identifier
	if exists {
		refreshClaims.TokenID = refreshID
		_, err = refreshClaims.Update(p.DB, uint64(refreshClaims.ExpiresAt.Unix()))
		if err != nil {
			log.Println("Error updating refresh token: ", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}
	} else {
		id, err := refreshClaims.InsertToDb(p.DB)
		if err != nil {
			log.Println("Error inserting refresh token to the database: ", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}
		refreshClaims.TokenID = id
	}

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

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"expires_at":   expiresAt.Unix(),
		"role":         ROLE,
	})

}
