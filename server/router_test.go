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

	router.addRoute("GET", "/route/:name", http.HandlerFunc(dynamicHandler))
	assert.Equal(t, "GET:\n/\n  route\n    :name\n\n\n", router.dumpTree())
	assert.Equal(t, 1, countDynamicNodes(router.trees.getRoot("GET")))

	router.addRoute("GET", "/route/:name/test/:param/name/", http.HandlerFunc(dynamicHandler))
	assert.Equal(t, "GET:\n/\n  route\n    :name\n      test\n        :param\n          name\n\n\n", router.dumpTree())
	assert.Equal(t, 2, countDynamicNodes(router.trees.getRoot("GET")))
}

func TestAddRouteMatchAll(t *testing.T) {
	router := newRouter()

	router.addRoute("GET", "/route/*path", http.HandlerFunc(matchallHandler))
	assert.Equal(t, "GET:\n/\n  route\n    *path\n\n\n", router.dumpTree())
	assert.Equal(t, 1, countMatchAllNodes(router.trees.getRoot("GET")))

	router.addRoute("GET", "/route/*path", http.HandlerFunc(matchallHandler))
	assert.Equal(t, "GET:\n/\n  route\n    *path\n\n\n", router.dumpTree())
	assert.Equal(t, 1, countMatchAllNodes(router.trees.getRoot("GET")))
}

func TestResolveTreeNotFound(t *testing.T) {
	router := newRouter()
	req, _ := http.NewRequest("GET", "/", nil)

	node, req, _ := router.resolve(req)
	assert.Nil(t, node)
	assert.NotNil(t, req)
}

func TestResolveNotFound(t *testing.T) {
	router := newRouter()
	router.addRoute("GET", "/testroute", http.HandlerFunc(simpleHandler))

	req, _ := http.NewRequest("GET", "/notfound", nil)

	node, req, _ := router.resolve(req)
	assert.Nil(t, node)
	assert.NotNil(t, req)

	req, _ = http.NewRequest("GET", "/*2<\\$/-\"5 ", nil)
	node, req, _ = router.resolve(req)
	assert.Nil(t, node)
	assert.NotNil(t, req)
}

func TestResolveMethodTreeNotFound(t *testing.T) {
	router := newRouter()
	router.addRoute("GET", "/testroute", http.HandlerFunc(simpleHandler))

	req, _ := http.NewRequest("POST", "/testroute", nil)
	node, req, _ := router.resolve(req)
	assert.Nil(t, node)
	assert.NotNil(t, req)
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
	node, _, _ := router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, simpleHandler, node.fn)

	req, _ = http.NewRequest("GET", "/", nil)
	node, _, _ = router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, simpleHandler, node.fn)

	req, _ = http.NewRequest("GET", "/route/test/name", nil)
	node, _, _ = router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, routeMiddleware(http.HandlerFunc(simpleHandler)), node.fn)

	req, _ = http.NewRequest("POST", "/route", nil)
	node, _, _ = router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, simpleHandler, node.fn)

	req, _ = http.NewRequest("POST", "/route/core/test/attr/name/", nil)
	node, req, _ = router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, routeMiddleware(http.HandlerFunc(dynamicHandler)), node.fn)

	w := httptest.NewRecorder()
	node.fn.ServeHTTP(w, req)
	assert.Equal(t, "route-startroute-startdynamic core attrroute-endroute-end", w.Body.String())

	req, _ = http.NewRequest("POST", "/route/core", nil)
	node, req, _ = router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, dynamicHandler, node.fn)

	w = httptest.NewRecorder()
	node.fn.ServeHTTP(w, req)
	assert.Equal(t, "dynamic core ", w.Body.String())
}

func TestResolveDynamicPathIndex(t *testing.T) {
	router := newRouter()
	router.addRoute("GET", "/static/:path", http.HandlerFunc(matchallHandler))

	req, _ := http.NewRequest("GET", "/static/", nil)
	node, req, _ := router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, matchallHandler, node.fn)

	w := httptest.NewRecorder()
	node.fn.ServeHTTP(w, req)
	assert.Equal(t, "match-all ", w.Body.String())
}

