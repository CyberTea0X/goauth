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
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if !jwt.Valid {
		return nil, result
	}
	return jwt, result
}

// Extracts Token from query or "Authorization" request header
func ExtractToken(c *gin.Context) (string, error) {
	auth := c.Request.Header.Get("Authorization")
	// Authorization usually starts on "Bearer", then comes single whitespace and then token string
	authSplit := strings.Split(auth, " ")
	fmt.Println(authSplit)
	if (len(authSplit) == 2) && (authSplit[1] != "") {
		return authSplit[1], nil
	} else if auth != "" {
		return auth, nil
	}
	token := c.Query("token")
	if token == "" {
		return "", errors.New("token not found in request or request headers")
	}
	return token, nil
}
