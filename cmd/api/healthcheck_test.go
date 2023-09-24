package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"testing"

	"github.com/chema0/greenlight/internal/assert"
)

func TestHealthcheckHandler(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())

	healthcheckBody, err := json.Marshal(envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": "test",
			"version":     "",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		urlPath     string
		wantCode    int
		wantBody    string
		wantHeaders map[string]string
		varyHeaders []string
	}{
		{
			name:     "Available",
			urlPath:  "/v1/healthcheck",
			wantCode: http.StatusOK,
			wantBody: string(healthcheckBody),
			wantHeaders: map[string]string{
				"Content-Type": "application/json",
			},
			varyHeaders: []string{"Access-Control-Request-Method", "Authorization", "Origin"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, headers, body := ts.get(t, tt.urlPath)

			for k, v := range tt.wantHeaders {
				assert.Equal(t, headers.Get(k), v)
			}

			varyHeaders := headers.Values("Vary")
			for _, v := range tt.varyHeaders {
				assert.Equal(t, slices.Contains(varyHeaders, v), true)
			}

			assert.Equal(t, code, tt.wantCode)
			assert.Equal(t, body, tt.wantBody)
		})
	}
}
