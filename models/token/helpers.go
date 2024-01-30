package token

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// While parsing checks if token is valid by checking its sign method and signature
// signature is token specific secret key
func ParseWithClaims(token string, claims jwt.Claims, secret string) (*jwt.Token, error) {
	jwt, result := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if !jwt.Valid {
		return nil, result
	}
	return jwt, result
}

// Extracts Token from query or request header
func ExtractToken(c *gin.Context) (string, error) {
	token := c.Query("token")
	if token != "" {
		return token, nil
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1], nil
	}
	if bearerToken == "" {
		return "", errors.New("Token not found in request or request headers")
	}
	return bearerToken, nil
}
