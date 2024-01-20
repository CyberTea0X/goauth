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
	TokenID  int64  `json:"token_id"`
	DeviceID uint   `json:"device_id"`
	UserID   int64  `json:"user_id"`
	Role     string `json:"role"`
}

func NewRefresh(tokenId int64, deviceId uint, userId int64, role string, expiresAt time.Time) *RefreshToken {
	t := new(RefreshToken)
    t.TokenID = tokenId
	t.DeviceID = deviceId
	t.UserID = userId
	t.ExpiresAt = jwt.NewNumericDate(expiresAt)
	t.Role = role
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
	const query = "INSERT INTO refresh_tokens (device_id, expires_at, user_id) VALUES (?,?,?)"
	res, err := db.Exec(query, t.DeviceID, t.ExpiresAt.Unix(), t.UserID)
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
		"WHERE id=? AND user_id=? AND device_id=? AND expires_at=?)"
	res, err := db.Query(query, t.TokenID, t.UserID, t.DeviceID, t.ExpiresAt.Unix())
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
		"WHERE user_id=? AND device_id=?"
	res, err := db.Query(query, t.UserID, t.DeviceID)
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
	const q = "UPDATE refresh_tokens SET expires_at=? WHERE id=? AND user_id =? AND device_id =?"
	_, err := db.Exec(q, expiresAt, t.TokenID, t.UserID, t.DeviceID)
	return t, err
}

// Creates refresh token table if it doesn't already exist
func CreateRefreshTable(db *sql.DB) error {
	_, err := db.Exec(refresh_table_ddl)
	return err
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
