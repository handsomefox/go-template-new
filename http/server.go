package http

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"project-template/config"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	mux        *chi.Mux
	httpServer *http.Server
	errc       chan error
}

func NewServer(c *config.Server) (*Server, error) {
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)

	mux := chi.NewMux()
	return &Server{
		mux: mux,
		httpServer: &http.Server{
			Addr:              addr,
			Handler:           mux,
			ReadTimeout:       10 * time.Second,
			ReadHeaderTimeout: 3 * time.Second,
			WriteTimeout:      120 * time.Second,
			IdleTimeout:       30 * time.Second,
			ErrorLog:          log.Default(),
		},
		errc: make(chan error),
	}, nil
}

func (s *Server) Run(ctx context.Context) {
	slog.LogAttrs(ctx, slog.LevelInfo, "Starting server", slog.String("addr", s.httpServer.Addr))
	if err := s.httpServer.ListenAndServe(); err != nil {
		s.errc <- err
	}
}

func (s *Server) Close() error {
	return s.httpServer.Close()
}

func (s *Server) Notify() <-chan error {
	return s.errc
}

func (s *Server) MountSubrouter(path string, fn func(r *chi.Mux)) {
	subrouter := chi.NewRouter()
	fn(subrouter)
	s.mux.Mount(path, subrouter)
}

func (s *Server) Mount(fn func(r *chi.Mux)) {
	fn(s.mux)
}

func (s *Server) Use(middlewares ...func(h http.Handler) http.Handler) {
	s.mux.Use(middlewares...)
}
