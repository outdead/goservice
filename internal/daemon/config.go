package daemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"time"

	"github.com/outdead/echo-skeleton/internal/logger"
	"gopkg.in/yaml.v3"
)

var (
	// ErrInvalidConfig is a basic configuration validation error. It is wrapped
	// if config validation fails.
	ErrInvalidConfig = errors.New("config validation error")

	// ErrParseConfig is returned when parsing the config from the file fails.
	ErrParseConfig = errors.New("config parse error")

	// ErrInvalidConfigExtension is returned when parsing a config from a file
	// when the file has an unsupported extension.
	ErrInvalidConfigExtension = errors.New("invalid config extension")
)

// Config is main service config structure.
type Config struct {
	App struct {
		Port                     string        `json:"port" yaml:"port"`
		ProfilerPort             string        `json:"profiler_port" yaml:"profiler_port"`
		CheckConnectionsInterval time.Duration `json:"check_connections_interval" yaml:"check_connections_interval"`
		ErrorBuffer              int           `json:"error_buffer" yaml:"error_buffer"`
		Log                      logger.Config `json:"log" yaml:"log"`
	} `json:"app" yaml:"app"`
}

// NewConfig creates new config from `name` file data.
func NewConfig(name string) (*Config, error) {
	cfg := new(Config)
	if err := cfg.ParseFromFile(name); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParseConfig, err)
	}

	return cfg, nil
}

// ParseFromFile reads config text data and binds to config struct.
func (cfg *Config) ParseFromFile(name string) error {
	file, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}

	switch ext := path.Ext(name); ext {
	case ".yaml":
		err = yaml.Unmarshal(file, cfg)
	case ".json":
		err = json.Unmarshal(file, cfg)
	default:
		err = fmt.Errorf("%w: %s", ErrInvalidConfigExtension, ext)
	}

	return err
}

// Validate checks config to required fields.
func (cfg *Config) Validate() (err error) {
	if cfg.App.Port == "" {
		return fmt.Errorf("%w: app.port is empty", ErrInvalidConfig)
	}

	if cfg.App.CheckConnectionsInterval == 0 {
		return fmt.Errorf("%w: app.check_connections_interval is empty", ErrInvalidConfig)
	}

	if cfg.App.ErrorBuffer == 0 {
		return fmt.Errorf("%w: app.error_buffer is empty", ErrInvalidConfig)
	}

	return
}

// Print print config to console.
func (cfg *Config) Print() error {
	js, err := json.MarshalIndent(cfg, "", "  ")
	if err == nil {
		fmt.Println(string(js))
	}

	return err
}
