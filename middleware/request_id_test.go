package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/fabysdev/fabyscore-go/server"
	"github.com/stretchr/testify/assert"
)

func TestRequestIDRandReadFailedPanic(t *testing.T) {
	defer func() {
		osHostname = os.Hostname
		randRead = rand.Read

		if r := recover(); r == nil {
			t.Errorf("RequestID did not panic")
		}
	}()

	osHostname = func() (string, error) { return "", nil }
	randRead = func(b []byte) (n int, err error) { return 0, errors.New("error") }

	RequestID("test")
}

func TestRequestIDRandomLengthPanic(t *testing.T) {
	defer func() {
		osHostname = os.Hostname
		randRead = rand.Read
		base64Encode = base64.StdEncoding.EncodeToString

		if r := recover(); r == nil {
			t.Errorf("RequestID did not panic")
		}
	}()

	osHostname = func() (string, error) { return "", nil }
	randRead = func(b []byte) (n int, err error) { return 0, nil }
	base64Encode = func(src []byte) string { return "" }

	RequestID("test")
}

func TestRequestIDAlreadyExists(t *testing.T) {
	defer func() {
		osHostname = os.Hostname
	}()

	osHostname = func() (string, error) { return "", nil }

	srv := server.New()
	srv.Use(RequestID("test"))
	srv.GET("/", func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(RequestIDContextKey)
		w.Write([]byte(id.(string)))
	})

	req, _ := http.NewRequest("GET", "/", nil)
	req = req.WithContext(context.WithValue(req.Context(), RequestIDContextKey, "existingkey"))

	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)
	assert.Equal(t, "existingkey", w.Body.String())
}

func TestRequestID(t *testing.T) {
	defer func() {
		osHostname = os.Hostname
		base64Encode = base64.StdEncoding.EncodeToString
	}()

	osHostname = func() (string, error) { return "", nil }
	base64Encode = func(src []byte) string { return "randomstring" }

	srv := server.New()
	srv.Use(RequestID("prefix"))
	srv.GET("/", func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(RequestIDContextKey)
		w.Write([]byte(id.(string)))
	})

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)
	assert.Equal(t, "prefix-localhost-randomstring-1", w.Body.String())
}

func TestRequestIDCounterIncreases(t *testing.T) {
	defer func() {
		osHostname = os.Hostname
		base64Encode = base64.StdEncoding.EncodeToString
	}()

	osHostname = func() (string, error) { return "", nil }
	base64Encode = func(src []byte) string { return "randomstring" }

	srv := server.New()
	srv.Use(RequestID("prefix"))
	srv.GET("/", func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(RequestIDContextKey)
		w.Write([]byte(id.(string)))
	})

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	assert.Equal(t, "prefix-localhost-randomstring-1", w.Body.String())

	req, _ = http.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	assert.Equal(t, "prefix-localhost-randomstring-2", w.Body.String())
}

func TestGetRequestIDNoEntry(t *testing.T) {
	ctx := context.Background()

	v := GetRequestID(ctx)
	assert.Equal(t, v, "")
}

func TestGetRequestIDNotAString(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, RequestIDContextKey, 123)

	v := GetRequestID(ctx)
	assert.Equal(t, v, "")
}

func TestGetRequestID(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, RequestIDContextKey, "requestid")

	v := GetRequestID(ctx)
	assert.Equal(t, v, "requestid")
}
