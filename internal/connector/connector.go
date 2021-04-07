package connector

import (
	"io"

	"github.com/outdead/goservice/internal/utils/driver/clickhouse"
	"github.com/outdead/goservice/internal/utils/driver/elasticsearch"
	"github.com/outdead/goservice/internal/utils/driver/postgres"
	"github.com/outdead/goservice/internal/utils/driver/rabbit"
	"github.com/outdead/goservice/internal/utils/driver/redis"
	"github.com/outdead/goservice/internal/utils/multierror"
)

// Connector is the interface for databases accessing.
type Connector interface {
	io.Closer
	CheckConnections() error
	IsErrNotFound(err error) bool

	PG() *postgres.DB
	CH() *clickhouse.DB
	ELA() *elasticsearch.Client
	Redis() *redis.Client
	RMQ() *rabbit.Client
}

type connector struct {
	pg    *postgres.DB
	ch    *clickhouse.DB
	ela   *elasticsearch.Client
	redis *redis.Client
	rmq   *rabbit.Client
}

// New establishes new connections from configuration parameters.
func New(cfg *Config) (Connector, error) {
	conn := connector{}
	var err error

	if conn.pg, err = postgres.NewDB(&cfg.Postgres); err != nil {
		return nil, conn.close(err)
	}

	if conn.ch, err = clickhouse.NewDB(&cfg.Clickhouse); err != nil {
		return nil, conn.close(err)
	}

	if conn.ela, err = elasticsearch.NewClient(&cfg.Elasticsearch); err != nil {
		return nil, conn.close(err)
	}

	if conn.redis, err = redis.NewClient(&cfg.Redis); err != nil {
		return nil, conn.close(err)
	}

	if conn.rmq, err = rabbit.NewClient(&cfg.RabbitMQ); err != nil {
		return nil, conn.close(err)
	}

	return &conn, nil
}

// CheckConnections checks database connections.
// TODO: Add multierror.
func (conn *connector) CheckConnections() error {
	if ok := conn.PG().IsConnected(); !ok {
		return postgres.ErrLostConnection
	}

	if ok := conn.CH().IsConnected(); !ok {
		return clickhouse.ErrLostConnection
	}

	if ok := conn.Redis().IsConnected(); !ok {
		return redis.ErrLostConnection
	}

	return nil
}

// IsErrNotFound returns true if the passed error indicates that there is
// no data in the database.
func (conn *connector) IsErrNotFound(err error) bool {
	return conn.ELA().IsErrNotFound(err) ||
		conn.PG().IsErrNoRows(err) ||
		conn.Redis().IsErrNoRows(err)
}

// CH returns pointer to clickhouse.DB.
func (conn *connector) CH() *clickhouse.DB {
	return conn.ch
}

// PG returns pointer to postgres.DB.
func (conn *connector) PG() *postgres.DB {
	return conn.pg
}

// ELA returns pointer to elasticsearch.Client.
func (conn *connector) ELA() *elasticsearch.Client {
	return conn.ela
}

// Redis returns pointer to redis.Client.
func (conn *connector) Redis() *redis.Client {
	return conn.redis
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

	if conn.ela != nil {
		conn.ela.Close()
	}

	if conn.redis != nil {
		if err := conn.redis.Close(); err != nil {
			errs.Append(err)
		}
	}

	if errs.Len() != 0 {
		return errs
	}

	return nil
}