func TestResolveMatchAll(t *testing.T) {
	router := newRouter()
	router.addRoute("GET", "/route/*path", http.HandlerFunc(matchallHandler))

	req, _ := http.NewRequest("GET", "/route/core", nil)
	node, req, _ := router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, matchallHandler, node.fn)

	w := httptest.NewRecorder()
	node.fn.ServeHTTP(w, req)
	assert.Equal(t, "match-all core", w.Body.String())

	req, _ = http.NewRequest("GET", "/route/core/", nil)
	node, req, _ = router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, matchallHandler, node.fn)

	w = httptest.NewRecorder()
	node.fn.ServeHTTP(w, req)
	assert.Equal(t, "match-all core/", w.Body.String())

	req, _ = http.NewRequest("GET", "/route/core/test/attr", nil)
	node, req, _ = router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, matchallHandler, node.fn)

	w = httptest.NewRecorder()
	node.fn.ServeHTTP(w, req)
	assert.Equal(t, "match-all core/test/attr", w.Body.String())

	req, _ = http.NewRequest("GET", "/route/core/test/attr/name/", nil)
	node, req, _ = router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, matchallHandler, node.fn)

	w = httptest.NewRecorder()
	node.fn.ServeHTTP(w, req)
	assert.Equal(t, "match-all core/test/attr/name/", w.Body.String())

	req, _ = http.NewRequest("GET", "/route/core/test/attr/name/file.txt", nil)
	node, req, _ = router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, matchallHandler, node.fn)

	w = httptest.NewRecorder()
	node.fn.ServeHTTP(w, req)
	assert.Equal(t, "match-all core/test/attr/name/file.txt", w.Body.String())
}

func TestResolveMatchAllNoEndingSlash(t *testing.T) {
	router := newRouter()
	router.addRoute("GET", "/route/*path", http.HandlerFunc(matchallHandler))

	req, _ := http.NewRequest("GET", "/route/core", nil)
	node, req, _ := router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, matchallHandler, node.fn)

	w := httptest.NewRecorder()
	node.fn.ServeHTTP(w, req)
	assert.Equal(t, "match-all core", w.Body.String())

	req, _ = http.NewRequest("GET", "/route/core/", nil)
	node, req, _ = router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, matchallHandler, node.fn)

	w = httptest.NewRecorder()
	node.fn.ServeHTTP(w, req)
	assert.Equal(t, "match-all core/", w.Body.String())

	req, _ = http.NewRequest("GET", "/route/core/test/attr", nil)
	node, req, _ = router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, matchallHandler, node.fn)

	w = httptest.NewRecorder()
	node.fn.ServeHTTP(w, req)
	assert.Equal(t, "match-all core/test/attr", w.Body.String())

	req, _ = http.NewRequest("GET", "/route/core/test/attr/name/", nil)
	node, req, _ = router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, matchallHandler, node.fn)

	w = httptest.NewRecorder()
	node.fn.ServeHTTP(w, req)
	assert.Equal(t, "match-all core/test/attr/name/", w.Body.String())

	req, _ = http.NewRequest("GET", "/route/core/test/attr/name/file.txt", nil)
	node, req, _ = router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, matchallHandler, node.fn)

	w = httptest.NewRecorder()
	node.fn.ServeHTTP(w, req)
	assert.Equal(t, "match-all core/test/attr/name/file.txt", w.Body.String())
}

func TestResolveMatchAllIndex(t *testing.T) {
	router := newRouter()
	router.addRoute("GET", "/*path", http.HandlerFunc(matchallHandler))

	req, _ := http.NewRequest("GET", "/", nil)
	node, req, _ := router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, matchallHandler, node.fn)

	w := httptest.NewRecorder()
	node.fn.ServeHTTP(w, req)
	assert.Equal(t, "match-all ", w.Body.String())
}

func TestResolveMatchAllPathIndex(t *testing.T) {
	router := newRouter()
	router.addRoute("GET", "/static/*path", http.HandlerFunc(matchallHandler))

	req, _ := http.NewRequest("GET", "/static/", nil)
	node, req, _ := router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, matchallHandler, node.fn)

	w := httptest.NewRecorder()
	node.fn.ServeHTTP(w, req)
	assert.Equal(t, "match-all ", w.Body.String())
}

func TestResolveMatchAllIndexWithMultipleRoutes(t *testing.T) {
	router := newRouter()
	router.addRoute("GET", "/test", http.HandlerFunc(simpleHandler))
	router.addRoute("GET", "/*path", http.HandlerFunc(matchallHandler))

	// match all req
	req, _ := http.NewRequest("GET", "/", nil)
	node, req, _ := router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, matchallHandler, node.fn)

	w := httptest.NewRecorder()
	node.fn.ServeHTTP(w, req)
	assert.Equal(t, "match-all ", w.Body.String())

	req, _ = http.NewRequest("GET", "/a", nil)
	node, req, _ = router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, matchallHandler, node.fn)

	w = httptest.NewRecorder()
	node.fn.ServeHTTP(w, req)
	assert.Equal(t, "match-all a", w.Body.String())

	req, _ = http.NewRequest("GET", "/a/b", nil)
	node, req, _ = router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, matchallHandler, node.fn)

	w = httptest.NewRecorder()
	node.fn.ServeHTTP(w, req)
	assert.Equal(t, "match-all a/b", w.Body.String())

	// simple req
	req, _ = http.NewRequest("GET", "/test", nil)
	node, req, _ = router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, simpleHandler, node.fn)

	w = httptest.NewRecorder()
	node.fn.ServeHTTP(w, req)
	assert.Equal(t, "simple", w.Body.String())
}

