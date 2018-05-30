package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUsePanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("srv.Use did not panic")
		}
	}()

	srv := New()

	srv.GET("/testroute", routeHandler)

	srv.Use(srvMiddleware)
}

func TestUse(t *testing.T) {
	srv := New()

	srv.Use(srvMiddleware)
	srv.UseWithSorting(srvMiddleware, -255)
	srv.UseWithSorting(srvMiddleware, 4)

	assert.Len(t, srv.middlewares, 3)
	assert.Equal(t, -255, srv.middlewares[0].sorting)
	assert.Equal(t, 0, srv.middlewares[1].sorting)
	assert.Equal(t, 4, srv.middlewares[2].sorting)
}

func TestRun(t *testing.T) {
	srv := New()
	srv.Use(srvMiddleware)
	srv.UseWithSorting(srvMiddleware, -255)
	srv.UseWithSorting(srvMiddleware, 4)

	assert.NotNil(t, srv.middlewares)
	assert.Len(t, srv.middlewares, 3)

	srv.Run(":1000000000", ReadHeaderTimeout(1*time.Second), IdleTimeout(1*time.Second), WriteTimeout(1*time.Second))

	assert.Nil(t, srv.middlewares)
}

func TestRoutes(t *testing.T) {
	srv := New()
	srv.GET("/testroute", routeHandler, srvRouteMiddleware)
	assert.Equal(t, "GET:\n/\n  testroute\n\n\n", srv.router.dumpTree())
	req, _ := http.NewRequest("GET", "/testroute", nil)
	node, req := srv.router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, srvRouteMiddleware(http.HandlerFunc(routeHandler)), node.fn)

	srv = New()
	srv.POST("/testroute", routeHandler)
	assert.Equal(t, "POST:\n/\n  testroute\n\n\n", srv.router.dumpTree())

	srv = New()
	srv.PUT("/testroute", routeHandler)
	assert.Equal(t, "PUT:\n/\n  testroute\n\n\n", srv.router.dumpTree())

	srv = New()
	srv.DELETE("/testroute", routeHandler)
	assert.Equal(t, "DELETE:\n/\n  testroute\n\n\n", srv.router.dumpTree())

	srv = New()
	srv.PATCH("/testroute", routeHandler)
	assert.Equal(t, "PATCH:\n/\n  testroute\n\n\n", srv.router.dumpTree())

	srv = New()
	srv.HEAD("/testroute", routeHandler)
	assert.Equal(t, "HEAD:\n/\n  testroute\n\n\n", srv.router.dumpTree())

	srv = New()
	srv.OPTIONS("/testroute", routeHandler)
	assert.Equal(t, "OPTIONS:\n/\n  testroute\n\n\n", srv.router.dumpTree())

	srv = New()
	srv.CONNECT("/testroute", routeHandler)
	assert.Equal(t, "CONNECT:\n/\n  testroute\n\n\n", srv.router.dumpTree())

	srv = New()
	srv.TRACE("/testroute", routeHandler)
	assert.Equal(t, "TRACE:\n/\n  testroute\n\n\n", srv.router.dumpTree())

	srv = New()
	srv.Any("/testroute", routeHandler)
	assert.Equal(t, "GET:\n/\n  testroute\n\n\nPOST:\n/\n  testroute\n\n\nPUT:\n/\n  testroute\n\n\nDELETE:\n/\n  testroute\n\n\nPATCH:\n/\n  testroute\n\n\nHEAD:\n/\n  testroute\n\n\nOPTIONS:\n/\n  testroute\n\n\nCONNECT:\n/\n  testroute\n\n\nTRACE:\n/\n  testroute\n\n\n", srv.router.dumpTree())
}

