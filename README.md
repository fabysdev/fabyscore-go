# FabysCore GO

[![Build Status](https://travis-ci.org/fabysdev/fabyscore-go.svg?branch=master)](https://travis-ci.org/fabysdev/fabyscore-go)
[![Coverage Status](https://coveralls.io/repos/github/fabysdev/fabyscore-go/badge.svg?branch=master)](https://coveralls.io/github/fabysdev/fabyscore-go?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/fabysdev/fabyscore-go)](https://goreportcard.com/report/github.com/fabysdev/fabyscore-go)

## Warning

The API is not yet stable! Some things will change.  
Please wait for a 1.0 release if you're not willing to update your code frequently.

## Components

### Server

HTTP Server Router

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

#### File Server

Only serves files and not the directory.

```go
srv.ServeFiles("/", http.Dir("./"))
```

### Middleware

`net/http` middleware handlers.

```go
package main

import "github.com/fabysdev/fabyscore-go/middleware"

func main() {
  timeout := middleware.Timeout(1*time.Second)
}
```

### Cache

In-Memory key-value store/cache.

```go
package main

import "github.com/fabysdev/fabyscore-go/cache"

func main() {
  // create new cache
  c := cache.New()

  // add an item
  c.Set("key", "value")

  // get an item
  item, found := c.Get("key")

  // delete an item
  c.Delete("key")
}
```

## License

Code and documentation released under the [MIT license](https://github.com/fabysdev/fabyscore-go/blob/master/LICENSE).
