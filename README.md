# FabysCore GO

[![Build Status](https://travis-ci.org/fabysdev/fabyscore-go.svg?branch=master)](https://travis-ci.org/fabysdev/fabyscore-go)

## Warning

The API is not yet stable! Some things will change.  
Please wait for a 1.0 release if you're not willing to update your code frequently.

## Usage

```
dep ensure -add github.com/fabysdev/fabyscore-go

go get github.com/fabysdev/fabyscore-go
```

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

### Middleware

`net/http` middleware handlers.

```go
package main

import "github.com/fabysdev/fabyscore-go/middleware"

func main() {
  timeout := middleware.Timeout(1*time.Second)
}
```

## License

Code and documentation released under the [MIT license](https://github.com/fabysdev/fabyscore-go/blob/master/LICENSE).
