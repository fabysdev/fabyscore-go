package fabyscore

import (
	"net/http"
	"time"
)

// ServerOption is a function type to modifiy the `http.Server` configuration.
type ServerOption func(*http.Server)

// ServerReadHeaderTimeout returns a ServerOption for setting the `http.Server`.`ReadHeaderTimeout`
func ServerReadHeaderTimeout(d time.Duration) ServerOption {
	return func(srv *http.Server) {
		srv.ReadHeaderTimeout = d
	}
}

// ServerIdleTimeout returns a ServerOption for setting the `http.Server`.`IdleTimeout`
func ServerIdleTimeout(d time.Duration) ServerOption {
	return func(srv *http.Server) {
		srv.IdleTimeout = d
	}
}
