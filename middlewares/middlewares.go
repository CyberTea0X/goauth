package middlewares

import (
	"net/http"
	"github.com/gin-gonic/gin"
    "github.com/CyberTea0X/delta_art/src/backend/utils/token"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := token.AccessParse(token.ExtractToken(c))
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
        c.Set("access_token", token)
		c.Next()
	}
}
