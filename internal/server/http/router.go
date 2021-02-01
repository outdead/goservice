package http

import "github.com/outdead/goservice/internal/server/http/handler/system"

func (s *Server) router() {
	root := s.echo.Group("")

	systemHandler := system.NewHandler()
	root.GET("/system/ping", systemHandler.Ping)
}
