package process

import (
	"errors"
	"time"
)

// Validation errors.
var (
	ErrEmptyStartInterval = errors.New("start_interval is empty")
)

// Config is config for process.
type Config struct {
	Disabled      bool          `yaml:"disabled" json:"disabled"`
	StartInterval time.Duration `yaml:"start_interval" json:"start_interval"`
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

	return nil
}
