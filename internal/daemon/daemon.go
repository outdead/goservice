package daemon

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/outdead/goservice/internal/connector"
	"github.com/outdead/goservice/internal/server/http"
	"github.com/outdead/goservice/internal/server/profiler"
	"github.com/outdead/goservice/internal/utils/logutils"
)

type Daemon struct {
	config *Config
	logger *logutils.Entry
	errors chan error

	conn   connector.Connector
	server struct {
		http *http.Server
	}
}

func NewDaemon(cfg *Config, log *logutils.Entry) *Daemon {
	d := Daemon{
		config: cfg,
		errors: make(chan error, cfg.App.ErrorBuffer),
		logger: log,
	}

	return &d
}

func (d *Daemon) Close() error {
	return d.close()
}

// Run starts the Daemon.
func (d *Daemon) Run() error {
	if err := d.init(); err != nil {
		return err
	}

	// Creates goroutine process for start profiler.
	profiler.Serve("0.0.0.0"+d.config.App.ProfilerPort, d.logger)

	// Creates goroutine process for start HTTP server.
	d.server.http.Serve(d.config.App.Port)

	interrupter := make(chan os.Signal, 1)
	signal.Notify(interrupter, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	ticker := time.NewTicker(d.config.App.CheckConnectionsInterval)
	defer ticker.Stop()

	d.logger.Info("start daemon success")

Loop:
	for {
		select {
		case <-interrupter:
			d.logger.Info("received an interrupt, unsubscribe and closing connections...")

			break Loop
		case err := <-d.Errors(): // down here
			d.logger.Info("daemon fatal error occurred, unsubscribe and closing connections...")

			return err
		case <-ticker.C:
			d.logger.Debug("check connections is not implemented")
		case err := <-d.server.http.Errors():
			d.reportError(err) // TODO: Try to recreate http server.
		}
	}

	return nil
}

func (d *Daemon) Errors() <-chan error {
	return d.errors
}

func (d *Daemon) init() error {
	if d.logger == nil {
		d.logger = logutils.New().NewEntry()
	}

	var err error

	if d.conn, err = connector.New(&d.config.Connections); err != nil {
		return fmt.Errorf("connector: %w", err)
	}

	d.server.http = http.NewServer(d.logger)

	return nil
}

func (d *Daemon) close() error {
	d.logger.Debug("stopping daemon...")

	var errs []error

	if err := d.server.http.Close(); err != nil {
		errs = append(errs, err)
	}

	if d.conn != nil {
		if err := d.conn.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		for _, e := range errs {
			d.logger.Errorf("close daemon error: %s", e)
		}

		d.logger.Errorf("daemon closed with %d errors", len(errs))

		return nil
	}

	d.logger.Info("stop daemon success")

	return nil
}

// reportError publishes error to the errors channel.
// if you do not read errors from the errors channel then after the channel
// buffer overflows the application exits with a fatal level and the
// os.Exit(1) exit code.
func (d *Daemon) reportError(err error) {
	if err != nil {
		select {
		case d.errors <- err:
		default:
			// IMPORTANT: This is a soft version of the application panic.
			d.logger.Fatalf("daemon error channel is locked: %v", err)
		}
	}
}
