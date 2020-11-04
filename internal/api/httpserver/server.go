package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/outdead/echo-skeleton/internal/api/httpserver/response"
	"github.com/outdead/echo-skeleton/internal/logger"
)

// ShutdownTimeOut is time to terminate queries when quit signal given.
const ShutdownTimeOut = 10 * time.Second

// ErrLockedServer возвращается при повторном вызове останвки HTTP сервер.
var ErrLockedServer = errors.New("http api server is locked")

type ServerInterface interface {
	Run()
	Close() error
	Errors() <-chan error
}

type Server struct {
	logger *logger.Entry
	errors chan error
	quit   chan bool
	wg     sync.WaitGroup

	echo *echo.Echo
}

func NewServer(log *logger.Entry) *Server {
	s := Server{
		logger: log,
		errors: make(chan error, 100),
		quit:   make(chan bool),
		echo:   echo.New(),
	}

	return &s
}

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

func (s *Server) Close() error {
	select {
	case s.quit <- true:
		s.wg.Wait()
		s.logger.Info("stop http api success")

		return nil
	default:
		return ErrLockedServer
	}
}

func (s *Server) Errors() <-chan error {
	return s.errors
}

func (s *Server) newEcho() *echo.Echo {
	e := echo.New()

	e.Logger.SetOutput(s.logger.Writer())
	e.Use(middleware.Recover())
	e.HideBanner = true
	e.HidePort = true

	e.HTTPErrorHandler = s.customHTTPErrorHandler

	return e
}

// customHTTPErrorHandler customizes error response.
// @source: https://github.com/labstack/echo/issues/325
func (s *Server) customHTTPErrorHandler(err error, c echo.Context) {
	switch t := err.(type) {
	case *echo.HTTPError:
		errorCode := t.Code
		switch errorCode {
		case http.StatusNotFound, http.StatusMethodNotAllowed:
			if err := response.ServeNotFoundError(c); err != nil {
				s.reportError(err)
			}
		default:
			s.logger.WithField("url", c.Path()).Errorf("unexpected http code: %d", errorCode)

			if err := response.ServeInternalServerError(c); err != nil {
				s.reportError(err)
			}
		}
	default:
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
			// IMPORTANT: Пердполагается, что канал ошибок пробрасывается
			// функцией Errors() и его читает вызывающая рутина. Если ошибки
			// никто не читает, то после переполнения буфера ошибок фаталимся.
			s.logger.Fatalf("http api server error channel is locked: %s", err)
		}
	}
}
