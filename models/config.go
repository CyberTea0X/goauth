package models

import (
	"errors"
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type AccessTokenCfg struct {
	Secret         string
	LifespanMinute uint `toml:"lifespan_minute"`
}

type RefreshTokenCfg struct {
	Secret       string
	LifespanHour uint `toml:"lifespan_hour"`
}

type TokensCfg struct {
	Access  AccessTokenCfg
	Refresh RefreshTokenCfg
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

type GuestService struct {
	Host string
	Port string
	Path string
}

type LoginService struct {
	Host string
	Port string
	Path string
}

type ExternalServicesConfig struct {
	Guest GuestService
	Login LoginService
}

type TomlConfig struct {
	Database DatabaseConfig
	Tokens   TokensCfg
	Services ExternalServicesConfig
}

func ParseConfig(filename string) (*TomlConfig, error) {
	configFile, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	tomlConfig := new(TomlConfig)

	err = toml.Unmarshal(configFile, tomlConfig)

	if err != nil {
		return nil, errors.Join(errors.New("Failed to parse config file"), err)
	}

	return tomlConfig, nil
}
