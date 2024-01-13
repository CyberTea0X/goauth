package token

import (
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RefreshToken struct {
	jwt.RegisteredClaims
	DeviceID uint `json:"device_id"`
	UserID   uint `json:"user_id"`
}

func NewRefresh(deviceId uint, userId uint, expiresAt time.Time) *RefreshToken {
	t := new(RefreshToken)
	t.DeviceID = deviceId
	t.UserID = userId
	t.ExpiresAt = jwt.NewNumericDate(expiresAt)
	return t
}

// Parses token from token string
func RefreshFromString(token string, secret string) (*RefreshToken, error) {
	t, err := ParseWithClaims(token, &RefreshToken{}, secret)
	if err != nil {
		return nil, err
	}
	if claims, ok := t.Claims.(*RefreshToken); ok {
		return claims, nil
	}
	return nil, errors.New("unknown claims type, cannot proceed")
}

// Encodes RefreshToken model into a JWT string
func (c *RefreshToken) TokenString(secret string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	token_string, err := token.SignedString([]byte(secret))

	return token_string, err
}

func (t *RefreshToken) InsertToDb(db *sql.DB) (*RefreshToken, error) {
	query := "INSERT INTO refresh_tokens (device_id, expires_at, user_id) VALUES (?,?,?)"
	_, err := db.Exec(query, t.DeviceID, t.ExpiresAt.Unix(), t.UserID)
	return t, err
}

func (t *RefreshToken) Exists(db *sql.DB) (bool, error) {
	query := "SELECT EXISTS(SELECT * FROM refresh_tokens WHERE user_id=? AND device_id=? AND expires_at=?)"
	res, err := db.Query(query, t.UserID, t.DeviceID, t.ExpiresAt.Unix())
	return res.Next(), err
}

// Deletes old refresh token from the database (not exactly token, but it's identifier)
func DeleteOldToken(db *sql.DB, user_id uint, device_id uint) error {
	query := "DELETE FROM refresh_tokens WHERE user_id =? AND device_id =?"
	_, err := db.Exec(query, user_id, device_id)
	return err
}
