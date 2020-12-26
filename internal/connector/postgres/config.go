package postgres

import (
	"errors"
	"fmt"
)

// Config validation errors.
var (
	ErrEmptyAddr     = errors.New("addr is empty")
	ErrEmptyDatabase = errors.New("database is empty")
	ErrEmptyUser     = errors.New("user is empty")
	ErrEmptyPassword = errors.New("password is empty")
)

// Config contains credentials for PostgreSQL database.
type Config struct {
	Addr         string            `yaml:"addr" json:"addr"`
	Database     string            `yaml:"database" json:"database"`
	User         string            `yaml:"username" json:"user"`
	Password     string            `yaml:"password" json:"password"`
	Notify       map[string]string `yaml:"notify" json:"notify"`
	Debug        bool              `yaml:"debug" json:"debug"`
	PoolSize     int               `yaml:"pool_size" json:"pool_size"`
	MaxIdleConns int               `yaml:"max_idle_conns" json:"max_idle_conns"`
	MaxOpenConns int               `yaml:"max_open_conns" json:"max_open_conns"`
}

// Validate checks required fields and validates for allowed values.
func (cfg *Config) Validate() error {
	if cfg.Addr == "" {
		return ErrEmptyAddr
	}

	if cfg.Database == "" {
		return ErrEmptyDatabase
	}

	if cfg.User == "" {
		return ErrEmptyUser
	}

	if cfg.Password == "" {
		return ErrEmptyPassword
	}

	return nil
}

// GetDataSourceName returns Data Source Name connection string to PostgreSQL database.
func (cfg *Config) GetDataSourceName() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Addr, cfg.Database)
}
