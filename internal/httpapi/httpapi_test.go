package httpapi_test

import (
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
