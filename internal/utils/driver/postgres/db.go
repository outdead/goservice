package postgres

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-pg/pg/v9"
)

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
		db.db.AddQueryHook(NewQueryLogger(os.Stdout))
	}

	if _, err := db.GetServerTime(); err != nil {
		_ = db.Close()

		return nil, err
	}

	return &db, nil
}

// Config returns a pointer to the Dialer with which the connection was made.
func (db *DB) Config() *Config {
	return db.config
}

// DB returns pointer to pg.DB.
func (db *DB) DB() *pg.DB {
	return db.db
}

// IsConnected checks connection status to database.
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

	if _, err := db.db.QueryOne(pg.Scan(&st), "SELECT now()"); err != nil {
		return st, fmt.Errorf("postgres: %w", err)
	}

	return st, nil
}

// Close closes database connections.
func (db *DB) Close() error {
	if db.db == nil {
		return nil
	}

	return db.db.Close()
}

// IsErrNoRows returns true if error is pg.ErrNoRows.
func (db *DB) IsErrNoRows(err error) bool {
	return errors.Is(err, pg.ErrNoRows)
}
