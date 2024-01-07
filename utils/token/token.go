package token

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Function that generates a JWT token that contains data:
// Authorized user ID;
// JWT token expiration date;
func GenerateAccessToken(user_id uint) (string, int64, error) {

	token_lifespan,err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_MINUTE_LIFESPAN"))
	if err != nil {
		return "", 0, err
	}
    expires := time.Now().Add(time.Minute * time.Duration(token_lifespan)).Unix()


	claims := jwt.MapClaims{
        "user_id": user_id,
        "exp": expires,
    }
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    token_string, err := token.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))

	return token_string, expires, err

}

// Function that generates a Refresh JWT token that contains data:
// User ID;
// Refresh JWT token expiration date;
func GenerateRefresh(user_id uint, device_id uint) (string, error) {

	token_lifespan,err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_HOUR_LIFESPAN"))
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
        "user_id": user_id,
        "device_id": device_id,
        "exp": time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    token_string, err := token.SignedString([]byte(os.Getenv("REFRESH_TOKEN_SECRET")))

	return token_string, err

}

// Parses token from token string
// While parsing checks if token is valid by checking its sign and signature
func Parse(token string, signature string) (*jwt.Token, error) {
	jwt, result := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signature), nil
	})
	return jwt, result
}

// Parses refresh token from token string
// While parsing checks if token is valid by checking its sign and signature
func RefreshParse(token string) (*jwt.Token, error) {
    return Parse(token, os.Getenv("REFRESH_TOKEN_SECRET"))
}

// Parses access token from token string
// While parsing checks if token is valid by checking its sign and signature
func AccessParse(token string) (*jwt.Token, error) {
    return Parse(token, os.Getenv("ACCESS_TOKEN_SECRET"))
}

// Extracts Token from query
func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

// Extracts uint value by key from token
// Error if token is invalid or parsing is failed
func ExtractUint(jwt_token *jwt.Token, key string) (uint, error) {
	claims, ok := jwt_token.Claims.(jwt.MapClaims)
	if ok && jwt_token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims[key]), 10, 64)
		return uint(uid), err
	}
	return 0, nil
}
