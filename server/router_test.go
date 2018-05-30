package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddRouteSimple(t *testing.T) {
	router := newRouter()
	router.addRoute("GET", "/testroute", http.HandlerFunc(simpleHandler))

	assert.Equal(t, "GET:\n/\n  testroute\n\n\n", router.dumpTree())

	router.addRoute("POST", "/route", http.HandlerFunc(simpleHandler))
	router.addRoute("GET", "/route", http.HandlerFunc(simpleHandler))

	assert.Equal(t, "GET:\n/\n  testroute\n  route\n\n\nPOST:\n/\n  route\n\n\n", router.dumpTree())
}

func TestAddRouteSimpleRoot(t *testing.T) {
	router := newRouter()
	router.addRoute("GET", "/", http.HandlerFunc(simpleHandler))
	router.addRoute("GET", "/testroute", http.HandlerFunc(simpleHandler))

	assert.Equal(t, "GET:\n/\n  testroute\n\n\n", router.dumpTree())
}

func TestAddRouteSimpleTree(t *testing.T) {
	router := newRouter()

	router.addRoute("GET", "/test/route/simple", http.HandlerFunc(simpleHandler))
	assert.Equal(t, "GET:\n/\n  test\n    route\n      simple\n\n\n", router.dumpTree())

	router.addRoute("GET", "/test/route/name", http.HandlerFunc(simpleHandler))
	assert.Equal(t, "GET:\n/\n  test\n    route\n      simple\n      name\n\n\n", router.dumpTree())

	router.addRoute("GET", "/route/core", http.HandlerFunc(simpleHandler))
	assert.Equal(t, "GET:\n/\n  test\n    route\n      simple\n      name\n  route\n    core\n\n\n", router.dumpTree())
}

func TestAddRouteDynamic(t *testing.T) {
	router := newRouter()

	router.addRoute("GET", "/route/:name", http.HandlerFunc(dynamicHandler))
	assert.Equal(t, "GET:\n/\n  route\n    :name\n\n\n", router.dumpTree())
	assert.Equal(t, 1, countDynamicNodes(router.trees.getRoot("GET")))

	router.addRoute("GET", "/route/:name/test/:param/name/", http.HandlerFunc(dynamicHandler))
	assert.Equal(t, "GET:\n/\n  route\n    :name\n      test\n        :param\n          name\n\n\n", router.dumpTree())
	assert.Equal(t, 2, countDynamicNodes(router.trees.getRoot("GET")))
}

func TestResolveTreeNotFound(t *testing.T) {
	router := newRouter()
	req, _ := http.NewRequest("GET", "/", nil)

	node, req := router.resolve(req)
	assert.Nil(t, node)
	assert.Nil(t, req)
}

func TestResolveNotFound(t *testing.T) {
	router := newRouter()
	router.addRoute("GET", "/testroute", http.HandlerFunc(simpleHandler))

	req, _ := http.NewRequest("GET", "/notfound", nil)

	node, req := router.resolve(req)
	assert.Nil(t, node)
	assert.Nil(t, req)

	req, _ = http.NewRequest("GET", "/*2<\\$/-\"5 ", nil)
	node, req = router.resolve(req)
	assert.Nil(t, node)
	assert.Nil(t, req)
}

func TestResolveMethodTreeNotFound(t *testing.T) {
	router := newRouter()
	router.addRoute("GET", "/testroute", http.HandlerFunc(simpleHandler))

	req, _ := http.NewRequest("POST", "/testroute", nil)
	node, req := router.resolve(req)
	assert.Nil(t, node)
	assert.Nil(t, req)
}

func TestResolve(t *testing.T) {
	router := newRouter()
	router.addRoute("GET", "/", http.HandlerFunc(simpleHandler))
	router.addRoute("GET", "/testroute", http.HandlerFunc(simpleHandler))
	router.addRoute("GET", "/route", http.HandlerFunc(simpleHandler))
	router.addRoute("GET", "/route/test/name", routeMiddleware(http.HandlerFunc(simpleHandler)))
	router.addRoute("POST", "/route", http.HandlerFunc(simpleHandler))
	router.addRoute("GET", "/route/:name", routeMiddleware(http.HandlerFunc(dynamicHandler)))
	router.addRoute("POST", "/route/:name", http.HandlerFunc(dynamicHandler))
	router.addRoute("POST", "/route/:name/test/:param/name/", routeMiddleware(routeMiddleware(http.HandlerFunc(dynamicHandler))))

	req, _ := http.NewRequest("GET", "/testroute", nil)
	w := httptest.NewRecorder()
	node, req := router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, simpleHandler, node.fn)

	req, _ = http.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	node, req = router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, simpleHandler, node.fn)

	req, _ = http.NewRequest("GET", "/route/test/name", nil)
	w = httptest.NewRecorder()
	node, req = router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, routeMiddleware(http.HandlerFunc(simpleHandler)), node.fn)

	req, _ = http.NewRequest("POST", "/route", nil)
	w = httptest.NewRecorder()
	node, req = router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, simpleHandler, node.fn)

	req, _ = http.NewRequest("POST", "/route/core/test/attr/name/", nil)
	w = httptest.NewRecorder()
	node, req = router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, routeMiddleware(http.HandlerFunc(dynamicHandler)), node.fn)

	node.fn.ServeHTTP(w, req)
	assert.Equal(t, "route-startroute-startdynamic core attrroute-endroute-end", w.Body.String())

	req, _ = http.NewRequest("POST", "/route/core", nil)
	w = httptest.NewRecorder()
	node, req = router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, dynamicHandler, node.fn)

	node.fn.ServeHTTP(w, req)
	assert.Equal(t, "dynamic core ", w.Body.String())

}

func BenchmarkResolve(b *testing.B) {
	router := newRouter()
	router.addRoute("GET", "/", http.HandlerFunc(simpleHandler))
	router.addRoute("GET", "/testroute", http.HandlerFunc(simpleHandler))
	router.addRoute("GET", "/route", http.HandlerFunc(simpleHandler))
	router.addRoute("GET", "/route/test/name", http.HandlerFunc(simpleHandler))
	router.addRoute("POST", "/route", http.HandlerFunc(simpleHandler))
	router.addRoute("GET", "/route/:name", http.HandlerFunc(dynamicHandler))
	router.addRoute("POST", "/route/:name", http.HandlerFunc(dynamicHandler))
	router.addRoute("POST", "/route/:name/test/:param/name/", http.HandlerFunc(dynamicHandler))
	router.addRoute("GET", "/testroute/abc/def/ghioj/klmn", http.HandlerFunc(simpleHandler))

	req, _ := http.NewRequest("GET", "/testroute/abc/def/ghioj/klmn", nil)

	for i := 0; i < b.N; i++ {
		router.resolve(req)
	}
}

//----------------------------------------------------------------------------------------------------------------------
func countDynamicNodes(n *node) int {
	count := 0
	if n.isDynamic {
		count++
	}

	for _, c := range n.children {
		count += countDynamicNodes(c)
	}

	return count
}

//----------------------------------------------------------------------------------------------------------------------
func simpleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "simple")
}

func dynamicHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	name := ctx.Value("name")
	nameStr := ""
	if name != nil {
		nameStr = name.(string)
	}

	param := ctx.Value("param")
	paramStr := ""
	if param != nil {
		paramStr = param.(string)
	}

	fmt.Fprint(w, "dynamic "+nameStr+" "+paramStr)
}

func routeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "route-start")
		next.ServeHTTP(w, r)
		fmt.Fprint(w, "route-end")
	})
}
