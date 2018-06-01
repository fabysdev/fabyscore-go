# FabysCore GO - Server

HTTP Server Router

## Usage

```go
package main

import "github.com/fabysdev/fabyscore-go/server"
import "net/http"
import "fmt"

func main() {
  srv := server.New();

  srv.GET("/", fabyscoreHandler)

  srv.Run(":8080")
}

func fabyscoreHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Hello!")
}
```

### Routes

```go
// GET
srv.GET("/", fabyscoreHandler)

// POST
srv.POST("/", fabyscoreHandler)

// PUT
srv.PUT("/", fabyscoreHandler)

// DELETE
srv.DELETE("/", fabyscoreHandler)

// OPTIONS
srv.OPTIONS("/", fabyscoreHandler)

// Any
srv.Any("/", fabyscoreHandler)

// Group
srv.Group("/test", func(g *Group) {
  g.GET("/", fabyscoreHandler)
  g.GET("/route", fabyscoreHandler)

  g.POST("/route", fabyscoreHandler)
})
```

#### Not Found Handler

```go
func fabyscoreNotFoundHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "404 - Not Found")
}

srv.SetNotFoundHandler(fabyscoreNotFoundHandler)
```

### Middlewares

Middlewares are standard `net/http` middleware handlers.

```go
func srvMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Println("server start")

    next.ServeHTTP(w, r)

    fmt.Println("server end")
  })
}
```

#### Server Middlewares

Server middlewares are defined on server level and are executed for every request, except for non existing routes. (e.g. a request logger)

```go
srv.Use(srvMiddleware)
srv.UseWithSorting(srvMiddleware, 0)
```

#### Route Middlewares

Route middlewares are defined on the route level.

```go
srv.GET("/", fabyscoreHandler, routeMiddleware, routeMiddlewareTwo)

srv.Group("/test", func(g *Group) {
  g.Use(groupMiddleware)
  g.UseWithSorting(groupMiddleware, 0)

  g.GET("/", fabyscoreHandler, routeMiddleware, routeMiddlewareTwo)
})
```

### Options

Options are used to change settings of the http.Server.

```go
srv.Run(":8080", ReadHeaderTimeout(5*time.Second), IdleTimeout(120*time.Second), WriteTimeout(5*time.Second))
```

#### ReadHeaderTimeout

Option for setting the `http.Server`.`ReadHeaderTimeout`

#### IdleTimeout

Option for setting the `http.Server`.`IdleTimeout`

#### WriteTimeout

Option for setting the `http.Server`.`WriteTimeout`

### ContextKey

Is used to create unique context keys.

```go
var RequestIDContextKey = &ContextKey{"request-id"}
```