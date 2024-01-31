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

func (*TomlConfig) Validate() error {
	return nil
}

func ParseConfig(filename string) (*TomlConfig, error) {
	configFile, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	tomlConfig := new(TomlConfig)
	d := toml.NewDecoder(configFile)
	d.DisallowUnknownFields()
	if err = d.Decode(tomlConfig); err != nil {
		var details *toml.StrictMissingError
		if !errors.As(err, &details) {
			return nil, errors.Join(errors.New("Failed to parse config file"), err)
		}
		return nil, errors.Join(fmt.Errorf("Failed to parse config file\n%s", details.String()))
	}
	fmt.Println(tomlConfig.Services)

	return tomlConfig, nil
}
