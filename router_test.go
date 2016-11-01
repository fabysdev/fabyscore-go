package fabyscore

import (
    "testing"
    "fmt"
    "net/http"
    "github.com/stretchr/testify/assert"
    "net/http/httptest"
)

func simpleHandler(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request) {
    fmt.Fprint(w, "simple")
    return w, r
}

func dynamicHandler(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request) {
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
    return w, r
}

func routerMod(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request) {
    fmt.Fprint(w, "mod")
    return w, r
}

var modifiers = Modifiers{}


func TestAddRouteSimple(t *testing.T) {
    router := newRouter()
    router.addRoute("GET", "/testroute", simpleHandler, modifiers)
    
    assert.Equal(t, "GET:\n/\n  testroute\n\n\n", router.dumpTree())

    router.addRoute("POST", "/route", simpleHandler, modifiers)
    router.addRoute("GET", "/route", simpleHandler, modifiers)

    assert.Equal(t, "GET:\n/\n  testroute\n  route\n\n\nPOST:\n/\n  route\n\n\n", router.dumpTree())
}

func TestAddRouteSimpleRoot(t *testing.T) {
    router := newRouter()
    router.addRoute("GET", "/", simpleHandler, modifiers)
    router.addRoute("GET", "/testroute", simpleHandler, modifiers)
    
    assert.Equal(t, "GET:\n/\n  testroute\n\n\n", router.dumpTree())
}

func TestAddRouteSimpleTree(t *testing.T) {
    router := newRouter()
    
    router.addRoute("GET", "/test/route/simple", simpleHandler, modifiers)
    assert.Equal(t, "GET:\n/\n  test\n    route\n      simple\n\n\n", router.dumpTree())

    router.addRoute("GET", "/test/route/name", simpleHandler, modifiers)
    assert.Equal(t, "GET:\n/\n  test\n    route\n      simple\n      name\n\n\n", router.dumpTree())

    router.addRoute("GET", "/route/core", simpleHandler, modifiers)
    assert.Equal(t, "GET:\n/\n  test\n    route\n      simple\n      name\n  route\n    core\n\n\n", router.dumpTree())
}

func TestAddRouteDynamic(t *testing.T) {
    router := newRouter()

    router.addRoute("GET", "/route/:name", dynamicHandler, modifiers)
    assert.Equal(t, "GET:\n/\n  route\n    :name\n\n\n", router.dumpTree())
    assert.Equal(t, 1, countDynamicNodes(router.trees.getRoot("GET")))

    router.addRoute("GET", "/route/:name/test/:param/name/", dynamicHandler, modifiers)
    assert.Equal(t, "GET:\n/\n  route\n    :name\n      test\n        :param\n          name\n            \n\n\n", router.dumpTree())
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
    router.addRoute("GET", "/testroute", simpleHandler, modifiers)
    
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
    router.addRoute("GET", "/testroute", simpleHandler, modifiers)

    req, _ := http.NewRequest("POST", "/testroute", nil)
    node, req := router.resolve(req)
    assert.Nil(t, node)
    assert.Nil(t, req)
}

func TestResolve(t *testing.T) {
    router := newRouter()
    router.addRoute("GET", "/", simpleHandler, modifiers)
    router.addRoute("GET", "/testroute", simpleHandler, modifiers)
    router.addRoute("GET", "/route", simpleHandler, modifiers)
    router.addRoute("GET", "/route/test/name", simpleHandler, Modifiers{NewModifier(1, routerMod)})
    router.addRoute("POST", "/route", simpleHandler, modifiers)
    router.addRoute("GET", "/route/:name", dynamicHandler, Modifiers{NewModifier(1, routerMod)})
    router.addRoute("POST", "/route/:name", dynamicHandler, modifiers)
    router.addRoute("POST", "/route/:name/test/:param/name/", dynamicHandler, Modifiers{NewModifier(1, routerMod), NewModifier(2, routerMod)})


    req, _ := http.NewRequest("GET", "/testroute", nil)
    w := httptest.NewRecorder()
    node, req := router.resolve(req)
    assert.NotNil(t, node)
    assert.NotNil(t, node.fn)
    assert.Empty(t, node.modifiers)
    assertFuncEquals(t, simpleHandler, node.fn)


    req, _ = http.NewRequest("GET", "/", nil)
    w = httptest.NewRecorder()
    node, req = router.resolve(req)
    assert.NotNil(t, node)
    assert.NotNil(t, node.fn)
    assert.Empty(t, node.modifiers)
    assertFuncEquals(t, simpleHandler, node.fn)


    req, _ = http.NewRequest("GET", "/route/test/name", nil)
    w = httptest.NewRecorder()
    node, req = router.resolve(req)
    assert.NotNil(t, node)
    assert.NotNil(t, node.fn)
    assert.Len(t, node.modifiers, 1)
    assertFuncEquals(t, simpleHandler, node.fn)
    

    req, _ = http.NewRequest("POST", "/route", nil)
    w = httptest.NewRecorder()
    node, req = router.resolve(req)
    assert.NotNil(t, node)
    assert.NotNil(t, node.fn)
    assert.Empty(t, node.modifiers)
    assertFuncEquals(t, simpleHandler, node.fn)


    req, _ = http.NewRequest("POST", "/route/core/test/attr/name/", nil)
    w = httptest.NewRecorder()
    node, req = router.resolve(req)
    assert.NotNil(t, node)
    assert.NotNil(t, node.fn)
    assert.Len(t, node.modifiers, 2)
    assertFuncEquals(t, dynamicHandler, node.fn)
    
    node.fn(w, req)
    assert.Equal(t, "dynamic core attr", w.Body.String())


    req, _ = http.NewRequest("POST", "/route/core", nil)
    w = httptest.NewRecorder()
    node, req = router.resolve(req)
    assert.NotNil(t, node)
    assert.NotNil(t, node.fn)
    assert.Empty(t, node.modifiers)
    assertFuncEquals(t, dynamicHandler, node.fn)
    
    node.fn(w, req)
    assert.Equal(t, "dynamic core ", w.Body.String())
    
}

func TestModifierSort(t *testing.T) {
    router := newRouter()
    router.addRoute("GET", "/testroute", simpleHandler, Modifiers{NewModifier(4, routerMod), NewModifier(1, routerMod)})

    req, _ := http.NewRequest("GET", "/testroute", nil)
    node, req := router.resolve(req)
    assert.NotNil(t, node)
    assert.NotNil(t, node.fn)
    assert.Len(t, node.modifiers, 2)
    assertFuncEquals(t, simpleHandler, node.fn)
    
    assert.Equal(t, 1, node.modifiers[0].sort)
    assert.Equal(t, 4, node.modifiers[1].sort)
}

func countDynamicNodes(n *node) int {
    count := 0
    if n.isDynamic {
        count++
    }
    
    for _, c := range n.children {
        count += countDynamicNodes(c)
    }
    
    return count;
}


