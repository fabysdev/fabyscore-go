package middleware

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fabysdev/fabyscore-go/server"
	"github.com/stretchr/testify/assert"
)

func TestTimeout(t *testing.T) {
	srv := server.New()

	srv.UseWithSorting(Timeout(1*time.Millisecond), -255)

	srv.GET("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 2)

		w.Write([]byte("test"))
	})

	ts := httptest.NewServer(srv)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusServiceUnavailable, res.StatusCode)

	buf, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.Equal(t, "Request Timeout", string(buf))
}

func TestNoTimeout(t *testing.T) {
	srv := server.New()

	srv.UseWithSorting(Timeout(1*time.Second), -255)

	srv.GET("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	})

	ts := httptest.NewServer(srv)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	buf, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.Equal(t, "test", string(buf))
}
