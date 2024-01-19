package models

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
)

type User struct {
	Id       uint   `json:"id"`
	Username string `json:"username"` // not null
	Password string `json:"password"` // not null
	Email    string `json:"email"`    // not null unique
	Role     string `json:"role"`     // not null
}

func VerifyPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err == nil {
		return true
	}
	return false
}

// Calls rows.Next() and Scans the row into the user struct
func (u *User) FromRow(rows *sql.Rows) (*User, error) {
	exists := rows.Next()
	if !exists {
		return u, errors.New(fmt.Sprintf("Can't scan user from row"))
	}

	err := rows.Scan(&u.Id, &u.Username, &u.Email, &u.Password, &u.Role)
	if err != nil {
		return u, errors.New(fmt.Sprintf("Can't scan user from row: %s", err.Error()))
	}

	return u, nil
}

func GetUserByEmail(db *sql.DB, email string) (*User, error) {

	u := &User{}

	rows, err := db.Query("SELECT * FROM users WHERE email =? LIMIT 1", email)
	if err != nil {
		return u, errors.New(fmt.Sprintf("User with email=%s not found: %s", email, err.Error()))
	}
	exists := rows.Next()
	if !exists {
		return u, errors.New(fmt.Sprintf("User with email=%s not found", email))
	}
	return u.FromRow(rows)

}

func GetUserByUsername(db *sql.DB, username string) (*User, error) {

	u := &User{}

	row, err := db.Query("SELECT * FROM users WHERE username =? LIMIT 1", username)
	if err != nil {
		return u, errors.New(fmt.Sprintf("User with username=%s not found: %s", username, err.Error()))
	}
	return u.FromRow(row)
}

func (u *User) PrepareGive() {
	u.Password = ""
	u.Email = ""
}
