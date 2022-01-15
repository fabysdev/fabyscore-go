package server

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// Route is the route definition.
type Route struct {
	Fn     http.Handler
	Path   string
	Method string
}

// router is a http request router.
// Create a new instance by using newRouter().
type router struct {
	trees     methodTrees
	hasRoutes bool
	pool      *sync.Pool
}

// newRouter returns a router instance.
func newRouter() *router {
	r := &router{
		trees: make(methodTrees, 0, 9),
		pool:  &sync.Pool{},
	}

	r.pool.New = func() interface{} {
		return newRouteParams()
	}

	return r
}

// addRoute adds a new request handler for a given method/path combination.
func (r *router) addRoute(method string, path string, fn http.Handler) {
	method = strings.ToUpper(method)

	root := r.trees.getRoot(method)
	if root == nil {
		root = &node{path: "/"}

		t := &tree{
			method: method,
			root:   root,
		}

		r.trees = append(r.trees, t)
	}

	root.add(path, fn)

	r.hasRoutes = true
}

// resolve returns the tree node and the request containing the context(if the route has parameters) for a given request.
// Returns nil,nil if no node was found for the request.
func (r *router) resolve(req *http.Request) (*node, *http.Request, *routeParams) {
	root := r.trees.getRoot(req.Method)
	if root == nil {
		return nil, req, nil
	}

	return root.resolve(req, r.pool)
}

// resetParams resets the params object and adds it back to the pool.
func (r *router) resetParams(params *routeParams) {
	if params == nil {
		return
	}

	params.Reset()
	r.pool.Put(params)
}

// dumpTree returns all trees as a string.
func (r *router) dumpTree() string {
	var str string
	for _, t := range r.trees {
		str += fmt.Sprintf("%s:\n", t.method)
		str += t.root.dump("")
		str += "\n\n"
	}

	return str
}
