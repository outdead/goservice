package process

import (
	"github.com/outdead/goservice/internal/connector"
	"github.com/outdead/goservice/internal/utils/logutils"
)

// DataSupplier describes getting data methods.
type DataSupplier interface {
}

// DataModifier describes changing data methods.
type DataModifier interface {
}

// DataSupplyModifier describes getting and changing data methods.
type DataSupplyModifier interface {
	DataSupplier
	DataModifier
}

// Repository implements DataSupplyModifier.
type Repository struct {
	logger *logutils.Entry
	conn   connector.Connector
}

// NewRepository creates and returns new Repository.
func NewRepository(log *logutils.Entry, conn connector.Connector) *Repository {
	return &Repository{logger: log, conn: conn}
}
