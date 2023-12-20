package config

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
)

type AppConfig struct {
	API      *APIConfig      `mapstructure:"API"`
	Postgres *PostgresConfig `mapstructure:"POSTGRES"`
}

func (c *AppConfig) validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.API, validation.Required),
		validation.Field(&c.Postgres, validation.Required),
	)
}

func (c *AppConfig) validateConfig() error {
	if err := c.validate(); err != nil {
		return fmt.Errorf("c.validate() -> %w", err)
	}

	if err := c.API.validate(); err != nil {
		return fmt.Errorf("c.API.validate() -> %w", err)
	}

	if err := c.Postgres.validate(); err != nil {
		return fmt.Errorf("c.Postgres.validate() -> %w", err)
	}

	return nil
}

type APIConfig struct {
	Environment string `mapstructure:"ENV"`
	Host        string `mapstructure:"HOST"`
	Port        string `mapstructure:"PORT"`
}

func (c *APIConfig) validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.Environment, validation.Required),
		validation.Field(&c.Host, validation.Required),
		validation.Field(&c.Port, validation.Required),
	)
}

type PostgresConfig struct {
	Host     string `mapstructure:"HOST"`
	Port     string `mapstructure:"PORT"`
	User     string `mapstructure:"USER"`
	Password string `mapstructure:"PASSWORD"`
	DB       string `mapstructure:"DB"`
}

func (c *PostgresConfig) validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.Host, validation.Required),
		validation.Field(&c.Port, validation.Required),
		validation.Field(&c.User, validation.Required),
		validation.Field(&c.Password, validation.Required),
		validation.Field(&c.DB, validation.Required),
	)
}
