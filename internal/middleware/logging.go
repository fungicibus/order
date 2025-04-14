package middleware

import (
	"net/http"

	"github.com/fungicibus/order/internal/logger"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func NewLoggingMiddleware(logger *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := GetRequestID(r.Context())
			logger.Info().
				Str("requestID", requestId).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Msg("request")

			rw := newResponseWriter(w)
			next.ServeHTTP(rw, r)

			logger.Info().
				Str("requestID", requestId).
				Int("status", rw.status).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Msg("response")
		})
	}
}
