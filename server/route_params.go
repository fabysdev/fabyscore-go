package server

import "net/http"

// routeParamsContextKey context key for the params object.
var routeParamsContextKey = &ContextKey{"route-params"}

// routeParams holds the dynamic / match all params from the url path.
type routeParams struct {
	paramsKeys, paramsValues []string
}

// newRouteParams creates a new routeParams object.
func newRouteParams() *routeParams {
	return &routeParams{}
}

// Reset resets the params.
func (rp *routeParams) Reset() {
	rp.paramsKeys = rp.paramsKeys[:0]
	rp.paramsValues = rp.paramsValues[:0]
}

// Add adds a new param.
func (rp *routeParams) Add(key, value string) {
	rp.paramsKeys = append(rp.paramsKeys, key)
	rp.paramsValues = append(rp.paramsValues, value)
}

// Len returns count of params.
func (rp *routeParams) Len() int {
	return len(rp.paramsKeys)
}

// Get returns the corresponding value for the param name or an empty string.
func (rp *routeParams) Get(name string) string {
	for k := len(rp.paramsKeys) - 1; k >= 0; k-- {
		if rp.paramsKeys[k] == name {
			return rp.paramsValues[k]
		}
	}

	return ""
}

// Param returns the corresponding value for the param name or an empty string.
func Param(r *http.Request, name string) string {
	if params := r.Context().Value(routeParamsContextKey); params != nil {
		return params.(*routeParams).Get(name)
	}

	return ""
}
