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

func fabyscoreHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello!")
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
app.Group("/test", func(g *Group) {
  g.GET("/", fabyscoreHandler)
  g.GET("/route", fabyscoreHandler)

  g.POST("/route", fabyscoreHandler)
})
```

####Not Found Handler
```go
func fabyscoreNotFoundHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "404 - Not Found")
}

app.SetNotFoundHandler(fabyscoreNotFoundHandler)
```


###Middlewares

Middlewares are standard `net/http` middleware handlers.

```go
func appMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Println("app start")
    
    next.ServeHTTP(w, r)
    
		fmt.Println("app end")
	})
}
```

####App Middlewares

App middlewares are defined on application level and are executed for every request, except for non existing routes. (e.g. a request logger)

```go
app.Use(appMiddleware)
app.UseWithSort(appMiddleware, 0)
```

####Route Middlewares

Route middlewares are defined on the route level.

```go
app.GET("/", fabyscoreHandler, routeMiddleware, routeMiddlewareTwo)
  
app.Group("/test", func(g *Group) {
  app.Use(groupMiddleware)
  app.UseWithSort(groupMiddleware, 0)

  g.GET("/", fabyscoreHandler, routeMiddleware, routeMiddlewareTwo)
})
```


##License
Code and documentation released under the [MIT license](https://github.com/fabyscore/fabyscore-go/blob/master/LICENSE).
