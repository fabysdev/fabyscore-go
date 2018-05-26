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

	globalMiddlewares middlewares
}

// NewApp returns an App instance.
func NewApp() *App {
	app := &App{}

	app.router = newRouter()
	app.globalMiddlewares = middlewares{}

	return app
}

// Run starts a http.Server for the application with the given addr.
// This method blocks the calling goroutine.
func (a *App) Run(addr string) {
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

	node.fn(w, req)
}

// GET adds a new request handler for a GET request with the given path.
func (a *App) GET(path string, fn http.HandlerFunc) {
	a.addRoute("GET", path, fn)
}

// POST adds a new request handler for a POST request with the given path.
func (a *App) POST(route string, fn http.HandlerFunc) {
	a.addRoute("POST", route, fn)
}

// PUT adds a new request handler for a PUT request with the given path.
func (a *App) PUT(route string, fn http.HandlerFunc) {
	a.addRoute("PUT", route, fn)
}

// DELETE adds a new request handler for a DELETE request with the given path.
func (a *App) DELETE(route string, fn http.HandlerFunc) {
	a.addRoute("DELETE", route, fn)
}

// PATCH adds a new request handler for a PATCH request with the given path.
func (a *App) PATCH(route string, fn http.HandlerFunc) {
	a.addRoute("PATCH", route, fn)
}

// HEAD adds a new request handler for a HEAD request with the given path.
func (a *App) HEAD(route string, fn http.HandlerFunc) {
	a.addRoute("HEAD", route, fn)
}

// OPTIONS adds a new request handler for a OPTIONS request with the given path.
func (a *App) OPTIONS(route string, fn http.HandlerFunc) {
	a.addRoute("OPTIONS", route, fn)
}

// CONNECT adds a new request handler for a CONNECT request with the given path.
func (a *App) CONNECT(route string, fn http.HandlerFunc) {
	a.addRoute("CONNECT", route, fn)
}

// TRACE adds a new request handler for a TRACE request with the given path.
func (a *App) TRACE(route string, fn http.HandlerFunc) {
	a.addRoute("TRACE", route, fn)
}

// Any adds a route for all HTTP methods.
func (a *App) Any(route string, fn http.HandlerFunc) {
	a.GET(route, fn)
	a.POST(route, fn)
	a.PUT(route, fn)
	a.DELETE(route, fn)
	a.PATCH(route, fn)
	a.HEAD(route, fn)
	a.OPTIONS(route, fn)
	a.CONNECT(route, fn)
	a.TRACE(route, fn)
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

	a.globalMiddlewares = append(a.globalMiddlewares, middleware{
		fn:   fn,
		sort: sorting,
	})

	sort.Sort(a.globalMiddlewares)
}

// UseWithSort @todo
func (a *App) addRoute(method, path string, fn http.HandlerFunc) {
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
