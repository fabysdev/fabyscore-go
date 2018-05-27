package fabyscore

// ContextKey is the type of the key for predefined context keys (e.g. request-id).
// const RequestIdCtxKey = &ContextKey{"request-id"}
type ContextKey struct {
	name string
}

// String returns the string representation of the ContextKey.
func (ck *ContextKey) String() string {
	return "fabyscore-go_" + ck.name
}

// dynamicContextKey is the type of the key for dynamic context entries (Unknown string value as key. e.g. /hello/:name -> context entry with key 'name').
type dynamicContextKey interface{}
