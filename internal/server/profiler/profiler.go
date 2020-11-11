package profiler

import (
	"net/http"

	// Load the dependency to enable the profiler.
	_ "net/http/pprof"

	"github.com/outdead/echo-skeleton/internal/logger"
)

// Serve starts an HTTP server on the given port and allows you to profile the
// service by reference {host}:{port}/debug/pprof/.
func Serve(port string, log *logger.Entry) {
	if port != "" {
		go func() {
			log.Error(http.ListenAndServe("localhost:"+port, nil))
		}()

		log.Infof("profiler started on port %s", port)
	}
}
