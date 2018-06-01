package server

// ContextKey is the type of the key for predefined context keys (e.g. request-id).
// const RequestIdCtxKey = &ContextKey{"request-id"}
type ContextKey struct {
	Name string
}

// String returns the string representation of the ContextKey.
func (ck *ContextKey) String() string {
	return "fabyscore-go_" + ck.Name
}

// dynamicContextKey is the type of the key for dynamic context entries (Unknown string value as key. e.g. /hello/:name -> context entry with key 'name').
type dynamicContextKey interface{}
