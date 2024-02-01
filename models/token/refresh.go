package token

import (
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RefreshToken struct {
	jwt.RegisteredClaims
	TokenID  int64    `json:"token_id"`
	DeviceID uint     `json:"device_id"`
	UserID   int64    `json:"user_id"`
	Roles    []string `json:"role"`
}

func NewRefresh(tokenId int64, deviceId uint, userId int64, roles []string, expiresAt time.Time) *RefreshToken {
	t := new(RefreshToken)
	t.TokenID = tokenId
	t.DeviceID = deviceId
	t.UserID = userId
	t.ExpiresAt = jwt.NewNumericDate(expiresAt)
	t.Roles = roles
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

// Inserts row that identifies token into the database (not token string)
func (t *RefreshToken) InsertToDb(db *sql.DB) (int64, error) {
	const query = "INSERT INTO refresh_tokens (device_id, expires_at, user_id, role) VALUES (?,?,?,?)"
	res, err := db.Exec(query, t.DeviceID, t.ExpiresAt.Unix(), t.UserID, t.Roles)
	if err != nil {
		return 0, err
	}
	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastInsertId, err
}

func (t *RefreshToken) Exists(db *sql.DB) (bool, error) {
	var exists bool
	const query = "" +
		"SELECT EXISTS(SELECT * FROM refresh_tokens " +
		"WHERE id=? AND user_id=? AND device_id=? AND expires_at=? AND role=?)"
	res, err := db.Query(query, t.TokenID, t.UserID, t.DeviceID, t.ExpiresAt.Unix(), t.Roles)
	if err != nil {
		return false, err
	}
	if !res.Next() {
		return false, nil
	}
	res.Scan(&exists)
	return exists, err
}

func (t *RefreshToken) FindID(db *sql.DB) (int64, bool, error) {
	const query = "" +
		"SELECT id FROM refresh_tokens " +
		"WHERE user_id=? AND device_id=? AND role=?"
	res, err := db.Query(query, t.UserID, t.DeviceID, t.Roles)
	if err != nil {
		return 0, false, err
	}
	if !res.Next() {
		return 0, false, nil
	}
	var id int64
	res.Scan(&id)
	return id, true, nil
}

// Updates expiredAt field of token in the database
func (t *RefreshToken) Update(db *sql.DB, expiresAt uint64) (*RefreshToken, error) {
	const q = "UPDATE refresh_tokens SET expires_at=? WHERE id=? AND user_id =? AND device_id =? AND role=?"
	_, err := db.Exec(q, expiresAt, t.TokenID, t.UserID, t.DeviceID, t.Roles)
	return t, err
}

func (t *RefreshToken) InsertOrUpdate(db *sql.DB) (int64, error) {
	id, exists, err := t.FindID(db)
	if err != nil {
		return 0, errors.Join(errors.New("error trying to find refresh token ID"), err)
	}
	if !exists {
		if id, err = t.InsertToDb(db); err != nil {
			return 0, errors.Join(errors.New("error trying to insert refresh token to DB"), err)
		}
		return id, nil
	}
	old_id := t.TokenID
	t.TokenID = id
	if _, err = t.Update(db, uint64(t.ExpiresAt.Unix())); err != nil {
		return 0, errors.Join(errors.New("error trying to update refresh token in the DB"), err)
	}
	t.TokenID = old_id
	return id, nil
}
