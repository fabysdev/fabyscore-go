package fabyscore

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
			t.Errorf("app.Use did not panic")
		}
	}()

	app := NewApp()

	app.GET("/testroute", routeHandler)

	app.Use(appMiddleware)
}

func TestUse(t *testing.T) {
	app := NewApp()

	app.Use(appMiddleware)
	app.UseWithSort(appMiddleware, -255)
	app.UseWithSort(appMiddleware, 4)

	assert.Len(t, app.middlewares, 3)
	assert.Equal(t, -255, app.middlewares[0].sort)
	assert.Equal(t, 0, app.middlewares[1].sort)
	assert.Equal(t, 4, app.middlewares[2].sort)
}

func TestRun(t *testing.T) {
	app := NewApp()
	app.Use(appMiddleware)
	app.UseWithSort(appMiddleware, -255)
	app.UseWithSort(appMiddleware, 4)

	assert.NotNil(t, app.middlewares)
	assert.Len(t, app.middlewares, 3)

	app.Run(":1000000000", ServerReadHeaderTimeout(1*time.Second), ServerIdleTimeout(1*time.Second), ServerWriteTimeout(1*time.Second))

	assert.Nil(t, app.middlewares)
}

func TestRoutes(t *testing.T) {
	app := NewApp()
	app.GET("/testroute", routeHandler, appRouteMiddleware)
	assert.Equal(t, "GET:\n/\n  testroute\n\n\n", app.router.dumpTree())
	req, _ := http.NewRequest("GET", "/testroute", nil)
	node, req := app.router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, appRouteMiddleware(http.HandlerFunc(routeHandler)), node.fn)

	app = NewApp()
	app.POST("/testroute", routeHandler)
	assert.Equal(t, "POST:\n/\n  testroute\n\n\n", app.router.dumpTree())

	app = NewApp()
	app.PUT("/testroute", routeHandler)
	assert.Equal(t, "PUT:\n/\n  testroute\n\n\n", app.router.dumpTree())

	app = NewApp()
	app.DELETE("/testroute", routeHandler)
	assert.Equal(t, "DELETE:\n/\n  testroute\n\n\n", app.router.dumpTree())

	app = NewApp()
	app.PATCH("/testroute", routeHandler)
	assert.Equal(t, "PATCH:\n/\n  testroute\n\n\n", app.router.dumpTree())

	app = NewApp()
	app.HEAD("/testroute", routeHandler)
	assert.Equal(t, "HEAD:\n/\n  testroute\n\n\n", app.router.dumpTree())

	app = NewApp()
	app.OPTIONS("/testroute", routeHandler)
	assert.Equal(t, "OPTIONS:\n/\n  testroute\n\n\n", app.router.dumpTree())

	app = NewApp()
	app.CONNECT("/testroute", routeHandler)
	assert.Equal(t, "CONNECT:\n/\n  testroute\n\n\n", app.router.dumpTree())

	app = NewApp()
	app.TRACE("/testroute", routeHandler)
	assert.Equal(t, "TRACE:\n/\n  testroute\n\n\n", app.router.dumpTree())

	app = NewApp()
	app.Any("/testroute", routeHandler)
	assert.Equal(t, "GET:\n/\n  testroute\n\n\nPOST:\n/\n  testroute\n\n\nPUT:\n/\n  testroute\n\n\nDELETE:\n/\n  testroute\n\n\nPATCH:\n/\n  testroute\n\n\nHEAD:\n/\n  testroute\n\n\nOPTIONS:\n/\n  testroute\n\n\nCONNECT:\n/\n  testroute\n\n\nTRACE:\n/\n  testroute\n\n\n", app.router.dumpTree())
}