func TestGroup(t *testing.T) {
	srv := New()
	srv.Group("/test", func(g *Group) {
		g.Use(srvGroupMiddleware)

		g.GET("/", routeHandler, srvRouteMiddleware)
		g.GET("/route", routeHandler)

		g.POST("/route", routeHandler)
	})

	tree := srv.router.dumpTree()
	assert.Equal(t, "GET:\n/\n  test\n    route\n\n\nPOST:\n/\n  test\n    route\n\n\n", tree)
	req, _ := http.NewRequest("GET", "/test/", nil)
	node, req := srv.router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, srvGroupMiddleware(srvRouteMiddleware(http.HandlerFunc(routeHandler))), node.fn)

	req, _ = http.NewRequest("GET", "/test/route", nil)
	node, req = srv.router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, srvGroupMiddleware(http.HandlerFunc(routeHandler)), node.fn)

	req, _ = http.NewRequest("POST", "/test/route", nil)
	node, req = srv.router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, srvGroupMiddleware(http.HandlerFunc(routeHandler)), node.fn)

	req, _ = http.NewRequest("POST", "/test/", nil)
	node, req = srv.router.resolve(req)
	assert.NotNil(t, node)
	assert.Nil(t, node.fn)

	req, _ = http.NewRequest("POST", "/test", nil)
	node, req = srv.router.resolve(req)
	assert.NotNil(t, node)
	assert.Nil(t, node.fn)

	req, _ = http.NewRequest("GET", "/test", nil)
	node, req = srv.router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)

	srv = New()
	srv.Group("test/", func(g *Group) {
		g.GET("", routeHandler)
		g.GET("route", routeHandler)

		g.POST("route", routeHandler)
	})
	assert.Equal(t, tree, srv.router.dumpTree())

	srv = New()
	srv.Group("/test", func(g *Group) {
		g.GET("/", routeHandler)
		g.GET("route", routeHandler)

		g.POST("route", routeHandler)
	})
	assert.Equal(t, tree, srv.router.dumpTree())
}

func TestGroupMethods(t *testing.T) {
	srv := New()
	srv.Group("/test", func(g *Group) {
		g.GET("/route", routeHandler)
		g.POST("/route", routeHandler)
		g.PUT("/route", routeHandler)
		g.DELETE("/route", routeHandler)
		g.PATCH("/route", routeHandler)
		g.HEAD("/route", routeHandler)
		g.OPTIONS("/route", routeHandler)
		g.CONNECT("/route", routeHandler)
		g.TRACE("/route", routeHandler)
	})

	assert.Equal(t, "GET:\n/\n  test\n    route\n\n\nPOST:\n/\n  test\n    route\n\n\nPUT:\n/\n  test\n    route\n\n\nDELETE:\n/\n  test\n    route\n\n\nPATCH:\n/\n  test\n    route\n\n\nHEAD:\n/\n  test\n    route\n\n\nOPTIONS:\n/\n  test\n    route\n\n\nCONNECT:\n/\n  test\n    route\n\n\nTRACE:\n/\n  test\n    route\n\n\n", srv.router.dumpTree())
}

func TestGroupUse(t *testing.T) {
	srv := New()

	var group *Group
	srv.Group("/test", func(g *Group) {
		g.Use(srvMiddleware)
		g.UseWithSorting(srvMiddleware, -255)
		g.UseWithSorting(srvMiddleware, 4)

		group = g
	})

	assert.Len(t, group.middlewares, 3)
	assert.Equal(t, -255, group.middlewares[0].sorting)
	assert.Equal(t, 0, group.middlewares[1].sorting)
	assert.Equal(t, 4, group.middlewares[2].sorting)
}

func TestGroupUsePanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("group.Use did not panic")
		}
	}()

	srv := New()

	srv.Group("/test", func(g *Group) {
		g.GET("/route", routeHandler)
		g.Use(srvMiddleware)
	})
}

func TestServeHTTPMiddlewareNoNext(t *testing.T) {
	srv := New()
	srv.Use(srvMiddleware)
	srv.UseWithSorting(srvMiddlewareNoNext, -255)

	srv.GET("/testroute", routeHandler, srvRouteMiddleware, srvRouteMiddleware)

	assert.Len(t, srv.middlewares, 2)

	req, _ := http.NewRequest("GET", "/testroute", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)
	assert.Equal(t, "no-next", w.Body.String())
}

func TestServeHTTPMiddlewareNoNextWithsrvMiddleware(t *testing.T) {
	srv := New()
	srv.Use(srvMiddleware)
	srv.UseWithSorting(srvMiddlewareNoNext, 255)

	srv.GET("/testroute", routeHandler, srvRouteMiddleware, srvRouteMiddleware)

	assert.Len(t, srv.middlewares, 2)

	req, _ := http.NewRequest("GET", "/testroute", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)
	assert.Equal(t, "srv-startno-nextsrv-end", w.Body.String())
}

func TestServeHTTPNotFound(t *testing.T) {
	srv := New()
	srv.SetNotFoundHandler(srvTestNotFoundHandler)

	srv.GET("/testroute", routeHandler)

	req, _ := http.NewRequest("GET", "/notfound", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "404")
}

func TestServeHTTPNotFoundDefault(t *testing.T) {
	srv := New()

	srv.GET("/testroute", routeHandler)

	req, _ := http.NewRequest("GET", "/notfound", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "404")
}

