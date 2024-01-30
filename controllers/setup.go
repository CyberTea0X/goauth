package controllers

import (
	"database/sql"

	"github.com/CyberTea0X/goauth/src/backend/models"
	"github.com/gin-gonic/gin"
)

type PublicController struct {
	DB              *sql.DB
	AccessTokenCfg  models.AccessTokenCfg
	RefreshTokenCfg models.RefreshTokenCfg
}

func SetupController(tokensConfig models.TokensCfg, db *sql.DB) *PublicController {

	pCtrl := new(PublicController)

	pCtrl.AccessTokenCfg = tokensConfig.Access
	pCtrl.RefreshTokenCfg = tokensConfig.Refresh

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
