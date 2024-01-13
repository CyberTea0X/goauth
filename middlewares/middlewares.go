package middlewares

import (
	"github.com/CyberTea0X/goauth/src/backend/token"
	"github.com/gin-gonic/gin"
	"net/http"
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
