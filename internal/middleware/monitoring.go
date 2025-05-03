package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/fungicibus/order/internal/logger"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

func NewMonitoringMiddleware(appVersion string, logger *logger.Logger) func(next http.Handler) http.Handler {
	var (
		version = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "version",
			Help: "Version information about this binary",
			ConstLabels: map[string]string{
				"version": appVersion,
			},
		})
		httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Count of all HTTP requests",
		}, []string{"code", "method"})

		httpRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of all HTTP requests",
		}, []string{"code", "handler", "method"})
	)
	version.Set(1)

	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(version)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isExcludedFromMonitoring(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			requestID := GetRequestID(r.Context())
			logger.Info().
				Str("requestID", requestID).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Msg("request started")

			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			duration := time.Since(start)
			statusCode := strconv.Itoa(ww.Status())

			httpRequestsTotal.WithLabelValues(statusCode, r.Method).Inc()
			httpRequestDuration.WithLabelValues(
				statusCode,
				r.URL.Path,
				r.Method,
			).Observe(duration.Seconds())

			logEvent := logger.Info()
			if ww.Status() >= 400 {
				logEvent = logger.Error()
			}
			logEvent.
				Str("requestID", requestID).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status", ww.Status()).
				Dur("duration", duration).
				Str("userAgent", r.UserAgent()).
				Msg("request completed")
		})
	}
}
