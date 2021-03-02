package redis

import (
	"errors"
	"time"
)

// Config validation errors.
var (
	ErrEmptyAddr = errors.New("addr is empty")
)

// Config contains credentials for Redis database.
type Config struct {
	Addr         string        `yaml:"addr"`
	Password     string        `yaml:"password"`
	DB           int           `yaml:"db"`
	TTL          time.Duration `yaml:"ttl"`
	MaxRetries   int           `yaml:"max_retries"`
	DialTimeout  time.Duration `yaml:"dial_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	PoolSize     int           `yaml:"pool_size"`
}

// Validate checks required fields and validates for allowed values.
func (cfg *Config) Validate() error {
	if cfg.Addr == "" {
		return ErrEmptyAddr
	}

	return nil
}
