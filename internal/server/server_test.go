package server_test

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/LoyalMinion/api-streaming-availability/internal/config"
	"github.com/LoyalMinion/api-streaming-availability/internal/server"
)

func TestServer_Serve(t *testing.T) {
	t.Parallel()

	type args struct {
		path string
	}

	type fields struct {
		drain   time.Duration // how long shutdown waits for the request
		release time.Duration // how long the handler stays in flight after shutdown starts
	}

	type want struct {
		status int
		Error  string
	}

	testTable := map[string]struct {
		Args   args
		Fields fields
		Want   want
	}{
		"success - drains an in-flight request": {
			Args:   args{path: "/"},
			Fields: fields{drain: time.Second, release: 20 * time.Millisecond},
			Want:   want{status: http.StatusOK},
		},
		"error - drain timeout passes before the request finishes": {
			Args:   args{path: "/"},
			Fields: fields{drain: 10 * time.Millisecond, release: 200 * time.Millisecond},
			Want:   want{Error: "shutdown: context deadline exceeded"},
		},
	}

	for name, tt := range testTable {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// started fires once the request is inside the handler.
			// The handler then blocks until the test closes release, so the request stays in flight.
			started := make(chan struct{})
			release := make(chan struct{})
			handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				close(started)
				<-release
				w.WriteHeader(http.StatusOK)
			})

			listener, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatalf("listen: %v", err)
			}

			srv := server.New(config.Config{DrainTimeout: tt.Fields.drain}, handler, slog.New(slog.DiscardHandler))
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			served := make(chan error, 1)
			go func() {
				served <- srv.Serve(ctx, listener)
			}()

			responded := make(chan int, 1)
			go func() {
				resp, err := http.Get("http://" + listener.Addr().String() + tt.Args.path)
				if err != nil {
					t.Errorf("request: %v", err)
					responded <- 0

					return
				}
				if err := resp.Body.Close(); err != nil {
					t.Errorf("close body: %v", err)
				}
				responded <- resp.StatusCode
			}()

			// Cancel only after the request has reached the handler, so drain has something to wait for.
			// release is shorter than drain when the request should finish, and longer when the timeout should win.
			<-started
			cancel()
			time.AfterFunc(tt.Fields.release, func() { close(release) })

			err = <-served
			status := <-responded

			if tt.Want.Error != "" {
				if err == nil {
					t.Fatalf("error = nil, want %q", tt.Want.Error)
				}
				if err.Error() != tt.Want.Error {
					t.Errorf("error = %q, want %q", err.Error(), tt.Want.Error)
				}

				return
			}

			if err != nil {
				t.Fatalf("error = %v, want nil", err)
			}
			if status != tt.Want.status {
				t.Errorf("status = %d, want %d", status, tt.Want.status)
			}
		})
	}
}
