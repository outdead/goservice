package http

import "github.com/outdead/echo-skeleton/internal/server/http/handler/system"

func (s *Server) router() {
	root := s.echo.Group("")

	systemHandler := system.NewHandler()
	root.GET("/ping", systemHandler.Ping())
}
