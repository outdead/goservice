package clickhouse

import (
	"errors"
	"fmt"
)

// Config validation errors.
var (
	ErrEmptyAddr = errors.New("addr is empty")
)

// Config contains credentials for ClickHouse database.
type Config struct {
	Addr     string `yaml:"addr" json:"addr"`
	Database string `yaml:"database" json:"database"`
	Debug    bool   `yaml:"debug" json:"debug"`
	ZoneInfo string `yaml:"zoneinfo" json:"zone_info"`
}

// Validate checks required fields and validates for allowed values.
func (cfg Config) Validate() error {
	if cfg.Addr == "" {
		return ErrEmptyAddr
	}

	return nil
}

// GetDataSourceName returns Data Source Name connection string to ClickHouse database.
func (cfg *Config) GetDataSourceName() string {
	debug := "False"
	if cfg.Debug {
		debug = "True"
	}

	database := ""
	if cfg.Database != "" {
		database = "&database=" + cfg.Database
	}

	return fmt.Sprintf("tcp://%s?charset=utf8&parseTime=True&debug=%s%s", cfg.Addr, debug, database)
}
