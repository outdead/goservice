package process

import (
	"sync"
	"time"

	"github.com/outdead/goservice/internal/utils/logutils"
)

// Process performs work in a separate goroutine.
type Process struct {
	config *Config
	logger *logutils.Entry
	repo   Repository
	errors chan error

	// Sync.
	quit chan bool
	wg   sync.WaitGroup
}

// NewProcess creates and returns new Process.
func NewProcess(cfg *Config, log *logutils.Entry, repo Repository) *Process {
	return &Process{
		config: cfg,
		logger: log,
		repo:   repo,
		errors: make(chan error, 100),
	}
}

// Errors returns errors channel.
func (p *Process) Errors() <-chan error {
	return p.errors
}

// Run starts goroutine process.
func (p *Process) Run() {
	if p.config.Disabled {
		p.logger.Debug("cannot run disabled process")

		return
	}

	p.logger.Info("process started")

	ticker := time.NewTicker(p.config.StartInterval)
	defer ticker.Stop()

	p.quit = make(chan bool, 1)

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

// Quit stops all processes.
func (p *Process) Quit() {
	if p.config.Disabled {
		p.logger.Debug("cannot quit disabled process")

		return
	}

	if p.quit != nil {
		select {
		case p.quit <- true:
			p.wg.Wait()
			p.logger.Info("process stopped")
		default:
			p.logger.Debug("process quit already been called")
		}
	}
}
