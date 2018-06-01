package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/fabysdev/fabyscore-go/server"
)

var osHostname = os.Hostname
var randRead = rand.Read
var base64Encode = base64.StdEncoding.EncodeToString

// RequestIDContextKey is the request id context key.
var RequestIDContextKey = &server.ContextKey{"request-id"}

// RequestID adds a request id into the request context.
// Has the form prefix-hostname-random-counter, e.g. http-localhost-ueT39830ghyR-1
// The random is only generated once per RequestID middleware (panics if generating the random fails).
func RequestID(prefix string) func(http.Handler) http.Handler {
	// resolve hostname
	hostname, err := osHostname()
	if err != nil || hostname == "" {
		hostname = "localhost"
	}

	// generate random
	b := make([]byte, 16)
	_, err = randRead(b)
	if err != nil {
		panic(fmt.Sprintf("Generating the random failed. Error: %v", err))
	}

	random := base64Encode(b)
	random = strings.NewReplacer("+", "", "/", "", "=", "").Replace(random)

	if len(random) < 12 {
		panic("Could not generate random string")
	}

	// create final prefix
	prefix = fmt.Sprintf("%s-%s-%s", prefix, hostname, random[:12])

	// create middleware
	var counter uint64
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if ctx.Value(RequestIDContextKey) != nil {
				next.ServeHTTP(w, r)
				return
			}

			num := atomic.AddUint64(&counter, 1)
			ctx = context.WithValue(ctx, RequestIDContextKey, fmt.Sprintf("%s-%d", prefix, num))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetRequestID returns the request id from the given context.
// Returns an empty string if no request id was found.
func GetRequestID(ctx context.Context) string {
	v := ctx.Value(RequestIDContextKey)
	if v == nil {
		return ""
	}

	if id, ok := v.(string); ok {
		return id
	}

	return ""
}
