# FabysCore GO - Decode

Decode request data into a Go type based on the Content-Type header.

## Usage

```go
package main

import "github.com/fabysdev/fabyscore-go/decode"
import "net/http"

type userReq struct {
  Name  string `query:"name" form:"name" json:"name"`
  Email string `query:"email" form:"email" json:"email"`
}

func fabyscoreHandler(w http.ResponseWriter, r *http.Request) {
  ur := new(userReq)

  // based on Content-Type header (query + content-type)
  err := decode.Request(r, ur)

  // query
  err := decode.Query(r, ur)

  // form
  err := decode.Form(r, ur)

  // json
  err := decode.JSON(r, ur)
}
```

## Supported Struct Field Types

- bool
- float (float32, float64)
- int (int, int8, int16, int32, int64)
- string
- uint (uint, uint8, uint16, uint32, uint64)
- pointer to one of the above types
