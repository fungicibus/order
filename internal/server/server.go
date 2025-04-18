package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"

	"github.com/fungicibus/order/config"
	"github.com/fungicibus/order/internal/logger"
	"github.com/fungicibus/order/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Server struct {
	cfg    *config.Config
	logger *logger.Logger
	srv    *http.Server
	v1     http.Handler
}

func New(cfg *config.Config, logger *logger.Logger, v1 http.Handler) *Server {
	return &Server{
		cfg:    cfg,
		logger: logger,
		v1:     v1,
	}
}

func (s *Server) Run(ctx context.Context) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", s.cfg.Server.Port),
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		Handler:      s.getRouter(),
		ReadTimeout:  s.cfg.Server.ReadTimeout,
		WriteTimeout: s.cfg.Server.WriteTimeout,
	}
	s.srv = srv

	s.logger.Info().Msgf("server started on port %d", s.cfg.Server.Port)
	return srv.ListenAndServe()
}

func (s *Server) getRouter() *chi.Mux {
	router := chi.NewMux()

	// Middleware
	router.Use(middleware.RequestIDMiddleware)
	router.Use(middleware.NewMonitoringMiddleware(s.cfg.App.Name, s.logger))

	// Profiler
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/pprof/block", pprof.Handler("block"))

	// Metrics
	router.Handle("/metrics", promhttp.Handler())

	// Swagger
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	router.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, s.cfg.OpenapiPath)
	})

	// API
	router.Mount("/api/v1", s.v1)

	return router
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info().Msg("graceful server shutdown")
	err := s.srv.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("failed to shutdown http server: %w", err)
	}
	return nil
}
