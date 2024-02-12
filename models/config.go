package models

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pelletier/go-toml/v2"
)

type AppConfig struct {
	Port string `validate:"required"`
	// http.Client timeout
	TimeoutSeconds uint `toml:"timeout_seconds" validate:"required"`
}

type LifespanCfg struct {
	LifespanMinute uint `toml:"lifespan_minute"`
	LifespanHour   uint `toml:"lifespan_hour"`
	LifespanDay    uint `toml:"lifespan_day"`
	lifespan       *time.Duration
}

// Just calculates time.Now() + lifespan
func (c *LifespanCfg) Lifespan() time.Duration {
	if c.lifespan != nil {
		return *c.lifespan
	}
	lifespan := new(time.Duration)
	*lifespan = time.Minute * time.Duration(c.LifespanMinute)
	*lifespan += time.Hour * time.Duration(c.LifespanHour)
	*lifespan += time.Hour * time.Duration(c.LifespanDay) * 24
	c.lifespan = lifespan
	return *lifespan
}

type AccessTokenCfg struct {
	LifespanCfg
	Secret string `validate:"required"`
}

type RefreshTokenCfg struct {
	LifespanCfg
	Secret string `validate:"required"`
}

type TokensCfg struct {
	Access  AccessTokenCfg  `validate:"required"`
	Refresh RefreshTokenCfg `validate:"required"`
}

type DatabaseConfig struct {
	Driver   string `validate:"required"`
	Host     string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	Name     string `toml:"database" validate:"required"`
	Port     string `validate:"required"`
}

func (c *DatabaseConfig) GetUrl() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.Name)
}

type GuestService struct {
	Host string `validate:"required"`
	Port string `validate:"required"`
	Path string `validate:"required"`
}

type LoginService struct {
	Host string `validate:"required"`
	Port string `validate:"required"`
	Path string `validate:"required"`
}

type ExternalServicesConfig struct {
	Guest GuestService `validate:"required"`
	Login LoginService `validate:"required"`
}

type TomlConfig struct {
	App      AppConfig              `validate:"required"`
	Database DatabaseConfig         `validate:"required"`
	Tokens   TokensCfg              `validate:"required"`
	Services ExternalServicesConfig `validate:"required"`
}

func (*TomlConfig) Validate() error {
	return nil
}

// Parses and validates config
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
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(tomlConfig)
	if err != nil {
		return nil, errors.Join(errors.New("Failed to parse config file"), err)
	}

	return tomlConfig, nil
}
