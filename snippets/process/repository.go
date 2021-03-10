package process

import (
	"github.com/outdead/goservice/internal/utils/connector"
	"github.com/outdead/goservice/internal/utils/logutil"
)

// Repository describes getting and changing data methods.
type Repository interface {
}

// repository implements Repository.
type repository struct {
	logger *logutil.Entry
	conn   connector.Connector
}

// NewRepository creates and returns new repository.
func NewRepository(log *logutil.Entry, conn connector.Connector) Repository {
	return &repository{logger: log, conn: conn}
}
