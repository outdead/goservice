package elasticsearch

import (
	"errors"
	"time"
)

const DefaultHealthcheckInterval = 5 * time.Second

// Config validation errors.
var (
	ErrEmptyAddr           = errors.New("addr is empty")
	ErrEmptyDatabase       = errors.New("database is empty")
	ErrHealthcheckInterval = errors.New("healthcheck_interval must be positive number or zero")
)

// Config contains credentials for Elasticsearch database.
type Config struct {
	Addr                string        `yaml:"addr" json:"addr"`
	Database            string        `yaml:"database" json:"database"`
	HealthcheckInterval time.Duration `json:"healthcheck_interval"`
}

// Validate checks required fields and validates for allowed values.
func (cfg Config) Validate() error {
	if cfg.Addr == "" {
		return ErrEmptyAddr
	}

	if cfg.Database == "" {
		return ErrEmptyDatabase
	}

	if cfg.HealthcheckInterval < 0 {
		return ErrHealthcheckInterval
	}

	return nil
}
