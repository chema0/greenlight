package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chema0/greenlight/config"
	"github.com/chema0/greenlight/internal/data"
	"github.com/chema0/greenlight/testhelpers"
)

func newTestApplication(t *testing.T) *application {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	db := testhelpers.NewTestingDB()

	return &application{
		config: config.NewConfig("test"),
		logger: logger,
		models: data.NewModels(db),
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}
