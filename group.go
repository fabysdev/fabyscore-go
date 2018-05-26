package fabyscore

import (
	"net/http"
	"sort"
	"strings"
)

// GroupSetupFunc @todo
type GroupSetupFunc func(g *Group)

// Group @todo
type Group struct {
	basePath    string
	middlewares middlewares
	app         *App

	hasRoutes bool
}

// Use @todo
// Defaults to a sort of 0. Use `UseWithSort` to set an sort for a middleware.
func (g *Group) Use(fn MiddlewareFunc) {
	g.UseWithSort(fn, 0)
}

// UseWithSort @todo
func (g *Group) UseWithSort(fn MiddlewareFunc, sorting int) {
	if g.hasRoutes {
		panic("Group middlewares must be defined before the routes")
	}

	g.middlewares = append(g.middlewares, middleware{
		fn:   fn,
		sort: sorting,
	})

	sort.Sort(g.middlewares)
}

// GET adds a new request handler for a GET request with the given path.
func (g *Group) GET(path string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	g.addRoute("GET", path, fn, middlewares)
}

// POST adds a new request handler for a POST request with the given path.
func (g *Group) POST(route string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	g.addRoute("POST", route, fn, middlewares)
}

// PUT adds a new request handler for a PUT request with the given path.
func (g *Group) PUT(route string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	g.addRoute("PUT", route, fn, middlewares)
}

// DELETE adds a new request handler for a DELETE request with the given path.
func (g *Group) DELETE(route string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	g.addRoute("DELETE", route, fn, middlewares)
}

// PATCH adds a new request handler for a PATCH request with the given path.
func (g *Group) PATCH(route string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	g.addRoute("PATCH", route, fn, middlewares)
}

// HEAD adds a new request handler for a HEAD request with the given path.
func (g *Group) HEAD(route string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	g.addRoute("HEAD", route, fn, middlewares)
}

// OPTIONS adds a new request handler for a OPTIONS request with the given path.
func (g *Group) OPTIONS(route string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	g.addRoute("OPTIONS", route, fn, middlewares)
}

// CONNECT adds a new request handler for a CONNECT request with the given path.
func (g *Group) CONNECT(route string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	g.addRoute("CONNECT", route, fn, middlewares)
}

// TRACE adds a new request handler for a TRACE request with the given path.
func (g *Group) TRACE(route string, fn http.HandlerFunc, middlewares ...MiddlewareFunc) {
	g.addRoute("TRACE", route, fn, middlewares)
}

// addRoute @todo
func (g *Group) addRoute(method, path string, fn http.Handler, middlewares []MiddlewareFunc) {
	groupRouteMiddlewares := []MiddlewareFunc{}
	for _, middleware := range g.middlewares {
		groupRouteMiddlewares = append(groupRouteMiddlewares, middleware.fn)
	}

	groupRouteMiddlewares = append(groupRouteMiddlewares, middlewares...)

	g.app.addRoute(method, g.basePath+"/"+strings.TrimLeft(path, "/"), fn, groupRouteMiddlewares)

	g.hasRoutes = true
}
