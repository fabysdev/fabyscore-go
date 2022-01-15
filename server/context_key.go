package server

// ContextKey is the type of the key for predefined context keys (e.g. request-id).
// var RequestIDContextKey = &ContextKey{"request-id"}
type ContextKey struct {
	Name string
}

// String returns the string representation of the ContextKey.
func (ck *ContextKey) String() string {
	return "fabyscore-go_" + ck.Name
}
