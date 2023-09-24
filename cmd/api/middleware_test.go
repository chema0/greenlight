package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chema0/greenlight/config"
	"github.com/chema0/greenlight/internal/assert"
)

func TestRateLimit(t *testing.T) {
	app := newTestApplication(t)
	cfg := config.NewConfig("test")

	for i := 0; i < cfg.Limiter.Burst; i++ {
		r, err := http.NewRequest(http.MethodGet, "/v1/healthcheck", nil)
		if err != nil {
			t.Fatal(err)
		}

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		})

		app.rateLimit(next).ServeHTTP(httptest.NewRecorder(), r)
	}

	r, err := http.NewRequest(http.MethodGet, "/v1/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	app.rateLimit(next).ServeHTTP(httptest.NewRecorder(), r)

	rr := httptest.NewRecorder()
	app.rateLimit(next).ServeHTTP(rr, r)

	rs := rr.Result()

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "rate limit exceeded")

	// tests := []struct {
	// 	name     string
	// 	urlPath  string
	// 	wantCode int
	// 	wantBody string
	// }{
	// 	{
	// 		name:     "Rate limit exceeded",
	// 		urlPath:  "/v1/healthcheck",
	// 		wantCode: http.StatusTooManyRequests,
	// 		wantBody: "rate limit exceeded",
	// 	},
	// }

	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		ts.get(t, tt.urlPath)
	// 		ts.get(t, tt.urlPath)
	// 		code, _, body := ts.get(t, tt.urlPath)

	// 		assert.Equal(t, code, tt.wantCode)
	// 		assert.Equal(t, body, tt.wantBody)
	// 	})
	// }
}
