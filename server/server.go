package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"
)

// Server is the main instance.
// Create a new instance by using New().
type Server struct {
	router          *router
	notFoundHandler http.HandlerFunc
	middlewares     middlewares

	quit chan os.Signal
}

// New returns an Server instance.
func New() *Server {
	s := &Server{}

	s.router = newRouter()
	s.middlewares = middlewares{}

	s.quit = make(chan os.Signal, 1)
	signal.Notify(s.quit, os.Interrupt, syscall.SIGTERM)

	return s
}

// Run starts a http.Server for the application with the given addr.
// This method blocks the calling goroutine.
func (s *Server) Run(addr string, options ...Option) error {
	return s.run(addr, options, "", "")
}

// RunTLS starts a https http.Server for the application with the given addr and certificate files.
// This method blocks the calling goroutine.
func (s *Server) RunTLS(addr, certFile, keyFile string, options ...Option) error {
	return s.run(addr, options, certFile, keyFile)
}

// See http.Handler interface's ServeHTTP.
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	node, req := s.router.resolve(req)
	if node == nil || node.fn == nil {
		if s.notFoundHandler != nil {
			s.notFoundHandler(w, req)
		} else {
			http.NotFound(w, req)
		}

		return
	}

	node.fn.ServeHTTP(w, req)
}

// GET adds a new request handler for a GET request with the given path.
func (s *Server) GET(path string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	s.addRoute("GET", path, fn, middlewares)
}

// POST adds a new request handler for a POST request with the given path.
func (s *Server) POST(route string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	s.addRoute("POST", route, fn, middlewares)
}

// PUT adds a new request handler for a PUT request with the given path.
func (s *Server) PUT(route string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	s.addRoute("PUT", route, fn, middlewares)
}

// DELETE adds a new request handler for a DELETE request with the given path.
func (s *Server) DELETE(route string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	s.addRoute("DELETE", route, fn, middlewares)
}

// PATCH adds a new request handler for a PATCH request with the given path.
func (s *Server) PATCH(route string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	s.addRoute("PATCH", route, fn, middlewares)
}

// HEAD adds a new request handler for a HEAD request with the given path.
func (s *Server) HEAD(route string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	s.addRoute("HEAD", route, fn, middlewares)
}

// OPTIONS adds a new request handler for a OPTIONS request with the given path.
func (s *Server) OPTIONS(route string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	s.addRoute("OPTIONS", route, fn, middlewares)
}

// CONNECT adds a new request handler for a CONNECT request with the given path.
func (s *Server) CONNECT(route string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	s.addRoute("CONNECT", route, fn, middlewares)
}

// TRACE adds a new request handler for a TRACE request with the given path.
func (s *Server) TRACE(route string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	s.addRoute("TRACE", route, fn, middlewares)
}

// Any adds a route for all available methods.
func (s *Server) Any(route string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	s.GET(route, fn, middlewares...)
	s.POST(route, fn, middlewares...)
	s.PUT(route, fn, middlewares...)
	s.DELETE(route, fn, middlewares...)
	s.PATCH(route, fn, middlewares...)
	s.HEAD(route, fn, middlewares...)
	s.OPTIONS(route, fn, middlewares...)
	s.CONNECT(route, fn, middlewares...)
	s.TRACE(route, fn, middlewares...)
}

// Group adds multiple routes with a common path prefix.
func (s *Server) Group(path string, fn GroupSetupFunc) {
	basePath := strings.Trim(path, "/")
	if basePath != "" {
		basePath = "/" + basePath
	}

	group := &Group{
		basePath:    basePath,
		srv:         s,
		middlewares: middlewares{},
	}

	fn(group)
}

// SetNotFoundHandler sets the http.HandlerFunc executed if no handler is found for the request.
func (s *Server) SetNotFoundHandler(fn http.HandlerFunc) {
	s.notFoundHandler = fn
}

// Use adds an middleware on server level.
// Defaults to a sorting of 0. Use `UseWithSort` to set an sorting for a middleware.
func (s *Server) Use(fn MiddlewareFunc) {
	s.UseWithSorting(fn, 0)
}

// UseWithSorting adds an middleware with an custom sorting value on server level.
func (s *Server) UseWithSorting(fn MiddlewareFunc, sorting int) {
	if s.router.hasRoutes {
		panic("Server middlewares must be defined before the routes")
	}

	s.middlewares = append(s.middlewares, middleware{
		fn:      fn,
		sorting: sorting,
	})

	sort.Sort(s.middlewares)
}

// ServeFiles serves the files from the given root at the given path.
// The given path is converted into a match-all path (e.g. /static/ => /static/*file)
// The default http.NotFound is used for 404s.
// Will not serve the directory, only files.
func (s *Server) ServeFiles(path string, root http.FileSystem) {
	fileServer := http.FileServer(FileSystem{root})

	s.GET(strings.TrimSuffix(path, "/")+"/*file", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		file := ctx.Value("file")

		r.URL.Path = file.(string)

		fileServer.ServeHTTP(w, r)
	})
}

// run starts and creates the http.Server and does the graceful shutdown.
func (s *Server) run(addr string, options []Option, certFile, keyFile string) error {
	// unset middlewares, they are only used during setup to create the final handler functions
	s.middlewares = nil

	// create http.Server
	srv := &http.Server{
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	for _, option := range options {
		option(srv)
	}

	srv.Addr = addr
	srv.Handler = s

	// graceful shutdown
	done := make(chan bool)
	go func() {
		<-s.quit

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		srv.Shutdown(ctx)

		close(done)
	}()

	var err error
	if certFile != "" && keyFile != "" {
		err = srv.ListenAndServeTLS(certFile, keyFile)
	} else {
		err = srv.ListenAndServe()
	}

	if err != http.ErrServerClosed {
		close(done)
		return err
	}

	<-done
	return nil
}

// addRoute adds a route to the router with the middleware aware handler.
func (s *Server) addRoute(method, path string, fn http.Handler, middlewares []MiddlewareFunc) {
	// create handler with route middlewares
	middlewaresLen := len(middlewares)
	if middlewaresLen > 0 {
		for i := middlewaresLen - 1; i >= 0; i-- {
			fn = middlewares[i](fn)
		}
	}

	// create handler function with server middlewares
	middlewaresLen = len(s.middlewares)
	if middlewaresLen > 0 {
		for i := middlewaresLen - 1; i >= 0; i-- {
			fn = s.middlewares[i].fn(fn)
		}
	}

	// add route to router
	s.router.addRoute(method, path, fn)
}
