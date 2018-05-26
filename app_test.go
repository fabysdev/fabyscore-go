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

func TestAddModifier(t *testing.T) {
	app := NewApp()

	app.AddModifier(0, appMod)
	app.AddModifier(-255, appMod)
	app.AddModifier(4, appMod)

	assert.Len(t, app.modifiers, 3)
	assert.Equal(t, 0, app.modifiers[0].sort)
	assert.Equal(t, -255, app.modifiers[1].sort)
	assert.Equal(t, 4, app.modifiers[2].sort)
}

func TestAddBeforeModifier(t *testing.T) {
	app := NewApp()

	app.AddBeforeModifier(512, appBeforeMod)
	app.AddBeforeModifier(-255, appBeforeMod)
	app.AddBeforeModifier(4, appBeforeMod)

	assert.Len(t, app.beforeModifiers, 3)
	assert.Equal(t, 512, app.beforeModifiers[0].sort)
	assert.Equal(t, -255, app.beforeModifiers[1].sort)
	assert.Equal(t, 4, app.beforeModifiers[2].sort)
}

func TestRun(t *testing.T) {
	app := NewApp()
	app.AddModifier(0, appMod)
	app.AddModifier(-255, appMod)
	app.AddModifier(4, appMod)

	app.AddBeforeModifier(512, appBeforeMod)
	app.AddBeforeModifier(-255, appBeforeMod)
	app.AddBeforeModifier(4, appBeforeMod)

	app.Run(":1000000000")

	assert.Len(t, app.modifiers, 3)
	assert.Equal(t, -255, app.modifiers[0].sort)
	assert.Equal(t, 0, app.modifiers[1].sort)
	assert.Equal(t, 4, app.modifiers[2].sort)

	assert.Len(t, app.beforeModifiers, 3)
	assert.Equal(t, -255, app.beforeModifiers[0].sort)
	assert.Equal(t, 4, app.beforeModifiers[1].sort)
	assert.Equal(t, 512, app.beforeModifiers[2].sort)
}

func TestRoutes(t *testing.T) {
	app := NewApp()
	app.GET("/testroute", routeHandler, NewModifier(0, appRouteMod), NewModifier(0, appRouteMod))
	assert.Equal(t, "GET:\n/\n  testroute\n\n\n", app.router.dumpTree())
	req, _ := http.NewRequest("GET", "/testroute", nil)
	node, req := app.router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assert.Len(t, node.modifiers, 2)
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
	app.Group("/test", Modifiers{NewModifier(0, routerMod)},
		&Route{Method: "GET", Path: "/", Fn: routeHandler, Modifiers: Modifiers{NewModifier(1, routerMod)}},
		&Route{Method: "GET", Path: "/route", Fn: routeHandler},

		&Route{Method: "POST", Path: "/route", Fn: routeHandler},
	)

	tree := app.router.dumpTree()
	assert.Equal(t, "GET:\n/\n  test\n    \n    route\n\n\nPOST:\n/\n  test\n    route\n\n\n", tree)
	req, _ := http.NewRequest("GET", "/test/", nil)
	node, req := app.router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assert.Len(t, node.modifiers, 2)
	assertFuncEquals(t, routeHandler, node.fn)

	req, _ = http.NewRequest("GET", "/test/route", nil)
	node, req = app.router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assert.Len(t, node.modifiers, 1)
	assertFuncEquals(t, routeHandler, node.fn)

	req, _ = http.NewRequest("POST", "/test/route", nil)
	node, req = app.router.resolve(req)
	assert.NotNil(t, node)
	assert.NotNil(t, node.fn)
	assert.Len(t, node.modifiers, 1)
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
	app.Group("test/", Modifiers{NewModifier(0, routerMod)},
		&Route{Method: "GET", Path: "", Fn: routeHandler},
		&Route{Method: "GET", Path: "route", Fn: routeHandler},

		&Route{Method: "POST", Path: "route", Fn: routeHandler},
	)
	assert.Equal(t, tree, app.router.dumpTree())

	app = NewApp()
	app.Group("/test", Modifiers{NewModifier(0, routerMod)},
		&Route{Method: "GET", Path: "/", Fn: routeHandler},
		&Route{Method: "GET", Path: "route", Fn: routeHandler},

		&Route{Method: "POST", Path: "route", Fn: routeHandler},
	)
	assert.Equal(t, tree, app.router.dumpTree())
}

func TestServeHTTPBeforeReturns(t *testing.T) {
	app := NewApp()
	app.AddBeforeModifier(0, appBeforeModReturns)

	app.GET("/testroute", routeHandler, NewModifier(0, appRouteMod), NewModifier(0, appRouteMod))

	assert.Len(t, app.beforeModifiers, 1)

	req, _ := http.NewRequest("GET", "/testroute", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)
	assert.Equal(t, "returned", w.Body.String())
}

