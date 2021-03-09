package process

import (
	"github.com/outdead/goservice/internal/connector"
	"github.com/outdead/goservice/internal/utils/logutils"
)

// Repository describes getting and changing data methods.
type Repository interface {
}

// repository implements Repository.
type repository struct {
	logger *logutils.Entry
	conn   connector.Connector
}

// NewRepository creates and returns new repository.
func NewRepository(log *logutils.Entry, conn connector.Connector) Repository {
	return &repository{logger: log, conn: conn}
}
