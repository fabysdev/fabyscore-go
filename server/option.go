package server

import (
	"net/http"
	"time"
)

// Option is a function type to modify the `http.Server` configuration.
type Option func(*http.Server)

// ReadHeaderTimeout returns an Option for setting the `http.Server`.`ReadHeaderTimeout`
func ReadHeaderTimeout(d time.Duration) Option {
	return func(srv *http.Server) {
		srv.ReadHeaderTimeout = d
	}
}

// IdleTimeout returns an Option for setting the `http.Server`.`IdleTimeout`
func IdleTimeout(d time.Duration) Option {
	return func(srv *http.Server) {
		srv.IdleTimeout = d
	}
}

// WriteTimeout returns an Option for setting the `http.Server`.`WriteTimeout`
func WriteTimeout(d time.Duration) Option {
	return func(srv *http.Server) {
		srv.WriteTimeout = d
	}
}
