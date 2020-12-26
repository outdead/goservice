package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-pg/pg/v9"
)

// QueryLogger is a queries logger.
type QueryLogger struct{}

// BeforeQuery implements BeforeQuery of the pg.QueryHook interface.
// Called before every query and does nothing.
func (d QueryLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

// AfterQuery prints the executed query to stdout.
// Called after the query has completed.
func (d QueryLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	fmt.Println(q.FormattedQuery())

	return nil
}

// ErrLostConnection is returned when connection to database was lost.
var ErrLostConnection = errors.New("postgres: connection is lost")

// DB is a wrapper around pg.DB which keeps track of the PostgreSQL database.
type DB struct {
	config *Config
	db     *pg.DB
}

// NewDB creates new connection to PostgreSQL using pg.v9.
func NewDB(cfg *Config) (*DB, error) {
	db := DB{config: cfg, db: pg.Connect(&pg.Options{
		Addr:        cfg.Addr,
		User:        cfg.User,
		Password:    cfg.Password,
		Database:    cfg.Database,
		PoolSize:    cfg.PoolSize,
		PoolTimeout: time.Hour,
	})}

	if cfg.Debug {
		db.db.AddQueryHook(QueryLogger{})
	}

	if _, err := db.GetServerTime(); err != nil {
		_ = db.Close()

		return nil, err
	}

	return &db, nil
}

// Config returns a pointer to the Config with which the connection was made.
func (db *DB) Config() *Config {
	return db.config
}

// DB returns pointer to pg.DB.
func (db *DB) DB() *pg.DB {
	return db.db
}

// IsConnected() checks connection status to database.
func (db *DB) IsConnected() bool {
	if db == nil {
		return false
	}

	if _, err := db.GetServerTime(); err != nil {
		return false
	}

	return true
}

// GetServerTime returns database server time or error.
func (db *DB) GetServerTime() (time.Time, error) {
	var st time.Time

	if db.db == nil {
		return st, ErrLostConnection
	}

	_, err := db.db.QueryOne(pg.Scan(&st), "SELECT now()")

	return st, err
}

// Close closes database connections.
func (db *DB) Close() error {
	if db.db == nil {
		return nil
	}

	return db.db.Close()
}
