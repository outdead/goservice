package rabbit

import (
	"errors"
	"fmt"

	"github.com/streadway/amqp"
)

// Validation errors.
var (
	ErrNoConsumers = errors.New("consumers are not set")

	ErrEmptyServer                = errors.New("server is empty")
	ErrEmptyExchangeType          = errors.New("exchange.type is empty")
	ErrEmptyRoutingKeyOrQueueName = errors.New("routing_key or queue.name is empty")
	ErrEmptyRoutingKey            = errors.New("routing_key is empty")
	ErrEmptyQueueName             = errors.New("queue name is empty")
)

// Config contains credentials for RabbitMQ.
type Config struct {
	Server    ServerConfig              `yaml:"server" json:"server"`
	Consumers map[string]ConsumerConfig `yaml:"consumers" json:"consumers"`
}

// Validate checks required fields and validates for allowed values.
func (cfg *Config) Validate() error {
	if err := cfg.Server.Validate(); err != nil {
		return fmt.Errorf("server: %w", err)
	}

	if cfg.Consumers == nil {
		return fmt.Errorf("consumers: %w", ErrNoConsumers)
	}

	for name, consumer := range cfg.Consumers {
		if err := consumer.Validate(); err != nil {
			return fmt.Errorf("consumers.%s: %w", name, err)
		}
	}

	return nil
}

// Config contains credentials for RabbitMQ queue.
type ConsumerConfig struct {
	QueueName  string `yaml:"queue_name" json:"queue_name"`
	RoutingKey string `yaml:"routing_key" json:"routing_key"`
}

// Validate checks required fields and validates for allowed values.
func (cfg *ConsumerConfig) Validate() error {
	if cfg.QueueName == "" {
		return ErrEmptyQueueName
	}

	if cfg.RoutingKey == "" {
		return ErrEmptyRoutingKey
	}

	return nil
}

// Config contains credentials for RabbitMQ server.
type ServerConfig struct {
	Server   string `yaml:"server" json:"server"`
	Exchange struct {
		Name       string `yaml:"name" json:"name"`
		Type       string `yaml:"type" json:"type"`
		AutoDelete bool   `yaml:"auto_delete" json:"auto_delete"`
		Durable    bool   `yaml:"durable" json:"durable"`
	} `yaml:"exchange" json:"exchange"`
	Queue struct {
		Name       string     `yaml:"name" json:"name"`
		AutoDelete bool       `yaml:"auto_delete" json:"auto_delete"`
		Durable    bool       `yaml:"durable" json:"durable"`
		Exclusive  bool       `yaml:"exclusive" json:"exclusive"`
		Args       amqp.Table `yaml:"arguments" json:"args"`
	} `yaml:"queue" json:"queue"`
	RoutingKey string `yaml:"routing_key" json:"routing_key"`
	Qos        int    `yaml:"qos" json:"qos"`
}

// Validate checks required fields and validates for allowed values.
func (cfg *ServerConfig) Validate() error {
	if cfg.Server == "" {
		return ErrEmptyServer
	}

	if cfg.Exchange.Type == "" {
		return ErrEmptyExchangeType
	}

	// If we want to publish, then routing_key must not be empty, if we want to
	// consume, then queue.name must not be empty
	if cfg.RoutingKey == "" && cfg.Queue.Name == "" {
		return ErrEmptyRoutingKeyOrQueueName
	}

	return nil
}
