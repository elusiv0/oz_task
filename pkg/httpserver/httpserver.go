package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type HttpServer struct {
	shutdownTimeout time.Duration
	server          *http.Server
}

const (
	defaultAddr            = ":80"
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	deafultShutdownTimeout = 3 * time.Second
)

func New(h http.Handler, opts ...Option) *HttpServer {
	s := &http.Server{
		Handler:      h,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		Addr:         defaultAddr,
	}

	httpserver := &HttpServer{
		shutdownTimeout: deafultShutdownTimeout,
		server:          s,
	}

	for _, opt := range opts {
		opt(httpserver)
	}

	return httpserver
}

func (s *HttpServer) Start() error {
	if err := s.server.ListenAndServe(); err != nil {
		return fmt.Errorf("starting http server: %w", err)
	}

	return nil
}

func (s *HttpServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)

	defer cancel()

	return s.server.Shutdown(ctx)
}
