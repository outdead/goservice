package clickhouse

import (
	"errors"
	"fmt"
	"os"
	"time"

	// Import ClickHouse driver.
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"
	"github.com/outdead/goservice/internal/utils/multierror"
)

// ErrLostConnection is returned when connection to database was lost.
var ErrLostConnection = errors.New("clickhouse: connection is lost")

// DB is a wrapper around sqlx.DB which keeps track of the ClickHouse database.
type DB struct {
	config *Config
	db     *sqlx.DB
}

// NewDB creates new connection to ClickHouse using sqlx.
func NewDB(cfg *Config) (*DB, error) {
	if cfg.ZoneInfo != "" {
		if err := os.Setenv("ZONEINFO", cfg.ZoneInfo); err != nil {
			return nil, fmt.Errorf("clickhouse: %w", err)
		}
	}

	db, err := sqlx.Open("clickhouse", cfg.GetDataSourceName())
	if err != nil {
		return nil, fmt.Errorf("clickhouse: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()

		return nil, fmt.Errorf("clickhouse: %w", err)
	}

	return &DB{config: cfg, db: db}, nil
}

// Config returns a pointer to the Config with which the connection was made.
func (db *DB) Config() *Config {
	return db.config
}

// DB returns pointer to sqlx.DB.
func (db *DB) DB() *sqlx.DB {
	return db.db
}

// IsConnected() checks connection status to database.
func (db *DB) IsConnected() bool {
	if db.db == nil {
		return false
	}

	if err := db.db.Ping(); err != nil {
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

	row := db.db.QueryRow("SELECT now()")
	err := row.Scan(&st)

	return st, fmt.Errorf("clickhouse: %w", err)
}

// MultiInsert performs a transactional insert of multiple records.
func (db *DB) MultiInsert(query string, rows [][]interface{}) error {
	if db.db == nil {
		return ErrLostConnection
	}

	tx, err := db.db.Begin()
	if err != nil {
		return fmt.Errorf("clickhouse: %w", err)
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return fmt.Errorf("clickhouse multiple errors: %w", multierror.New(err, err2))
		}

		return fmt.Errorf("clickhouse: %w", err)
	}

	defer stmt.Close()

	for i := range rows {
		if _, err := stmt.Exec(rows[i]...); err != nil {
			if err2 := tx.Rollback(); err2 != nil {
				return fmt.Errorf("clickhouse multiple errors: %w", multierror.New(err, err2))
			}

			return fmt.Errorf("clickhouse: %w", err)
		}
	}

	return tx.Commit()
}

// Close closes database connections.
func (db *DB) Close() error {
	if db.db == nil {
		return nil
	}

	return db.db.Close()
}
