package http

import "github.com/outdead/echo-skeleton/internal/server/http/handler/health"

func (s *Server) router() {
	root := s.echo.Group("")

	healthHandler := health.NewHandler()
	root.GET("/ping", healthHandler.Ping())
}
