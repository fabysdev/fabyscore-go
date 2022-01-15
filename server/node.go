package server

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
)

// node is a tree node for a specific path part.
type node struct {
	path       string
	children   []*node
	parent     *node
	isDynamic  bool
	isMatchAll bool
	fn         http.Handler
}

// add adds a new node with a given path.
func (n *node) add(path string, fn http.Handler) {
	if path == "/" {
		n.path = "/"
		n.fn = fn
		return
	}

	parts := strings.Split(path, "/")[1:]
	if parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}

	var resolvedNode *node
	for i := 0; i < len(parts); i++ {
		part := parts[i]

		resolvedNode = n.load(part)

		// check if the route can be added to tree
		if resolvedNode != nil {
			// the route can not be added if a match all route already exists for this part
			if resolvedNode.isMatchAll && (part == "" || part[0] != '*') {
				panic(fmt.Sprintf("Route '%s' can not be added. Match-All route '%s' conflicts with it. Check the route registration order.", path, resolvedNode.resolvePath()))
			}

			// the route can not be added if a dynamic route with a different path part already exists for this part
			if resolvedNode.isDynamic && part != resolvedNode.path {
				panic(fmt.Sprintf("Route '%s' can not be added. Dynamic route '%s' conflicts with it. Check the route registration order.", path, resolvedNode.resolvePath()))
			}
		}

		if resolvedNode == nil {
			resolvedNode = &node{
				path:   part,
				parent: n,
			}

			// dynamic node
			if len(part) > 0 && part[0] == ':' {
				resolvedNode.isDynamic = true
			}

			// match all node
			if len(part) > 0 && part[0] == '*' {
				resolvedNode.isMatchAll = true
			}

			n.children = append(n.children, resolvedNode)
		}

		// break if current node is a match-all
		if resolvedNode.isMatchAll {
			if i != len(parts)-1 {
				panic(fmt.Sprintf("Route '%s' has ineffective parts. Everything after the Match-All part '%s' is ignored. Remove the ineffective parts from the route.", path, resolvedNode.resolvePath()))
			}

			break
		}

		n = resolvedNode
	}

	resolvedNode.fn = fn
}

// resolve returns the node and the request with context for a given request.
// Returns nil, nil if no node was found for the request.
func (n *node) resolve(req *http.Request, paramsPool *sync.Pool) (*node, *http.Request, *routeParams) {
	if req.URL.Path == "/" {
		if n.fn == nil && len(n.children) != 0 {
			// check if a dynamic or match all route exists as root child
			for _, child := range n.children {
				if child.isDynamic || child.isMatchAll {
					return child, req, nil
				}
			}
		}

		return n, req, nil
	}

	path := req.URL.Path
	pathLen := len(path)
	var params *routeParams
	startIndex := 1
	for i := 1; i < pathLen; i++ {
		// skip until next /
		if path[i] != '/' {
			continue
		}

		// load the next node with the current path part
		n = n.load(path[startIndex:i])
		if n == nil {
			return nil, req, nil
		}

		// if the node is a match-all, add the remaining path as param and return the node
		if n.isMatchAll {
			if params == nil {
				params = paramsPool.Get().(*routeParams)
			}

			params.Add(n.path[1:], path[startIndex:])
			req = req.WithContext(context.WithValue(req.Context(), routeParamsContextKey, params))
			return n, req, params
		}

		// add the part as param if the node is dynamic
		if n.isDynamic {
			if params == nil {
				params = paramsPool.Get().(*routeParams)
			}

			params.Add(n.path[1:], path[startIndex:i])
		}

		startIndex = i + 1
	}

	// if the path does not end with an / startIndex won't be at the end of the path (because the loop skips to /)
	if startIndex != pathLen {
		n = n.load(path[startIndex:])
		if n == nil {
			return nil, req, params
		}

		if n.isDynamic || n.isMatchAll {
			if params == nil {
				params = paramsPool.Get().(*routeParams)
			}

			params.Add(n.path[1:], path[startIndex:])
		}
	}

	if params != nil {
		req = req.WithContext(context.WithValue(req.Context(), routeParamsContextKey, params))
	}

	if n.fn == nil && len(n.children) == 1 && (n.children[0].isDynamic || n.children[0].isMatchAll) {
		return n.children[0], req, params
	}

	return n, req, params
}

// load returns the child node matching the given path, nil if no matching node was found.
func (n *node) load(path string) *node {
	for _, node := range n.children {
		if node.path == path || node.isDynamic || node.isMatchAll {
			return node
		}
	}

	return nil
}

// resolvePath returns the full path to the node.
func (n *node) resolvePath() string {
	path := []string{
		n.path,
	}

	for n.parent != nil {
		n = n.parent
		if n.path != "/" {
			path = append(path, n.path)
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(path)))

	return "/" + strings.Join(path, "/")
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