func TestServeHTTPNotFound(t *testing.T) {
	app := NewApp()
	app.SetNotFoundHandler(appTestNotFoundHandler)
	app.AddBeforeModifier(0, appBeforeMod)

	app.GET("/testroute", routeHandler, NewModifier(0, appRouteMod), NewModifier(0, appRouteMod))

	assert.Len(t, app.beforeModifiers, 1)

	req, _ := http.NewRequest("GET", "/notfound", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "404")
}

func TestServeHTTPNotFoundDefault(t *testing.T) {
	app := NewApp()
	app.AddBeforeModifier(0, appBeforeMod)

	app.GET("/testroute", routeHandler, NewModifier(0, appRouteMod), NewModifier(0, appRouteMod))

	assert.Len(t, app.beforeModifiers, 1)

	req, _ := http.NewRequest("GET", "/notfound", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "404")
}

func TestServeHTTPAppModReturns(t *testing.T) {
	app := NewApp()
	app.AddModifier(0, appModReturns)

	app.GET("/testroute", routeHandler, NewModifier(0, appRouteMod), NewModifier(0, appRouteMod))

	req, _ := http.NewRequest("GET", "/testroute", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)
	assert.Equal(t, "returned", w.Body.String())
}

func TestServeHTTPRouteModReturns(t *testing.T) {
	app := NewApp()

	app.GET("/testroute", routeHandler, NewModifier(0, appRouteMod), NewModifier(0, appRouteModReturns))

	req, _ := http.NewRequest("GET", "/testroute", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)
	assert.Equal(t, "routereturned", w.Body.String())
}

func TestServeHTTP(t *testing.T) {
	app := NewApp()

	app.GET("/testroute", routeHandler, NewModifier(0, appRouteMod), NewModifier(0, appRouteMod))

	req, _ := http.NewRequest("GET", "/testroute", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)
	assert.Equal(t, "routerouter", w.Body.String())
}

func TestApp(t *testing.T) {
	app := NewApp()

	app.GET("/he", fabyscoreHandler, NewModifier(1, ma))
	app.GET("/he/2", fabyscoreHandler, NewModifier(1, mb))
	app.GET("/hello", fabyscoreHandler, NewModifier(1, ma))
	app.GET("/dyn/", fabyscoreHandler, NewModifier(1, ma))
	app.GET("/hello/test", fabyscoreHandler, NewModifier(1, mb))
	app.GET("/hello/test/it", fabyscoreHandler, NewModifier(1, mc))
	app.GET("/dyn/add/:id", fabyscoreHandler, NewModifier(1, ma))
	app.GET("/dyn/change/:id/mod/:mod", fabyscoreHandler, NewModifier(1, mb))
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
func appMod(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request) {
	fmt.Fprint(w, "mod")
	return w, r
}

func appBeforeMod(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request) {
	fmt.Fprint(w, "before")
	return w, r
}

func appRouteMod(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request) {
	fmt.Fprint(w, "route")
	return w, r
}

func routeHandler(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request) {
	fmt.Fprint(w, "r")
	return w, r
}

func appBeforeModReturns(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request) {
	fmt.Fprint(w, "returned")
	return w, nil
}

func appModReturns(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request) {
	fmt.Fprint(w, "returned")
	return w, nil
}

func appRouteModReturns(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request) {
	fmt.Fprint(w, "returned")
	return w, nil
}

func appTestNotFoundHandler(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte("404"))
	return w, r
}

func fabyscoreHandler(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request) {
	fmt.Fprint(w, r.Context().Value("id"))
	fmt.Fprint(w, r.Context().Value("mod"))
	fmt.Fprint(w, "contr")
	return w, r
}

func mb(w http.ResponseWriter, req *http.Request) (http.ResponseWriter, *http.Request) {
	fmt.Fprint(w, "b")
	return w, req
}

func ma(w http.ResponseWriter, req *http.Request) (http.ResponseWriter, *http.Request) {
	fmt.Fprint(w, "a")
	return w, req
}

func mc(w http.ResponseWriter, req *http.Request) (http.ResponseWriter, *http.Request) {
	fmt.Fprint(w, "c")
	return w, req
}

//----------------------------------------------------------------------------------------------------------------------
func assertFuncEquals(t *testing.T, expected interface{}, actual interface{}) {
	if reflect.ValueOf(expected).Pointer() != reflect.ValueOf(actual).Pointer() {
		_, file, line, _ := runtime.Caller(1)
		t.Error("Not equal functions. " + fmt.Sprintf("%s:%d", file, line))
	}
}
