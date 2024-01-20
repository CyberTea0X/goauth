package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessToken struct {
	jwt.RegisteredClaims
	UserId uint   `json:"user_id"`
	Role   string `json:"role"`
}

func NewAccess(userId uint, role string, expiresAt time.Time) *AccessToken {
	t := new(AccessToken)
	t.UserId = userId
	t.ExpiresAt = jwt.NewNumericDate(expiresAt)
	t.Role = role
	return t
}

// Parses token from token string
func AccessFromString(token string, secret string) (*AccessToken, error) {
	t, err := ParseWithClaims(token, &AccessToken{}, secret)
	if err != nil {
		return nil, err
	}
	if claims, ok := t.Claims.(*AccessToken); ok {
		return claims, nil
	}
	return nil, errors.New("unknown claims type, cannot proceed")
}

// Encodes RefreshToken model into a JWT string
func (c *AccessToken) TokenString(secret string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	token_string, err := token.SignedString([]byte(secret))

	return token_string, err
}