func TestServeHTTPRouteMiddlewareNoNext(t *testing.T) {
	srv := New()

	srv.GET("/testroute", routeHandler, srvRouteMiddleware, srvMiddlewareNoNext)

	req, _ := http.NewRequest("GET", "/testroute", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)
	assert.Equal(t, "route-startno-nextroute-end", w.Body.String())
}

func TestServeHTTPRouteMiddlewareNoNextFirst(t *testing.T) {
	srv := New()

	srv.GET("/testroute", routeHandler, srvMiddlewareNoNext, srvRouteMiddleware)

	req, _ := http.NewRequest("GET", "/testroute", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)
	assert.Equal(t, "no-next", w.Body.String())
}

func TestServeHTTP(t *testing.T) {
	srv := New()

	srv.GET("/testroute", routeHandler)

	req, _ := http.NewRequest("GET", "/testroute", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)
	assert.Equal(t, "r", w.Body.String())
}

func TestServer(t *testing.T) {
	srv := New()

	srv.GET("/he", fabyscoreHandler, ma)
	srv.GET("/he/2", fabyscoreHandler, mb)
	srv.GET("/hello", fabyscoreHandler, ma)
	srv.GET("/dyn/", fabyscoreHandler, ma)
	srv.GET("/hello/test", fabyscoreHandler, mb)
	srv.GET("/hello/test/it", fabyscoreHandler, mc)
	srv.GET("/dyn/add/:id", fabyscoreHandler, ma)
	srv.GET("/dyn/change/:id/mod/:mod", fabyscoreHandler, mb)
	srv.GET("/", fabyscoreHandler)

	// /
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	assert.Equal(t, "<nil><nil>contr", w.Body.String())

	// /he
	req, _ = http.NewRequest("GET", "/he", nil)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	assert.Equal(t, "a<nil><nil>contr", w.Body.String())

	// /he/2
	req, _ = http.NewRequest("GET", "/he/2", nil)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	assert.Equal(t, "b<nil><nil>contr", w.Body.String())

	// /hello
	req, _ = http.NewRequest("GET", "/hello", nil)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	assert.Equal(t, "a<nil><nil>contr", w.Body.String())

	// /dyn/
	req, _ = http.NewRequest("GET", "/dyn/", nil)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	assert.Equal(t, "a<nil><nil>contr", w.Body.String())

	// /hello/test
	req, _ = http.NewRequest("GET", "/hello/test", nil)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	assert.Equal(t, "b<nil><nil>contr", w.Body.String())

	// /hello/test/it
	req, _ = http.NewRequest("GET", "/hello/test/it", nil)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	assert.Equal(t, "c<nil><nil>contr", w.Body.String())

	// /dyn/add/:id
	req, _ = http.NewRequest("GET", "/dyn/add/123", nil)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	assert.Equal(t, "a123<nil>contr", w.Body.String())

	// /dyn/change/:id/mod/:mod
	req, _ = http.NewRequest("GET", "/dyn/change/123/mod/asdf", nil)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	assert.Equal(t, "b123asdfcontr", w.Body.String())
}

//----------------------------------------------------------------------------------------------------------------------
func srvMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "srv-start")
		next.ServeHTTP(w, r)
		fmt.Fprint(w, "srv-end")
	})
}

func srvRouteMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "route-start")
		next.ServeHTTP(w, r)
		fmt.Fprint(w, "route-end")
	})
}

func srvGroupMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "group-start")
		next.ServeHTTP(w, r)
		fmt.Fprint(w, "group-end")
	})
}

func srvMiddlewareNoNext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "no-next")
	})
}

func routeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "r")
}

func srvTestNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte("404"))
}

func fabyscoreHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, r.Context().Value("id"))
	fmt.Fprint(w, r.Context().Value("mod"))
	fmt.Fprint(w, "contr")
}

func ma(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "a")
		next.ServeHTTP(w, r)
	})
}

func mb(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "b")
		next.ServeHTTP(w, r)
	})
}

func mc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "c")
		next.ServeHTTP(w, r)
	})
}

//----------------------------------------------------------------------------------------------------------------------
func assertFuncEquals(t *testing.T, expected interface{}, actual interface{}) {
	if reflect.ValueOf(expected).Pointer() != reflect.ValueOf(actual).Pointer() {
		_, file, line, _ := runtime.Caller(1)
		t.Error("Not equal functions. " + fmt.Sprintf("%s:%d", file, line))
	}
}