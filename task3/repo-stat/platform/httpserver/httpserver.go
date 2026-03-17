package httpserver

import (
	"context"
	"net/http"
	"time"
)

type Config struct {
	Address string `yaml:"address" default:"8080"`
	Timeout time.Duration `yaml:"timeout" default:"5s"`
}

type Server struct {
	srv *http.Server
}

func New(cfg Config, handler http.Handler) *Server {
	return &Server{
		srv: &http.Server{
			Addr:         cfg.Address,
			ReadTimeout:  cfg.Timeout,
			WriteTimeout: cfg.Timeout,
			IdleTimeout:  30 * time.Second,
			Handler:      handler,
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = s.srv.Shutdown(ctxShutdown)
	}()
	return s.srv.ListenAndServe()
}