package profiler

import (
	"net/http"

	// Load the dependency to enable the profiler.
	_ "net/http/pprof"

	"github.com/outdead/goservice/internal/utils/logutils"
)

// Serve starts an HTTP server on the given port and allows you to profile the
// service by reference {host}:{port}/debug/pprof/.
func Serve(addr string, log *logutils.Entry) {
	if addr != "" {
		go func() {
			log.Error(http.ListenAndServe(addr, nil))
		}()

		log.Infof("profiler started on %s", addr)
	}
}