func TestResolveDynamicWithMultipleRoutes(t *testing.T) {
	router := newRouter()
	router.addRoute("GET", "/test", http.HandlerFunc(simpleHandler))
	router.addRoute("GET", "/:name/a", http.HandlerFunc(dynamicHandler))
	router.addRoute("GET", "/:name/b", http.HandlerFunc(dynamicHandler))

	// dynamic
	req, _ := http.NewRequest("GET", "/name_a/a", nil)
	node, req, _ := router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, dynamicHandler, node.fn)

	w := httptest.NewRecorder()
	node.fn.ServeHTTP(w, req)
	assert.Equal(t, "dynamic name_a ", w.Body.String())

	req, _ = http.NewRequest("GET", "/name_b/b", nil)
	node, req, _ = router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, dynamicHandler, node.fn)

	w = httptest.NewRecorder()
	node.fn.ServeHTTP(w, req)
	assert.Equal(t, "dynamic name_b ", w.Body.String())

	// simple req
	req, _ = http.NewRequest("GET", "/test", nil)
	node, req, _ = router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, simpleHandler, node.fn)

	w = httptest.NewRecorder()
	node.fn.ServeHTTP(w, req)
	assert.Equal(t, "simple", w.Body.String())
}

func TestAddRoutePanicsConflictingMatchAll(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("router.addRoute did not panic for conflicting match-all route")
		}

		assert.Equal(t, "Route '/test/a' can not be added. Match-All route '/*path' conflicts with it. Check the route registration order.", r)
	}()

	router := newRouter()
	router.addRoute("GET", "/*path", http.HandlerFunc(matchallHandler))
	router.addRoute("GET", "/test/a", http.HandlerFunc(simpleHandler))
}

func TestAddRoutePanicsConflictingMatchAllComplex(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("router.addRoute did not panic for conflicting match-all route")
		}

		assert.Equal(t, "Route '/route/test/a' can not be added. Match-All route '/route/*path' conflicts with it. Check the route registration order.", r)
	}()

	router := newRouter()
	router.addRoute("GET", "/route/*path", http.HandlerFunc(matchallHandler))
	router.addRoute("GET", "/route/test/a", http.HandlerFunc(simpleHandler))
}

func TestAddRoutePanicsMatchAllIneffectiveParts(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("router.addRoute did not panic for ineffective parts of a match-all route")
		}

		assert.Equal(t, "Route '/*path/a' has ineffective parts. Everything after the Match-All part '/*path' is ignored. Remove the ineffective parts from the route.", r)
	}()

	router := newRouter()
	router.addRoute("GET", "/*path/a", http.HandlerFunc(matchallHandler))
}

func TestAddRoutePanicsConflictingDynamic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("router.addRoute did not panic for conflicting dynamic route")
		}

		assert.Equal(t, "Route '/test/a' can not be added. Dynamic route '/:name' conflicts with it. Check the route registration order.", r)
	}()

	router := newRouter()
	router.addRoute("GET", "/:name", http.HandlerFunc(dynamicHandler))
	router.addRoute("GET", "/test/a", http.HandlerFunc(simpleHandler))
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

func countMatchAllNodes(n *node) int {
	count := 0
	if n.isMatchAll {
		count++
	}

	for _, c := range n.children {
		count += countMatchAllNodes(c)
	}

	return count
}

//----------------------------------------------------------------------------------------------------------------------
func simpleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "simple")
}

func dynamicHandler(w http.ResponseWriter, r *http.Request) {
	nameStr := Param(r, "name")
	paramStr := Param(r, "param")

	fmt.Fprint(w, "dynamic "+nameStr+" "+paramStr)
}

func matchallHandler(w http.ResponseWriter, r *http.Request) {
	pathStr := Param(r, "path")

	fmt.Fprint(w, "match-all "+pathStr)
}

func routeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "route-start")
		next.ServeHTTP(w, r)
		fmt.Fprint(w, "route-end")
	})
}
