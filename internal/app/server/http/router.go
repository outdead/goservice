package http

import (
	_ "github.com/outdead/goservice/internal/app/server/http/docs"
	"github.com/outdead/goservice/internal/app/server/http/handler/system"
	swagger "github.com/swaggo/echo-swagger"
)

//go:generate swag init -g router.go

//
// @title Go Service Example API
// @version 0.0.0-develop
// @description This is a goservice sample server.
// @host localhost:8080
// @BasePath /
//
// router creates routing.
func (s *Server) router() {
	root := s.echo.Group("")

	systemHandler := system.NewHandler()
	root.GET("/system/ping", systemHandler.Ping)

	root.GET("/swagger/*", swagger.WrapHandler)
}
