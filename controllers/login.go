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

type LoginInput struct {
	Username string `form:"username"`
	Password string `form:"password" binding:"required"`
	Email    string `form:"email"`
	DeviceId uint   `form:"device_id" binding:"required"`
}

type LoginOutput struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresAt    int64    `json:"expires_at"`
	Roles        []string `json:"role"`
}

// Function that is responsible for user authorization.
//
// In response to a successful authorization request, returns
// access token and refresh token, as well as the time of death of the access token
func (p *PublicController) Login(c *gin.Context) {

	var input LoginInput

	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrToMap(models.ErrInvalidQuery))
		return
	}
	user, err := models.LoginUser(p.Client, p.LoginServiceURL, input.Username, input.Password, input.Email)

	if err != nil {
		externalErr := new(models.ExternalServiceError)
		serviceErr := new(models.ConstantError)
		if errors.As(err, &externalErr) {
			c.JSON(externalErr.Status, models.ErrToMap(externalErr))
		} else if errors.As(err, &serviceErr) {
			c.JSON(http.StatusBadRequest, models.ErrToMap(err))
		} else {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	expiresAt := time.Now().Add(p.RefreshTokenCfg.Lifespan())
	refreshClaims := token.NewRefresh(-1, input.DeviceId, user.Id, user.Roles, expiresAt)

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
	accessClaims := token.NewAccess(user.Id, user.Roles, expiresAt)
	accessToken, err := accessClaims.TokenString(p.AccessTokenCfg.Secret)

	if err != nil {
		log.Println("Error generating access token: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt.Unix(),
		Roles:        user.Roles,
	})

}
