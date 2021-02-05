package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/outdead/goservice/internal/server/http/middleware"
	"github.com/outdead/goservice/internal/server/http/response"
	"github.com/outdead/goservice/internal/utils/logutils"
)

// ShutdownTimeOut is time to terminate queries when quit signal given.
const ShutdownTimeOut = 10 * time.Second

// ErrLockedServer returned on repeated call Close() the HTTP server.
var ErrLockedServer = errors.New("http api server is locked")

// Server defines parameters for running an HTTP server.
type Server struct {
	logger *logutils.Entry
	errors chan error
	quit   chan bool
	wg     sync.WaitGroup

	echo *echo.Echo
}

// NewServer allocates and returns a new Server.
func NewServer(log *logutils.Entry) *Server {
	s := Server{
		logger: log,
		errors: make(chan error, 100),
		quit:   make(chan bool),
		echo:   echo.New(),
	}

	return &s
}

// Serve initializes HTTP Server and runs it on received port.
func (s *Server) Serve(port string) {
	s.echo = s.newEcho()
	s.router()

	go func() {
		if err := s.echo.Start(":" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			// Report error if server is not closed by Echo#Shutdown.
			s.reportError(fmt.Errorf("start http server error: %w", err))
		}
	}()

	s.logger.Infof("http server started on port %s", port)

	s.quit = make(chan bool)
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		<-s.quit
		s.logger.Debug("stopping http api server...")

		ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeOut)
		defer cancel()

		if s.echo != nil {
			if err := s.echo.Shutdown(ctx); err != nil {
				s.reportError(fmt.Errorf("shutdown http server error: %w", err))
			}
		}
	}()
}

// Close stops HTTP Server.
func (s *Server) Close() error {
	if s.quit == nil {
		return ErrLockedServer
	}

	select {
	case s.quit <- true:
		s.wg.Wait()
		s.logger.Info("stop http api success")

		return nil
	default:
		return ErrLockedServer
	}
}

// Errors returns errors channel.
func (s *Server) Errors() <-chan error {
	return s.errors
}

func (s *Server) newEcho() *echo.Echo {
	e := echo.New()

	e.Logger.SetOutput(s.logger.Writer())
	e.Use(middleware.Recover())
	e.Use(middleware.WithContext())
	e.HideBanner = true
	e.HidePort = true

	e.HTTPErrorHandler = s.httpErrorHandler

	return e
}

// httpErrorHandler customizes error response.
// @source: https://github.com/labstack/echo/issues/325
func (s *Server) httpErrorHandler(err error, c echo.Context) {
	var t *echo.HTTPError
	if errors.As(err, &t) {
		switch t.Code {
		case http.StatusNotFound, http.StatusMethodNotAllowed:
			if err := response.ServeNotFoundError(c); err != nil {
				s.reportError(err)
			}
		default:
			s.logger.WithField("url", c.Path()).Errorf("unexpected http code: %d", t.Code)

			if err := response.ServeInternalServerError(c); err != nil {
				s.reportError(err)
			}
		}
	} else {
		s.logger.WithField("url", c.Path()).Errorf("unexpected http error: %s", err)

		if err := response.ServeInternalServerError(c); err != nil {
			s.reportError(err)
		}
	}
}

func (s *Server) reportError(err error) {
	if err != nil {
		select {
		case s.errors <- err:
		default:
			// IMPORTANT: It is assumed that the error channel is forwarded by
			// the Errors() function and the calling routine reads it. If no one
			// reads the errors, then after the error buffer overflows we are
			// exit with fatal level.
			s.logger.Fatalf("http api server error channel is locked: %s", err)
		}
	}
}
