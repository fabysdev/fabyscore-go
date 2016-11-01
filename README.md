#FabysCore GO



##Warning

The API is not yet stable! Some things will change.  
Please wait for a 1.0 release if you're not willing to update your code frequently.


##Usage

```
go get github.com/fabyscore/fabyscore-go
```

```go
package main
  
import "github.com/fabyscore/fabyscore-go"
import "net/http"
import "fmt"
  
func main() {
    app := fabyscore.NewApp();
        
    app.GET("/", fabyscoreHandler)
    
    app.Run(":8080")
}
  
func fabyscoreHandler(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request) {
    fmt.Fprint(w, "Hello!")
    return w, r
}

```

###Routes
```go
// GET
app.GET("/", fabyscoreHandler)
  
// POST
app.POST("/", fabyscoreHandler)
  
// PUT
app.PUT("/", fabyscoreHandler)
  
// DELETE
app.DELETE("/", fabyscoreHandler)
  
// OPTIONS
app.OPTIONS("/", fabyscoreHandler)
  
// Any
app.Any("/", fabyscoreHandler)
  
// Group
app.Group("/test", Modifiers{}, 
    &Route{Method: "GET", Path: "/", Fn: fabyscoreHandler},
    &Route{Method: "GET", Path: "/route", Fn: fabyscoreHandler},
    
    &Route{Method: "POST", Path: "/route", Fn: fabyscoreHandler},
)
```

####Not Found Handler
```go
func fabyscoreNotFoundHandler(w http.ResponseWriter, req *http.Request) (http.ResponseWriter, *http.Request) {
    fmt.Fprint(w, "404 - Not Found")
    return w, req
}

app.SetNotFoundHandler(fabyscoreNotFoundHandler)
```


###Modifiers

A modifier can change the response and/or the request before the route handler is invoked.   
No further modifier (or the route handler) will be invoked if a modifier does not return a request.

```go
func fabyscoreModifier(w http.ResponseWriter, req *http.Request) (http.ResponseWriter, *http.Request) {
    req.Header.Add("X-Test", "Test")
    return w, req
}
```

####App Modifiers

App modifiers are defined on application level and are executed for every request. (e.g. a request logger)

```go
app.AddBeforeModifier(0, fabyscoreModifier)
```

#####BeforeModifiers

Before Modifiers are executed before the route is resolved.

```go
app.AddModifier(0, fabyscoreModifier)
```

####Route Modifiers

Route modifiers are defined on the route level.

```go
app.GET("/", fabyscoreHandler, fabyscore.NewModifier(0, fabyscoreModifier), fabyscore.NewModifier(1, fabyscoreModifier))
  
app.Group("/test", fabyscore.Modifiers{fabyscore.NewModifier(0, fabyscoreModifier)},
    &fabyscore.Route{Method: "GET", Path: "/", Fn: fabyscoreHandler, Modifiers: fabyscore.Modifiers{fabyscore.NewModifier(0, fabyscoreModifier)}},
)
```



##License
Code and documentation released under the [MIT license](https://github.com/fabyscore/fabyscore-go/blob/master/LICENSE).
