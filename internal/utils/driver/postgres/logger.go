package postgres

import (
	"context"
	"fmt"
	"io"

	"github.com/go-pg/pg/v9"
)

// QueryLogger is a queries logger.
type QueryLogger struct {
	w io.Writer
}

// NewQueryLogger creates and returns pointer to QueryLogger.
func NewQueryLogger(w io.Writer) *QueryLogger {
	return &QueryLogger{w: w}
}

// BeforeQuery implements BeforeQuery of the pg.QueryHook interface.
// Called before every query and does nothing.
func (d QueryLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

// AfterQuery prints the executed query to stdout.
// Called after the query has completed.
func (d QueryLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	query, err := q.FormattedQuery()

	fmt.Fprintln(d.w, query, err)

	return nil
}
