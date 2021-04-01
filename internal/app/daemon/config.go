package daemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"time"

	"github.com/outdead/goservice/internal/connector"
	"github.com/outdead/goservice/internal/utils/logutil"
	"gopkg.in/yaml.v3"
)

// Validation errors.
var (
	ErrEmptyPort                     = errors.New("app.port is empty")
	ErrEmptyCheckConnectionsInterval = errors.New("app.check_connections_interval is empty")
	ErrEmptyErrorBuffer              = errors.New("app.error_buffer is empty")

	// ErrInvalidConfigExtension is returned when parsing a config from a file
	// when the file has an unsupported extension.
	ErrInvalidConfigExtension = errors.New("invalid config extension")
)

// Config is main service config structure.
type Config struct {
	App struct {
		Port                     string         `json:"port" yaml:"port"`
		ProfilerAddr             string         `json:"profiler_addr" yaml:"profiler_addr"`
		CheckConnectionsInterval time.Duration  `json:"check_connections_interval" yaml:"check_connections_interval"`
		ErrorBuffer              int            `json:"error_buffer" yaml:"error_buffer"`
		Log                      logutil.Config `json:"log" yaml:"log"`
	} `json:"app" yaml:"app"`
	Connections connector.Config `yaml:"connections" json:"connections"`
}

// NewConfig creates new config from `name` file data.
func NewConfig(name string) (*Config, error) {
	cfg := new(Config)
	if err := cfg.ParseFromFile(name); err != nil {
		return nil, err
	}

	return cfg, nil
}

// ParseFromFile reads config text data and binds to config struct.
func (cfg *Config) ParseFromFile(name string) error {
	file, err := ioutil.ReadFile(name)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
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
func (cfg *Config) Validate() error {
	if cfg.App.Port == "" {
		return ErrEmptyPort
	}

	if cfg.App.CheckConnectionsInterval == 0 {
		return ErrEmptyCheckConnectionsInterval
	}

	if cfg.App.ErrorBuffer == 0 {
		return ErrEmptyErrorBuffer
	}

	if err := cfg.Connections.Validate(); err != nil {
		return fmt.Errorf("connections: %w", err)
	}

	return nil
}

// Print print config to console.
func (cfg *Config) Print() error {
	js, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("print config: %w", err)
	}

	fmt.Println(string(js))

	return nil
}
