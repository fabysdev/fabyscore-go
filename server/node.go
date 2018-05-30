package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// node is a tree node for a specific path part.
type node struct {
	path      string
	children  []*node
	isDynamic bool
	fn        http.Handler
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

	path := req.URL.Path
	pathLen := len(path)

	startIndex := 1
	var ctx context.Context
	for i := 1; i < pathLen; i++ {
		if path[i] != '/' {
			continue
		}

		n = n.load(path[startIndex:i])
		if n == nil {
			return nil, nil
		}

		if n.isDynamic {
			if ctx == nil {
				ctx = req.Context()
			}

			ctx = context.WithValue(ctx, dynamicContextKey(n.path[1:]), path[startIndex:i])
		}

		startIndex = i + 1
	}

	if startIndex != pathLen {
		n = n.load(path[startIndex:])
		if n == nil {
			return nil, nil
		}

		if n.isDynamic {
			if ctx == nil {
				ctx = req.Context()
			}

			ctx = context.WithValue(ctx, dynamicContextKey(n.path[1:]), path[startIndex:])
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
