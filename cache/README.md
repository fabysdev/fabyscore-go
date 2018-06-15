# FabysCore GO - Cache

In-Memory key-value store/cache.
It is basically a thread-safe `map[string]interface{}`.

## Usage

```go
package main

import "github.com/fabysdev/fabyscore-go/cache"

func main() {
  // create new cache
  c := cache.New()

  // create new cache with a cleanup routine
  c, stop := cache.NewWithCleanup(10 * time.Minute)

  // add an item
  c.Set("key", "value")

  // add an item with expiration
  c.Set("key", "value", cache.Expire(1 * time.Minute))

  // get an item
  item, found := c.Get("key")

  // delete an item
  c.Delete("key")

  // clear the cache
  c.Clear()

  // get all keys
  keys := c.Keys()

  // delete expired items
  c.DeleteExpired()
}
```
