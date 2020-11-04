package profiler

import (
	"net/http"

	// Для включения профилировщика подгружаем зависимость.
	_ "net/http/pprof"

	"github.com/outdead/echo-skeleton/internal/logger"
)

// Serve запускает HTTP сервер на переданном порту и позволяет
// профилировать сервис по ссылке {host}:{port}/debug/pprof/.
func Serve(port string, log *logger.Entry) {
	if port != "" {
		go func() {
			log.Error(http.ListenAndServe("localhost:"+port, nil))
		}()

		log.Infof("profiler started on port %s", port)
	}
}
