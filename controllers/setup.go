package controllers

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/CyberTea0X/goauth/src/backend/models"
	"github.com/gin-gonic/gin"
)

type PublicController struct {
	DB              *sql.DB
	AccessTokenCfg  models.AccessTokenCfg
	RefreshTokenCfg models.RefreshTokenCfg
	GuestServiceURL url.URL
	LoginServiceURL url.URL
	Client          models.HTTPClient
}

func NewController(tokensConfig models.TokensCfg, servicesConfig models.ExternalServicesConfig, client models.HTTPClient, db *sql.DB) *PublicController {

	pCtrl := new(PublicController)

	pCtrl.AccessTokenCfg = tokensConfig.Access
	pCtrl.RefreshTokenCfg = tokensConfig.Refresh
	g := servicesConfig.Guest
	pCtrl.GuestServiceURL = url.URL{
		Host:   fmt.Sprintf("%s:%s", g.Host, g.Port),
		Path:   g.Path,
		Scheme: "http",
	}
	l := servicesConfig.Guest
	pCtrl.LoginServiceURL = url.URL{
		Host:   fmt.Sprintf("%s:%s", l.Host, l.Port),
		Path:   l.Path,
		Scheme: "http",
	}
	pCtrl.Client = client

	pCtrl.DB = db
	return pCtrl
}

// We need separate function for router setup to do testing properly
func SetupRouter(c *PublicController) *gin.Engine {
	router := gin.Default()

	public := router.Group("api")
	public.GET("health_check", c.HealthCheck)
	public.GET("login", c.Login)
	public.GET("refresh", c.Refresh)
	public.GET("auth", c.Auth)
	public.GET("guest", c.Guest)
	return router
}
