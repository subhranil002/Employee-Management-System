package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(
	statusCode int,
) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *statusWriter) Write(
	data []byte,
) (int, error) {

	if w.status == 0 {
		w.status = http.StatusOK
	}

	return w.ResponseWriter.Write(data)
}

// Logger logs HTTP request method, path, status, and duration
func Logger(
	logger *slog.Logger,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				start := time.Now()

				// Wrap ResponseWriter to capture HTTP status code
				sw := &statusWriter{
					ResponseWriter: w,
				}

				next.ServeHTTP(sw, r)

				status := sw.status

				if status == 0 {
					status = http.StatusOK
				}

				logger.Info(
					"http request",
					"method", r.Method,
					"path", r.URL.Path,
					"status", status,
					"duration", time.Since(start),
				)
			},
		)
	}
}
