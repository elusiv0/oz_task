package httpserver

import (
	"time"
)

type Option func(s *HttpServer)

func ReadTimeout(t time.Duration) Option {
	return func(s *HttpServer) {
		s.server.ReadTimeout = t
	}
}

func Port(p string) Option {
	return func(s *HttpServer) {
		s.server.Addr = p
	}
}

func WriteTimeout(t time.Duration) Option {
	return func(s *HttpServer) {
		s.server.WriteTimeout = t
	}
}

func ShutdownTimeout(t time.Duration) Option {
	return func(s *HttpServer) {
		s.shutdownTimeout = t
	}
}
