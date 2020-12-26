package clickhouse

import (
	"fmt"
	"strings"
)

// Model is an interface for ClickHouse data structures.
// Describes methods for getting the table name, field names and their values.
type Model interface {
	// TableName returns the table name.
	TableName() string

	// GetFields returns the names of the table fields.
	GetFields() []string

	// GetValues returns the values of the table fields.
	GetValues() []interface{}
}

// PrepareInsertionSQL returns a SQL prepare statement string to insert records
// into the database.
func PrepareInsertionSQL(model Model) string {
	fields := model.GetFields()
	binds := strings.Repeat("?,", len(fields))
	binds = binds[:len(binds)-1]

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", model.TableName(), strings.Join(fields, ", "), binds)
}
