package models

import (
	"fmt"
)

type TomlConfig struct {
	Database  DatabaseConfig
	TokensCfg TokensCfg `toml:"tokens"`
}

type TokensCfg struct {
	AccessTokenCfg  AccessTokenCfg  `toml:"access"`
	RefreshTokenCfg RefreshTokenCfg `toml:"refresh"`
}

type AccessTokenCfg struct {
	Secret         string
	LifespanMinute uint `toml:"lifespan_minute"`
}

type RefreshTokenCfg struct {
	Secret       string
	LifespanHour uint `toml:"lifespan_hour"`
}

type DatabaseConfig struct {
	Driver   string
	Host     string
	User     string
	Password string
	Name     string `toml:"database"`
	Port     string
}

func (c *DatabaseConfig) GetUrl() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.Name)
}
