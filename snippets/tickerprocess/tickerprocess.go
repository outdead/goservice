package tickerprocess

import (
	"sync"
	"time"

	"github.com/outdead/goservice/internal/utils/logutil"
)

// Repository describes getting and changing data methods.
type Repository interface {
}

// Process performs work in a separate goroutine.
type Process struct {
	config *Config
	logger *logutil.Entry
	errors chan error

	repo Repository

	// Sync.
	quit    chan bool
	started bool
	wg      sync.WaitGroup
}

// NewProcess creates and returns new Process.
func NewProcess(cfg *Config, repo Repository, log *logutil.Entry) *Process {
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

	if p.started {
		p.logger.Warning("process already been started")

		return
	}

	p.quit = make(chan bool, 1)
	p.started = true

	p.wg.Add(1)

	go p.run()
}

// Quit stops all processes.
func (p *Process) Quit() {
	if p.config.Disabled {
		p.logger.Debug("cannot quit disabled process")

		return
	}

	if p.quit == nil || !p.started {
		p.logger.Debug("cannot quit stopped process")

		return
	}

	select {
	case p.quit <- true:
		p.wg.Wait()
	default:
		p.logger.Debug("process quit already been called")
	}
}

// ReportError publishes error to the errors channel.
// if you do not read errors from the errors channel then after the channel
// buffer overflows the application exits with a fatal level and the
// os.Exit(1) exit code.
func (p *Process) ReportError(err error) {
	if err != nil {
		select {
		case p.errors <- err:
		default:
			// IMPORTANT: Фактически это мягкий вариант паники приложения.
			p.logger.Fatalf("process error channel is locked: %v", err)
		}
	}
}

func (p *Process) run() {
	defer func() {
		p.started = false
		p.logger.Info("process stopped")
		p.wg.Done()
	}()

	ticker := time.NewTicker(p.config.StartInterval)
	defer ticker.Stop()

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
