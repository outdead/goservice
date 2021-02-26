package elasticsearch

import (
	"context"
	"errors"
	"fmt"

	"github.com/olivere/elastic"
)

// DefaultBatchLimit contains the default value for the
// multi-insertion elements limit.
const DefaultBatchLimit = 5000

// ErrLostConnection is returned when connection to database was lost.
var ErrLostConnection = errors.New("elasticsearch: connection is lost")

// Client is a wrapper around elastic.Client which keeps track of the Elasticsearch
// database.
type Client struct {
	config *Config
	conn   *elastic.Client
	ctx    context.Context
}

// NewDB creates new connection to Elasticsearch using olivere/elastic.
func NewClient(cfg *Config) (*Client, error) {
	if cfg.HealthcheckInterval == 0 {
		cfg.HealthcheckInterval = DefaultHealthcheckInterval
	}

	conn, err := elastic.NewClient(
		elastic.SetSniff(true),
		elastic.SetURL(cfg.Addr),
		elastic.SetHealthcheckInterval(cfg.HealthcheckInterval),
	)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	client := Client{
		config: cfg,
		conn:   conn,
		ctx:    context.Background(),
	}

	return &client, nil
}

// Config returns a pointer to the Config with which the connection was made.
func (client *Client) Config() *Config {
	return client.config
}

// Conn returns pointer to elastic.Client.
func (client *Client) Conn() *elastic.Client {
	return client.conn
}

// MultiInsert performs a bulk insert of multiple records.
func (client *Client) MultiInsert(rows []Model) error {
	if client.conn == nil {
		return ErrLostConnection
	}

	bulk := client.conn.Bulk()

	for _, row := range rows {
		req := elastic.NewBulkIndexRequest().
			Index(client.config.Database).
			Type(row.TableName()).
			Id(row.CalculateID()).
			Doc(row)

		bulk = bulk.Add(req)
	}

	if _, err := bulk.Do(client.ctx); err != nil {
		return fmt.Errorf("elasticsearch: %w", err)
	}

	return nil
}

// Close tops the background processes that the client is running.
func (client *Client) Close() {
	if client.conn == nil {
		return
	}

	client.conn.Stop()
}
