package fabyscore

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

//----------------------------------------------------------------------------------------------------------------------

// Route is the route definition.
type Route struct {
	Fn     HandlerFunc
	Path   string
	Method string
}

// dynamicContextKey is the type of the key for dynamic context entries (Unknown string as key).
type dynamicContextKey interface{}

// contextKey is the type of the key for predefined context keys (e.g. request-id).
// const RequestIdCtxKey = &contextKey{"request-id"}
type contextKey struct {
	name string
}

func (ck *contextKey) String() string {
	return "fabyscore-go_" + ck.name
}

//----------------------------------------------------------------------------------------------------------------------

// router is a http request router.
// Create a new instance by using newRouter().
type router struct {
	trees     methodTrees
	hasRoutes bool
}

// newRouter returns a router instance.
func newRouter() *router {
	r := &router{
		trees: make(methodTrees, 0, 9),
	}

	return r
}

// addRoute adds a new request handler for a given method/path combination.
func (r *router) addRoute(method string, path string, fn HandlerFunc) {
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
func (r *router) resolve(req *http.Request) (*node, *http.Request) {
	root := r.trees.getRoot(req.Method)
	if root == nil {
		return nil, nil
	}

	return root.resolve(req)
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

//----------------------------------------------------------------------------------------------------------------------

// node is a tree node for a specific path part.
type node struct {
	path      string
	children  []*node
	isDynamic bool
	fn        HandlerFunc
}

// add adds a new node with a given path.
func (n *node) add(path string, fn HandlerFunc) {
	if path == "/" {
		n.path = "/"
		n.fn = fn
		return
	}

	parts := strings.Split(path, "/")[1:]

	var resolvedNode *node
	for i := 0; i < len(parts); i++ {
		part := parts[i]
		resolvedNode = n.load(part)
		if resolvedNode == nil {
			resolvedNode = &node{
				path: part,
			}

			if len(part) > 0 && part[0] == ':' {
				resolvedNode.isDynamic = true
			}

			n.children = append(n.children, resolvedNode)
		}

		n = resolvedNode
	}

	resolvedNode.fn = fn
}

// resolve returns the node and the request with context for a given request.
// Returns nil, nil if no node was found for the request.
func (n *node) resolve(req *http.Request) (*node, *http.Request) {
	if req.URL.Path == "/" {
		return n, req
	}

	parts := strings.Split(req.URL.Path, "/")[1:]

	var ctx context.Context
	for i := 0; i < len(parts); i++ {
		part := parts[i]
		n = n.load(part)
		if n == nil {
			return nil, nil
		}

		if n.isDynamic {
			if ctx == nil {
				ctx = req.Context()
			}

			ctx = context.WithValue(ctx, dynamicContextKey(n.path[1:]), part)
		}
	}

	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return n, req
}

// load returns the child node matching the given path, nil if no matching node was found.
func (n *node) load(path string) *node {
	for _, node := range n.children {
		if node.path == path || node.isDynamic {
			return node
		}
	}

	return nil
}

// dump returns the node and its children as string.
func (n *node) dump(prefix string) string {
	line := fmt.Sprintf("%s%s\n", prefix, n.path)
	prefix += "  "
	for _, node := range n.children {
		line += node.dump(prefix)
	}

	return line
}

//----------------------------------------------------------------------------------------------------------------------

// tree contains the http method and the root node.
type tree struct {
	method string
	root   *node
}

//----------------------------------------------------------------------------------------------------------------------

// methodTrees is a slice of trees.
type methodTrees []*tree

// getRoot returns the root node for the tree with the given method, nil if no tree for the given method exists.
func (mt methodTrees) getRoot(method string) *node {
	for _, t := range mt {
		if t.method == method {
			return t.root
		}
	}

	return nil
}