func TestGroup(t *testing.T) {
	app := NewApp()
	app.Group("/test", func(g *Group) {
		g.Use(appGroupMiddleware)

		g.GET("/", routeHandler, appRouteMiddleware)
		g.GET("/route", routeHandler)

		g.POST("/route", routeHandler)
	})

	tree := app.router.dumpTree()
	assert.Equal(t, "GET:\n/\n  test\n    \n    route\n\n\nPOST:\n/\n  test\n    route\n\n\n", tree)
	req, _ := http.NewRequest("GET", "/test/", nil)
	node, req := app.router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, appGroupMiddleware(appRouteMiddleware(http.HandlerFunc(routeHandler))), node.fn)

	req, _ = http.NewRequest("GET", "/test/route", nil)
	node, req = app.router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, appGroupMiddleware(http.HandlerFunc(routeHandler)), node.fn)

	req, _ = http.NewRequest("POST", "/test/route", nil)
	node, req = app.router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, appGroupMiddleware(http.HandlerFunc(routeHandler)), node.fn)

	req, _ = http.NewRequest("POST", "/test/", nil)
	node, req = app.router.resolve(req)
	assert.Nil(t, node)

	req, _ = http.NewRequest("POST", "/test", nil)
	node, req = app.router.resolve(req)
	assert.NotNil(t, node)
	assert.Nil(t, node.fn)

	req, _ = http.NewRequest("GET", "/test", nil)
	node, req = app.router.resolve(req)
	assert.NotNil(t, node)
	assert.Nil(t, node.fn)

	app = NewApp()
	app.Group("test/", func(g *Group) {
		g.GET("", routeHandler)
		g.GET("route", routeHandler)

		g.POST("route", routeHandler)
	})
	assert.Equal(t, tree, app.router.dumpTree())

	app = NewApp()
	app.Group("/test", func(g *Group) {
		g.GET("/", routeHandler)
		g.GET("route", routeHandler)

		g.POST("route", routeHandler)
	})
	assert.Equal(t, tree, app.router.dumpTree())
}

func TestGroupMethods(t *testing.T) {
	app := NewApp()
	app.Group("/test", func(g *Group) {
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

	assert.Equal(t, "GET:\n/\n  test\n    route\n\n\nPOST:\n/\n  test\n    route\n\n\nPUT:\n/\n  test\n    route\n\n\nDELETE:\n/\n  test\n    route\n\n\nPATCH:\n/\n  test\n    route\n\n\nHEAD:\n/\n  test\n    route\n\n\nOPTIONS:\n/\n  test\n    route\n\n\nCONNECT:\n/\n  test\n    route\n\n\nTRACE:\n/\n  test\n    route\n\n\n", app.router.dumpTree())
}

func TestGroupUse(t *testing.T) {
	app := NewApp()

	var group *Group
	app.Group("/test", func(g *Group) {
		g.Use(appMiddleware)
		g.UseWithSort(appMiddleware, -255)
		g.UseWithSort(appMiddleware, 4)

		group = g
	})

	assert.Len(t, group.middlewares, 3)
	assert.Equal(t, -255, group.middlewares[0].sort)
	assert.Equal(t, 0, group.middlewares[1].sort)
	assert.Equal(t, 4, group.middlewares[2].sort)
}

func TestGroupUsePanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("group.Use did not panic")
		}
	}()

	app := NewApp()

	app.Group("/test", func(g *Group) {
		g.GET("/route", routeHandler)
		g.Use(appMiddleware)
	})
}

func TestServeHTTPMiddlewareNoNext(t *testing.T) {
	app := NewApp()
	app.Use(appMiddleware)
	app.UseWithSort(appMiddlewareNoNext, -255)

	app.GET("/testroute", routeHandler, appRouteMiddleware, appRouteMiddleware)

	assert.Len(t, app.middlewares, 2)

	req, _ := http.NewRequest("GET", "/testroute", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)
	assert.Equal(t, "no-next", w.Body.String())
}

func TestServeHTTPMiddlewareNoNextWithAppMiddleware(t *testing.T) {
	app := NewApp()
	app.Use(appMiddleware)
	app.UseWithSort(appMiddlewareNoNext, 255)

	app.GET("/testroute", routeHandler, appRouteMiddleware, appRouteMiddleware)

	assert.Len(t, app.middlewares, 2)

	req, _ := http.NewRequest("GET", "/testroute", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)
	assert.Equal(t, "app-startno-nextapp-end", w.Body.String())
}

