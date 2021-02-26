package connector

import (
	"fmt"

	"github.com/outdead/goservice/internal/connector/clickhouse"
	"github.com/outdead/goservice/internal/connector/elasticsearch"
	"github.com/outdead/goservice/internal/connector/postgres"
	"github.com/outdead/goservice/internal/connector/rabbit"
)

// Config contains credentials for databases.
type Config struct {
	Postgres      postgres.Config      `yaml:"postgres" json:"postgres"`
	Clickhouse    clickhouse.Config    `yaml:"clickhouse" json:"clickhouse"`
	Elasticsearch elasticsearch.Config `yaml:"elasticsearch" json:"elasticsearch"`
	RabbitMQ      rabbit.Config        `yaml:"rabbitmq" json:"rabbit_mq"`
}

// Validate checks required fields and validates for allowed values.
func (cfg *Config) Validate() error {
	if err := cfg.Postgres.Validate(); err != nil {
		return fmt.Errorf("postgres: %w", err)
	}

	if err := cfg.Clickhouse.Validate(); err != nil {
		return fmt.Errorf("clickhouse: %w", err)
	}

	if err := cfg.Elasticsearch.Validate(); err != nil {
		return fmt.Errorf("elasticsearch: %w", err)
	}

	if err := cfg.RabbitMQ.Validate(); err != nil {
		return fmt.Errorf("rabbitmq: %w", err)
	}

	return nil
}
