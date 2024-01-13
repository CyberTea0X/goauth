package models

import (
	"fmt"
)

type DatabaseConfig struct {
	Dbdriver   string `env:"DB_DRIVER,required"`
	DbHost     string `env:"DB_HOST,required"`
	DbUser     string `env:"DB_USER,required"`
	DbPassword string `env:"DB_PASSWORD,required"`
	DbName     string `env:"DB_NAME,required"`
	DbPort     string `env:"DB_PORT,required"`
}

func (c *DatabaseConfig) GetUrl() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", c.DbUser, c.DbPassword, c.DbHost, c.DbPort, c.DbName)
}

type AccessTokenCfg struct {
	Secret   string `env:"ACCESS_TOKEN_SECRET,required"`
	Lifespan uint   `env:"ACCESS_TOKEN_MINUTE_LIFESPAN,required"`
}

type RefreshTokenCfg struct {
	Secret   string `env:"REFRESH_TOKEN_SECRET,required"`
	Lifespan uint   `env:"REFRESH_TOKEN_HOUR_LIFESPAN,required"`
}
