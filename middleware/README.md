# FabysCore GO - Middleware

`net/http` middleware handlers.


## Usage

```go
package main
  
import "github.com/fabysdev/fabyscore-go/middleware"
  
func main() {
  timeout := middleware.Timeout(1*time.Second)
}
```

## Middlewares

### Timeout

Converts creates a TimeoutHandler and updates the request context with the given timeout.

```go
timeout := middleware.Timeout(1*time.Second)

e.g.
srv.UseWithSorting(middleware.Timeout(1*time.Second), -255)
```