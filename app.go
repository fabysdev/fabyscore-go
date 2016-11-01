package fabyscore

import (
    "net/http"
    "sort"
    "strings"
)

//----------------------------------------------------------------------------------------------------------------------

// HandlerFunc is the fabyscore http.HandlerFunc type.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request)



//----------------------------------------------------------------------------------------------------------------------

// App is the main fabyscore instance, it contains the modifiers and the router.
// Create a new instance by using NewApp().
type App struct {
	router *router
	modifiers Modifiers
    beforeModifiers Modifiers
    notFoundHandler HandlerFunc
}

// NewApp returns an App instance.
func NewApp() *App {
	app := &App{}
	
    app.router = newRouter()
    app.modifiers = Modifiers{}
    app.beforeModifiers = Modifiers{}
    
	return app
}

// Run starts a http.Server for the application with the given addr. 
// This method blocks the calling goroutine.
func(a *App) Run(addr string) {
    sort.Sort(a.modifiers)
    sort.Sort(a.beforeModifiers)
    
    http.ListenAndServe(addr, a)
}

// See http.Handler interface's ServeHTTP.
func (a *App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    for _, mod := range a.beforeModifiers {
        w, req = mod.fn(w, req)
        if req == nil {
            return
        }
    }
    
    node, req := a.router.resolve(req)
    if node == nil || node.fn == nil {
        if a.notFoundHandler != nil {
            a.notFoundHandler(w, req)
        }else {
            http.NotFound(w, req)
        }
        
        return
    }
 
    for _, mod := range a.modifiers {
        w, req = mod.fn(w, req)
        if req == nil {
            return
        }
    }

    for _, mod := range node.modifiers {
        w, req = mod.fn(w, req)
        if req == nil {
            return
        }
    }

    node.fn(w, req)
}

// AddModifier adds a modifier which is executed after the route is resolved.
func(a *App) AddModifier(sort int, fn HandlerFunc) {
    a.modifiers = append(a.modifiers, Modifier{
        sort: sort,
        fn: fn,
    })
}

// AddBeforeModifier adds a modifier which is executed before the route is resolved.
func(a *App) AddBeforeModifier(sort int, fn HandlerFunc) {
    a.beforeModifiers = append(a.beforeModifiers, Modifier{
        sort: sort,
        fn: fn,
    })
}

// GET adds a new request handler (HandlerFunc and modifiers) for a GET request with the given path.
func(a *App) GET(path string, fn HandlerFunc, modifiers ...Modifier) {
    a.router.addRoute("GET", path, fn, modifiers)
}

// POST adds a new request handler (HandlerFunc and modifiers) for a POST request with the given path.
func(a *App) POST(route string, fn HandlerFunc, modifiers ...Modifier) {
    a.router.addRoute("POST", route, fn, modifiers)
}

// PUT adds a new request handler (HandlerFunc and modifiers) for a PUT request with the given path.
func(a *App) PUT(route string, fn HandlerFunc, modifiers ...Modifier) {
    a.router.addRoute("PUT", route, fn, modifiers)
}

// DELETE adds a new request handler (HandlerFunc and modifiers) for a DELETE request with the given path.
func(a *App) DELETE(route string, fn HandlerFunc, modifiers ...Modifier) {
    a.router.addRoute("DELETE", route, fn, modifiers)
}

// PATCH adds a new request handler (HandlerFunc and modifiers) for a PATCH request with the given path.
func(a *App) PATCH(route string, fn HandlerFunc, modifiers ...Modifier) {
    a.router.addRoute("PATCH", route, fn, modifiers)
}

// HEAD adds a new request handler (HandlerFunc and modifiers) for a HEAD request with the given path.
func(a *App) HEAD(route string, fn HandlerFunc, modifiers ...Modifier) {
    a.router.addRoute("HEAD", route, fn, modifiers)
}

// OPTIONS adds a new request handler (HandlerFunc and modifiers) for a OPTIONS request with the given path.
func(a *App) OPTIONS(route string, fn HandlerFunc, modifiers ...Modifier) {
    a.router.addRoute("OPTIONS", route, fn, modifiers)
}

// CONNECT adds a new request handler (HandlerFunc and modifiers) for a CONNECT request with the given path.
func(a *App) CONNECT(route string, fn HandlerFunc, modifiers ...Modifier) {
    a.router.addRoute("CONNECT", route, fn, modifiers)
}

// TRACE adds a new request handler (HandlerFunc and modifiers) for a TRACE request with the given path.
func(a *App) TRACE(route string, fn HandlerFunc, modifiers ...Modifier) {
    a.router.addRoute("TRACE", route, fn, modifiers)
}

// Any adds a route for all HTTP methods.
func(a *App) Any(route string, fn HandlerFunc, modifiers ...Modifier) {
    a.GET(route, fn, modifiers...)
    a.POST(route, fn, modifiers...)
    a.PUT(route, fn, modifiers...)
    a.DELETE(route, fn, modifiers...)
    a.PATCH(route, fn, modifiers...)
    a.HEAD(route, fn, modifiers...)
    a.OPTIONS(route, fn, modifiers...)
    a.CONNECT(route, fn, modifiers...)
    a.TRACE(route, fn, modifiers...)
}

// Group adds multiple routes with common modifiers and common path prefix.
func(a *App) Group(path string, modifiers Modifiers, routes ...*Route) {
    for _, route := range routes {
        route.Path = "/" + strings.Trim(path, "/") + "/" + strings.TrimLeft(route.Path, "/")
        route.Modifiers = append(route.Modifiers, modifiers...)

        a.router.addRoute(route.Method, route.Path, route.Fn, route.Modifiers)
    }
}

// SetNotFoundHandler sets the HandlerFunc executed if no handler is found for the request.
func(a *App) SetNotFoundHandler(fn HandlerFunc) {
    a.notFoundHandler = fn
}



//----------------------------------------------------------------------------------------------------------------------

// Modifier can change the request and the response before the route handler is called.
// The execution path is linear: Modifier1 -> Modifier2 -> RouteHandler.
// Create a new instance by using NewModifier().
type Modifier struct {
    sort int;
    fn   HandlerFunc
}

// NewModifier returns a new Modifier instance.
func NewModifier(sort int, fn HandlerFunc) Modifier {
    return Modifier{
        sort: sort,
        fn: fn,
    }
}



//----------------------------------------------------------------------------------------------------------------------

// Modifiers is a slice of Modifier implementing sort.Interface.
type Modifiers []Modifier

// See sort.Interface Len().
func (slice Modifiers) Len() int {
    return len(slice)
}

// See sort.Interface Less().
func (slice Modifiers) Less(i, j int) bool {
    return slice[i].sort < slice[j].sort;
}

// See sort.Interface Swap().
func (slice Modifiers) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}


