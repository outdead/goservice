package process

import (
	"sync"
	"time"

	"github.com/outdead/goservice/internal/connector"
	"github.com/outdead/goservice/internal/utils/logutils"
)

type Process struct {
	config *Config
	logger *logutils.Entry
	conn   connector.Connector
	errors chan error

	// Sync.
	quit chan bool
	wg   sync.WaitGroup
}

func NewProcess(cfg *Config, log *logutils.Entry, conn connector.Connector) *Process {
	return &Process{
		config: cfg,
		logger: log,
		conn:   conn,
		errors: make(chan error, 100),
	}
}

func (p *Process) Errors() <-chan error {
	return p.errors
}

func (p *Process) Run() {
	if p.config.Disabled {
		p.logger.Debug("cannot run disabled process")

		return
	}

	p.logger.Info("process started")
	defer p.logger.Info("process stopped")

	ticker := time.NewTicker(p.config.StartInterval)
	defer ticker.Stop()

	p.quit = make(chan bool)

	p.wg.Add(1)
	defer p.wg.Done()

	for {
		select {
		case <-ticker.C:
			p.logger.Debug("process tick...")
		case <-p.quit:
			p.logger.Debug("process quit...")

			return
		}
	}
}

func (p *Process) Quit() {
	if p.config.Disabled {
		p.logger.Debug("cannot quit disabled process")

		return
	}

	if p.quit != nil {
		select {
		case p.quit <- true:
			p.wg.Wait()
		default:
			p.logger.Debug("process quit already been called")
		}
	}
}
