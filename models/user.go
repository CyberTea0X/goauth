package models

import (
	"errors"
	"fmt"
	"html"
	"net/mail"
	"strings"
	"time"

	"github.com/CyberTea0X/delta_art/src/backend/utils/token"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:255;not null;" json:"username"`
	Password string `gorm:"size:255;not null;" json:"password"`
    Email string `gorm:"size:255;not null;unique;" json:"email"`
    RefreshTokens []RefreshToken `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
    u.CreatedAt = time.Now()
    u.UpdatedAt = time.Now()
	return u, db.Create(&u).Error
}

func (u *User) BeforeSave() error {

	//turn password into hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password),bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	//remove spaces in username 
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	return nil

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


func GetUserByID(db *gorm.DB, uid uint) (User,error) {

	var u User

	if err := db.First(&u,uid).Error; err != nil {
        return u, errors.New(fmt.Sprintf("User with id=%d not found", uid))
	}

	u.PrepareGive()
	
	return u,nil

}

func (u *User) PrepareGive(){
	u.Password = ""
    u.Email = ""
}

type TokensData struct {
    AccessToken string
    RefreshToken string
    AccessTokenExpires int64
}


func Login(db *gorm.DB, u User, device_id uint) (TokensData, error) {
	
	var err error

    var udb User

    if u.Username != "" {
	    err = db.Model(User{}).Where("username = ?", u.Username).Take(&udb).Error
    } else {
        err = db.Model(User{}).Where("email = ?", u.Email).Take(&udb).Error
    }

	if err != nil {
		return TokensData{}, err
	}

    err = VerifyPassword(u.Password, udb.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return TokensData{}, err
	}

    DeleteOldTokens(db, udb.ID, device_id)

	access_token, expires, err := token.GenerateAccessToken(udb.ID)

	if err != nil {
		return TokensData{}, err
	}

    refresh_token, err := token.GenerateRefresh(udb.ID, device_id)

    if err != nil {
		return TokensData{}, err
	}

    refresh_model := RefreshToken{}
    refresh_model.Token = refresh_token
    refresh_model.UserID = udb.ID
    refresh_model.DeviceID = device_id

    if _, err := refresh_model.SaveToken(db); err != nil {
        return TokensData{}, err
    }



	return TokensData{access_token, refresh_token, expires}, err
	
}
