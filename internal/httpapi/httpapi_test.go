package httpapi_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/LoyalMinion/api-streaming-availability/internal/httpapi"
)

func TestAPI_Health(t *testing.T) {
	t.Parallel()

	type args struct {
		method string
		path   string
	}

	type fields struct {
		started time.Time
		now     func() time.Time
	}

	type body struct {
		Status        string `json:"status"`
		UptimeSeconds int    `json:"uptime_seconds"`
	}

	type want struct {
		status int
		body   body
	}

	started := time.Date(2024, 5, 1, 10, 0, 0, 0, time.UTC)

	testTable := map[string]struct {
		Args   args
		Fields fields
		Want   want
	}{
		"success - reports ok and uptime": {
			Args: args{method: http.MethodGet, path: "/health"},
			Fields: fields{
				started: started,
				now:     func() time.Time { return started.Add(142 * time.Second) },
			},
			Want: want{
				status: http.StatusOK,
				body:   body{Status: "ok", UptimeSeconds: 142},
			},
		},
	}

	for name, tt := range testTable {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			api := httpapi.New(slog.New(slog.DiscardHandler), tt.Fields.started, tt.Fields.now)
			req := httptest.NewRequest(tt.Args.method, tt.Args.path, nil)
			rec := httptest.NewRecorder()

			api.Handler().ServeHTTP(rec, req)

			if rec.Code != tt.Want.status {
				t.Fatalf("status = %d, want %d", rec.Code, tt.Want.status)
			}

			var got body
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if got != tt.Want.body {
				t.Errorf("body = %+v, want %+v", got, tt.Want.body)
			}
		})
	}
}

func TestAPI_RequestLog(t *testing.T) {
	t.Parallel()

	type args struct {
		method string
		path   string
	}

	type entry struct {
		Level   string `json:"level"`
		Msg     string `json:"msg"`
		Method  string `json:"method"`
		Path    string `json:"path"`
		Pattern string `json:"pattern"`
	}

	type want struct {
		status int
		log    entry
	}

	testTable := map[string]struct {
		Args args
		Want want
	}{
		"success - logs the matched route": {
			Args: args{method: http.MethodGet, path: "/health"},
			Want: want{
				status: http.StatusOK,
				log: entry{
					Level:   "INFO",
					Msg:     "request",
					Method:  http.MethodGet,
					Path:    "/health",
					Pattern: "GET /health",
				},
			},
		},
		"success - logs an unmatched path": {
			Args: args{method: http.MethodGet, path: "/missing"},
			Want: want{
				status: http.StatusNotFound,
				log: entry{
					Level:   "INFO",
					Msg:     "request",
					Method:  http.MethodGet,
					Path:    "/missing",
					Pattern: "unmatched",
				},
			},
		},
	}

	for name, tt := range testTable {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			api := httpapi.New(slog.New(slog.NewJSONHandler(&buf, nil)), time.Time{}, time.Now)
			rec := httptest.NewRecorder()

			api.Handler().ServeHTTP(rec, httptest.NewRequest(tt.Args.method, tt.Args.path, nil))

			if rec.Code != tt.Want.status {
				t.Fatalf("status = %d, want %d", rec.Code, tt.Want.status)
			}

			var got entry
			dec := json.NewDecoder(&buf)
			if err := dec.Decode(&got); err != nil {
				t.Fatalf("decode log: %v", err)
			}
			if dec.More() {
				t.Errorf("extra log line")
			}
			if got != tt.Want.log {
				t.Errorf("log = %+v, want %+v", got, tt.Want.log)
			}
		})
	}
}
