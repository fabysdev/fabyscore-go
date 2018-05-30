package server

// tree contains the http method and the root node.
type tree struct {
	method string
	root   *node
}

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
