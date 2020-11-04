package httpserver

import "github.com/outdead/echo-skeleton/internal/api/httpserver/handler/health"

func (s *Server) router() {
	root := s.echo.Group("")

	healthHandler := health.NewHandler()
	root.GET("/ping", healthHandler.Ping())
}
