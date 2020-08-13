package authentication

import (
	"net/http"
	"runtime/debug"
	"time"

	log "github.com/go-kit/kit/log"
)

type LoggingMW struct {
	Logger log.Logger
}

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, status: http.StatusOK}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	return
}



// LoggingMiddlewareReal logs the incoming HTTP request & its duration.
func LoggingMiddlewareReal(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal Server Error: From logger middleware"))
					logger.Log(
						"err", err,
						"trace", debug.Stack(),
					)
				}
			}()

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			logger.Log(
				"Status", wrapped.status,
				"Device", r.Header["User-Agent"],
				"Method", r.Method,
				"Path", r.URL.EscapedPath(),
				"Duration", time.Since(start),
			)
		}

		return http.HandlerFunc(fn)
	}
}
