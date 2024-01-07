package models

import (
	"errors"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

func ConnectDataBase() (*gorm.DB, error) {

	err := godotenv.Load(".env")

	if err != nil {
        return &gorm.DB{}, errors.New("error opening .env file: " + err.Error())
	}	

	Dbdriver := os.Getenv("DB_DRIVER")
	DbHost := os.Getenv("DB_HOST")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DbPort := os.Getenv("DB_PORT")

	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
	
    DB, err := gorm.Open(Dbdriver, DBURL)

	if err != nil {
        return &gorm.DB{}, err
	} 

	DB.AutoMigrate(&User{})
    DB.AutoMigrate(&RefreshToken{})
    return DB, err
}
