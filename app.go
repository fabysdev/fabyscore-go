package fabyscore

import (
	"net/http"
	"sort"
	"strings"
)

//----------------------------------------------------------------------------------------------------------------------

// MiddlewareFunc @todo
type MiddlewareFunc func(http.Handler) http.Handler

// @todo, comment?
type middleware struct {
	fn   MiddlewareFunc
	sort int
}

// @todo, comment
type middlewares []middleware

// See sort.Interface Len().
func (slice middlewares) Len() int {
	return len(slice)
}

// See sort.Interface Less().
func (slice middlewares) Less(i, j int) bool {
	return slice[i].sort < slice[j].sort
}

// See sort.Interface Swap().
func (slice middlewares) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

//----------------------------------------------------------------------------------------------------------------------

// App is the main fabyscore instance.
// Create a new instance by using NewApp().
type App struct {
	router          *router
	notFoundHandler http.HandlerFunc

	middlewares middlewares
}

// NewApp returns an App instance.
func NewApp() *App {
	app := &App{}

	app.router = newRouter()
	app.middlewares = middlewares{}

	return app
}

// Run starts a http.Server for the application with the given addr.
// This method blocks the calling goroutine.
func (a *App) Run(addr string) {
	a.middlewares = nil

	http.ListenAndServe(addr, a)
}

// See http.Handler interface's ServeHTTP.
func (a *App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	node, req := a.router.resolve(req)
	if node == nil || node.fn == nil {
		if a.notFoundHandler != nil {
			a.notFoundHandler(w, req)
		} else {
			http.NotFound(w, req)
		}

		return
	}

	node.fn.ServeHTTP(w, req)
}

// GET adds a new request handler for a GET request with the given path.
func (a *App) GET(path string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	a.addRoute("GET", path, fn, middlewares)
}

// POST adds a new request handler for a POST request with the given path.
func (a *App) POST(route string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	a.addRoute("POST", route, fn, middlewares)
}

// PUT adds a new request handler for a PUT request with the given path.
func (a *App) PUT(route string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	a.addRoute("PUT", route, fn, middlewares)
}

// DELETE adds a new request handler for a DELETE request with the given path.
func (a *App) DELETE(route string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	a.addRoute("DELETE", route, fn, middlewares)
}

// PATCH adds a new request handler for a PATCH request with the given path.
func (a *App) PATCH(route string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	a.addRoute("PATCH", route, fn, middlewares)
}

// HEAD adds a new request handler for a HEAD request with the given path.
func (a *App) HEAD(route string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	a.addRoute("HEAD", route, fn, middlewares)
}

// OPTIONS adds a new request handler for a OPTIONS request with the given path.
func (a *App) OPTIONS(route string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	a.addRoute("OPTIONS", route, fn, middlewares)
}

// CONNECT adds a new request handler for a CONNECT request with the given path.
func (a *App) CONNECT(route string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	a.addRoute("CONNECT", route, fn, middlewares)
}

// TRACE adds a new request handler for a TRACE request with the given path.
func (a *App) TRACE(route string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	a.addRoute("TRACE", route, fn, middlewares)
}

// Any adds a route for all HTTP methods.
func (a *App) Any(route string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	a.GET(route, fn, middlewares...)
	a.POST(route, fn, middlewares...)
	a.PUT(route, fn, middlewares...)
	a.DELETE(route, fn, middlewares...)
	a.PATCH(route, fn, middlewares...)
	a.HEAD(route, fn, middlewares...)
	a.OPTIONS(route, fn, middlewares...)
	a.CONNECT(route, fn, middlewares...)
	a.TRACE(route, fn, middlewares...)
}

// Group adds multiple routes with common path prefix.
func (a *App) Group(path string, routes ...*Route) {
	for _, route := range routes {
		route.Path = "/" + strings.Trim(path, "/") + "/" + strings.TrimLeft(route.Path, "/")

		a.router.addRoute(route.Method, route.Path, route.Fn)
	}
}

// SetNotFoundHandler sets the http.HandlerFunc executed if no handler is found for the request.
func (a *App) SetNotFoundHandler(fn http.HandlerFunc) {
	a.notFoundHandler = fn
}

// Use @todo
// Defaults to a sort of 0. Use `UseWithSort` to set an sort for a middleware.
func (a *App) Use(fn MiddlewareFunc) {
	a.UseWithSort(fn, 0)
}

// UseWithSort @todo
func (a *App) UseWithSort(fn MiddlewareFunc, sorting int) {
	if a.router.hasRoutes {
		panic("App middlewares must be defined before the routes")
	}

	a.middlewares = append(a.middlewares, middleware{
		fn:   fn,
		sort: sorting,
	})

	sort.Sort(a.middlewares)
}

// UseWithSort @todo
func (a *App) addRoute(method, path string, fn http.Handler, middlewares []MiddlewareFunc) {
	// create handler with route middlewares
	middlewaresLen := len(middlewares)
	if middlewaresLen > 0 {
		for i := middlewaresLen - 1; i >= 0; i-- {
			fn = middlewares[i](fn)
		}
	}

	// create handler function with app middlewares
	middlewaresLen = len(a.middlewares)
	if middlewaresLen > 0 {
		for i := middlewaresLen - 1; i >= 0; i-- {
			fn = a.middlewares[i].fn(fn)
		}
	}

	// add route to router
	a.router.addRoute(method, path, fn)
}

//----------------------------------------------------------------------------------------------------------------------

// Modifier can change the request and the response before the route handler is called.
// The execution path is linear: Modifier1 -> Modifier2 -> RouteHandler.
// Create a new instance by using NewModifier().
type Modifier struct {
	sort int
	fn   http.HandlerFunc
}

// NewModifier returns a new Modifier instance.
func NewModifier(sort int, fn http.HandlerFunc) Modifier {
	return Modifier{
		sort: sort,
		fn:   fn,
	}
}
