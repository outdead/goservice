package process

import (
	"errors"
	"time"
)

// Validation errors.
var (
	ErrEmptyStartInterval = errors.New("start_interval is empty")
	ErrRecordsInterval    = errors.New("records_interval is empty")
)

// Config is config for process.
type Config struct {
	Disabled        bool          `yaml:"disabled" json:"disabled"`
	StartInterval   time.Duration `yaml:"start_interval" json:"start_interval"`
	RecordsInterval time.Duration `yaml:"records_interval" json:"records_interval"`
	TimeOffset      time.Duration `yaml:"time_offset" json:"time_offset"`
}

// Validate checks config to required fields.
func (cfg *Config) Validate() error {
	if cfg.Disabled {
		// Do not validate disabled component.
		return nil
	}

	if cfg.StartInterval == 0 {
		return ErrEmptyStartInterval
	}

	if cfg.RecordsInterval == 0 {
		return ErrRecordsInterval
	}

	return nil
}
