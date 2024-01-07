package controllers

import (
	"log"

	"github.com/CyberTea0X/delta_art/src/backend/models"
	"github.com/jinzhu/gorm"
)

type PublicController struct {
    DB *gorm.DB
}

func Setup() *PublicController {
    pCtrl := new(PublicController)
    db, err := models.ConnectDataBase()

    if err != nil {
        log.Fatal("failed to connect to database:", err)
    }

    pCtrl.DB = db
    return pCtrl
}