func TestServeHTTPNotFound(t *testing.T) {
	app := NewApp()
	app.SetNotFoundHandler(appTestNotFoundHandler)

	app.GET("/testroute", routeHandler)

	req, _ := http.NewRequest("GET", "/notfound", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "404")
}

func TestServeHTTPNotFoundDefault(t *testing.T) {
	app := NewApp()

	app.GET("/testroute", routeHandler)

	req, _ := http.NewRequest("GET", "/notfound", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "404")
}

func TestServeHTTPRouteMiddlewareNoNext(t *testing.T) {
	app := NewApp()

	app.GET("/testroute", routeHandler, appRouteMiddleware, appMiddlewareNoNext)

	req, _ := http.NewRequest("GET", "/testroute", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)
	assert.Equal(t, "route-startno-nextroute-end", w.Body.String())
}

func TestServeHTTPRouteMiddlewareNoNextFirst(t *testing.T) {
	app := NewApp()

	app.GET("/testroute", routeHandler, appMiddlewareNoNext, appRouteMiddleware)

	req, _ := http.NewRequest("GET", "/testroute", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)
	assert.Equal(t, "no-next", w.Body.String())
}

func TestServeHTTP(t *testing.T) {
	app := NewApp()

	app.GET("/testroute", routeHandler)

	req, _ := http.NewRequest("GET", "/testroute", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)
	assert.Equal(t, "r", w.Body.String())
}

func TestApp(t *testing.T) {
	app := NewApp()

	app.GET("/he", fabyscoreHandler, ma)
	app.GET("/he/2", fabyscoreHandler, mb)
	app.GET("/hello", fabyscoreHandler, ma)
	app.GET("/dyn/", fabyscoreHandler, ma)
	app.GET("/hello/test", fabyscoreHandler, mb)
	app.GET("/hello/test/it", fabyscoreHandler, mc)
	app.GET("/dyn/add/:id", fabyscoreHandler, ma)
	app.GET("/dyn/change/:id/mod/:mod", fabyscoreHandler, mb)
	app.GET("/", fabyscoreHandler)

	// /
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	assert.Equal(t, "<nil><nil>contr", w.Body.String())

	// /he
	req, _ = http.NewRequest("GET", "/he", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	assert.Equal(t, "a<nil><nil>contr", w.Body.String())

	// /he/2
	req, _ = http.NewRequest("GET", "/he/2", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	assert.Equal(t, "b<nil><nil>contr", w.Body.String())

	// /hello
	req, _ = http.NewRequest("GET", "/hello", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	assert.Equal(t, "a<nil><nil>contr", w.Body.String())

	// /dyn/
	req, _ = http.NewRequest("GET", "/dyn/", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	assert.Equal(t, "a<nil><nil>contr", w.Body.String())

	// /hello/test
	req, _ = http.NewRequest("GET", "/hello/test", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	assert.Equal(t, "b<nil><nil>contr", w.Body.String())

	// /hello/test/it
	req, _ = http.NewRequest("GET", "/hello/test/it", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	assert.Equal(t, "c<nil><nil>contr", w.Body.String())

	// /dyn/add/:id
	req, _ = http.NewRequest("GET", "/dyn/add/123", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	assert.Equal(t, "a123<nil>contr", w.Body.String())

	// /dyn/change/:id/mod/:mod
	req, _ = http.NewRequest("GET", "/dyn/change/123/mod/asdf", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	assert.Equal(t, "b123asdfcontr", w.Body.String())
}

//----------------------------------------------------------------------------------------------------------------------
func appMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "app-start")
		next.ServeHTTP(w, r)
		fmt.Fprint(w, "app-end")
	})
}

func appRouteMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "route-start")
		next.ServeHTTP(w, r)
		fmt.Fprint(w, "route-end")
	})
}

func appGroupMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "group-start")
		next.ServeHTTP(w, r)
		fmt.Fprint(w, "group-end")
	})
}

func appMiddlewareNoNext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "no-next")
	})
}

func routeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "r")
}

func appTestNotFoundHandler(w http.ResponseWriter, r *http.Request) {
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
