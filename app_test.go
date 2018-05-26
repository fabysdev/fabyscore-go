package fabyscore

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	app := NewApp()
	app.GET("/testroute", routeHandler)
	assert.Equal(t, "GET:\n/\n  testroute\n\n\n", app.router.dumpTree())
	req, _ := http.NewRequest("GET", "/testroute", nil)
	node, req := app.router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, routeHandler, node.fn)

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
	app.Group("/test",
		&Route{Method: "GET", Path: "/", Fn: http.HandlerFunc(routeHandler)},
		&Route{Method: "GET", Path: "/route", Fn: http.HandlerFunc(routeHandler)},

		&Route{Method: "POST", Path: "/route", Fn: http.HandlerFunc(routeHandler)},
	)

	tree := app.router.dumpTree()
	assert.Equal(t, "GET:\n/\n  test\n    \n    route\n\n\nPOST:\n/\n  test\n    route\n\n\n", tree)
	req, _ := http.NewRequest("GET", "/test/", nil)
	node, req := app.router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, routeHandler, node.fn)

	req, _ = http.NewRequest("GET", "/test/route", nil)
	node, req = app.router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, routeHandler, node.fn)

	req, _ = http.NewRequest("POST", "/test/route", nil)
	node, req = app.router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assertFuncEquals(t, routeHandler, node.fn)

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
	app.Group("test/",
		&Route{Method: "GET", Path: "", Fn: http.HandlerFunc(routeHandler)},
		&Route{Method: "GET", Path: "route", Fn: http.HandlerFunc(routeHandler)},

		&Route{Method: "POST", Path: "route", Fn: http.HandlerFunc(routeHandler)},
	)
	assert.Equal(t, tree, app.router.dumpTree())

	app = NewApp()
	app.Group("/test",
		&Route{Method: "GET", Path: "/", Fn: http.HandlerFunc(routeHandler)},
		&Route{Method: "GET", Path: "route", Fn: http.HandlerFunc(routeHandler)},

		&Route{Method: "POST", Path: "route", Fn: http.HandlerFunc(routeHandler)},
	)
	assert.Equal(t, tree, app.router.dumpTree())
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

	app.GET("/he", fabyscoreHandler)
	app.GET("/he/2", fabyscoreHandler)
	app.GET("/hello", fabyscoreHandler)
	app.GET("/dyn/", fabyscoreHandler)
	app.GET("/hello/test", fabyscoreHandler)
	app.GET("/hello/test/it", fabyscoreHandler)
	app.GET("/dyn/add/:id", fabyscoreHandler)
	app.GET("/dyn/change/:id/mod/:mod", fabyscoreHandler)
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
	assert.Equal(t, "<nil><nil>contr", w.Body.String())

	// /he/2
	req, _ = http.NewRequest("GET", "/he/2", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	assert.Equal(t, "<nil><nil>contr", w.Body.String())

	// /hello
	req, _ = http.NewRequest("GET", "/hello", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	assert.Equal(t, "<nil><nil>contr", w.Body.String())

	// /dyn/
	req, _ = http.NewRequest("GET", "/dyn/", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	assert.Equal(t, "<nil><nil>contr", w.Body.String())

	// /hello/test
	req, _ = http.NewRequest("GET", "/hello/test", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	assert.Equal(t, "<nil><nil>contr", w.Body.String())

	// /hello/test/it
	req, _ = http.NewRequest("GET", "/hello/test/it", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	assert.Equal(t, "<nil><nil>contr", w.Body.String())

	// /dyn/add/:id
	req, _ = http.NewRequest("GET", "/dyn/add/123", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	assert.Equal(t, "123<nil>contr", w.Body.String())

	// /dyn/change/:id/mod/:mod
	req, _ = http.NewRequest("GET", "/dyn/change/123/mod/asdf", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	assert.Equal(t, "123asdfcontr", w.Body.String())
}

//----------------------------------------------------------------------------------------------------------------------
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

//----------------------------------------------------------------------------------------------------------------------
func assertFuncEquals(t *testing.T, expected interface{}, actual interface{}) {
	if reflect.ValueOf(expected).Pointer() != reflect.ValueOf(actual).Pointer() {
		_, file, line, _ := runtime.Caller(1)
		t.Error("Not equal functions. " + fmt.Sprintf("%s:%d", file, line))
	}
}
