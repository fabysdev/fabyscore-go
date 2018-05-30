package server

import "net/http"

// MiddlewareFunc is the type for middleware functions.
type MiddlewareFunc func(http.Handler) http.Handler

// middleware is the type for an sort aware middleware.
type middleware struct {
	fn      MiddlewareFunc
	sorting int
}

// middlewares is a sortable middleware slice implementing sort.Interface.
type middlewares []middleware

// See sort.Interface Len().
func (slice middlewares) Len() int {
	return len(slice)
}

// See sort.Interface Less().
func (slice middlewares) Less(i, j int) bool {
	return slice[i].sorting < slice[j].sorting
}

// See sort.Interface Swap().
func (slice middlewares) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
