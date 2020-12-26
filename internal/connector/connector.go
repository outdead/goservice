package connector

import (
	"io"

	"github.com/outdead/goservice/internal/connector/clickhouse"
	"github.com/outdead/goservice/internal/connector/postgres"
	"github.com/outdead/goservice/internal/connector/rabbit"
	"github.com/outdead/goservice/internal/utils/multierror"
)

// Connector is the interface for databases accessing.
type Connector interface {
	io.Closer
	CH() *clickhouse.DB
	PG() *postgres.DB
	RMQ() *rabbit.Client
}

type connector struct {
	ch  *clickhouse.DB
	pg  *postgres.DB
	rmq *rabbit.Client
}

// New establishes new connections from configuration parameters.
func New(cfg *Config) (Connector, error) {
	conn := connector{}
	var err error

	if conn.ch, err = clickhouse.NewDB(&cfg.Clickhouse); err != nil {
		return nil, conn.close(err)
	}

	if conn.pg, err = postgres.NewDB(&cfg.Postgres); err != nil {
		return nil, err
	}

	if conn.rmq, err = rabbit.NewClient(&cfg.RabbitMQ); err != nil {
		return nil, conn.close(err)
	}

	return &conn, nil
}

// CH returns pointer to clickhouse.DB.
func (conn *connector) CH() *clickhouse.DB {
	return conn.ch
}

// PG returns pointer to postgres.DB.
func (conn *connector) PG() *postgres.DB {
	return conn.pg
}

// RMQ returns pointer to rabbit.Client.
func (conn *connector) RMQ() *rabbit.Client {
	return conn.rmq
}

// Close closes all databases connections.
func (conn *connector) Close() error {
	return conn.close()
}

func (conn *connector) close(prevErrs ...error) error {
	errs := multierror.New(prevErrs...)

	if conn.pg != nil {
		if err := conn.pg.Close(); err != nil {
			errs.Append(err)
		}
	}

	if conn.ch != nil {
		if err := conn.ch.Close(); err != nil {
			errs.Append(err)
		}
	}

	if errs.Len() != 0 {
		return errs
	}

	return nil
}
