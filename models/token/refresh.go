package token

import (
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const refresh_table_ddl = "" +
	"CREATE TABLE IF NOT EXISTS `refresh_tokens` (" +
	"`id` int(10) unsigned NOT NULL AUTO_INCREMENT," +
	"`device_id` int(10) unsigned NOT NULL," +
	"`user_id` int(10) unsigned NOT NULL," +
	"`expires_at` int(10) unsigned NOT NULL," +
	"PRIMARY KEY (`id`)," +
	"UNIQUE KEY `refresh_tokens_device_id_IDX` (`device_id`,`user_id`,`expires_at`) USING BTREE" +
	") ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;"

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
	var exists bool
	query := "SELECT EXISTS(SELECT * FROM refresh_tokens WHERE user_id=? AND device_id=? AND expires_at=?)"
	res, err := db.Query(query, t.UserID, t.DeviceID, t.ExpiresAt.Unix())
	res.Next()
	res.Scan(&exists)
	return exists, err
}

// Deletes old refresh token from the database (not exactly token, but it's identifier)
func DeleteOldToken(db *sql.DB, user_id uint, device_id uint) error {
	query := "DELETE FROM refresh_tokens WHERE user_id =? AND device_id =?"
	_, err := db.Exec(query, user_id, device_id)
	return err
}

// Creates refresh token table if it doesn't already exist
func CreateRefreshTable(db *sql.DB) error {
	_, err := db.Exec(refresh_table_ddl)
	return err
}
