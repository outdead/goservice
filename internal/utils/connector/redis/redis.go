package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// ErrNoRows is returned when Get returned zero records.
var ErrNoRows = errors.New("key not found")

// Client is a database handle representing connection to Redis.
type Client struct {
	config *Config
	conn   *redis.Client
}

// NewClient creates and returns new Redis Client.
func NewClient(cfg *Config) (*Client, error) {
	// The go-redis package used sets localhost: 6379 as default if no value is
	// set. Remove this unobvious behavior and require to always specify the value
	// of the address with Redis.
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolSize:     cfg.PoolSize,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		_ = client.Close()

		return nil, fmt.Errorf("redis: create connection: %w", err)
	}

	return &Client{config: cfg, conn: client}, nil
}

// IsConnected checks connection status to database.
func (client *Client) IsConnected() bool {
	if client.conn == nil {
		return false
	}

	if _, err := client.conn.Ping(context.Background()).Result(); err != nil {
		return false
	}

	return true
}

// Dialer возвращает указатель на Dialer, с которым было совершено подключение.
func (client *Client) Config() *Config {
	return client.config
}

// Conn возвращает указатель на базовое соединение с Redis.
func (client *Client) Conn() *redis.Client {
	return client.conn
}

// Set sets value by key to Redis with ttl.
func (client *Client) Set(key string, data interface{}) error {
	cmd := client.conn.Set(context.Background(), key, data, client.config.TTL)

	return cmd.Err()
}

// GetValue gets value by key from Redis.
func (client *Client) Get(key string) (string, error) {
	val, err := client.conn.Get(context.Background(), key).Result()
	if err == redis.Nil { //nolint // this is still required according to go-redis documentation
		return "", ErrNoRows
	} else if err != nil {
		return "", fmt.Errorf("redis: %w", err)
	}

	return val, nil
}

// Del deletes value from Redis by key.
func (client *Client) Del(key string) error {
	cmd := client.conn.Del(context.Background(), key)

	return cmd.Err()
}

// IsErrNoRows returns true if err is ErrNoRows.
func (client *Client) IsErrNoRows(err error) bool {
	return errors.Is(err, ErrNoRows)
}

// Close closes connection to database.
func (client *Client) Close() error {
	if client.conn == nil {
		return nil
	}

	return client.conn.Close()
}
