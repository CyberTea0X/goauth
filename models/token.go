package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type RefreshToken struct {
	gorm.Model
    DeviceID uint `gorm:"not null;" json:"device_id"`
	Token string `gorm:"size:255;not null;" json:"token"`
    UserID uint
}

func (t *RefreshToken) SaveToken(db *gorm.DB) (*RefreshToken, error) {
    t.CreatedAt = time.Now()
    t.UpdatedAt = time.Now()
	return t, db.Create(&t).Error
}

func DeleteOldTokens(db *gorm.DB, user_id uint, device_id uint) {
	var t RefreshToken
    db.Where("user_id = ? AND device_id = ?", user_id, device_id).Unscoped().Delete(&t)
}

func RefreshTokenExists(db *gorm.DB, user_id uint, device_id uint) bool {
	var t RefreshToken
    db.Where("user_id = ? AND device_id = ?", user_id, device_id).First(&t)
    return db.Error == nil
}
