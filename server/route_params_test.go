package server

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouteParams(t *testing.T) {
	params := newRouteParams()
	params.Add("test", "value")

	assert.Equal(t, params.Len(), 1)
	assert.Equal(t, params.Get("test"), "value")
	assert.Equal(t, params.Get("notfound"), "")

	params.Reset()
	assert.Equal(t, params.Len(), 0)
}

func TestParam(t *testing.T) {
	params := newRouteParams()
	params.Add("test", "value")

	req, _ := http.NewRequest("GET", "/", nil)
	req = req.WithContext(context.WithValue(req.Context(), routeParamsContextKey, params))

	assert.Equal(t, Param(req, "test"), "value")
	assert.Equal(t, Param(req, "notfound"), "")
}
