// Package httpapi serves the service's HTTP routes and owns its wire format.
package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// API serves the HTTP routes.
type API struct {
	logger   *slog.Logger
	started  time.Time
	now      func() time.Time
	requests *prometheus.HistogramVec
	reg      *prometheus.Registry
}

// New returns an API that reports uptime since started, read from now.
func New(logger *slog.Logger, started time.Time, now func() time.Time) *API {
	reg := prometheus.NewRegistry()
	requests := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests, labeled by the matched route pattern.",
		Buckets: prometheus.ExponentialBuckets(0.0001, 2, 16),
	}, []string{"pattern"})
	reg.MustRegister(requests)

	return &API{logger: logger, started: started, now: now, requests: requests, reg: reg}
}

// Handler returns the routes served by a.
func (a *API) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", a.health)
	mux.Handle("GET /metrics", promhttp.HandlerFor(a.reg, promhttp.HandlerOpts{}))

	return chain(mux,
		logRequests(a.logger),
		recordDuration(a.requests),
	)
}

func (a *API) health(w http.ResponseWriter, _ *http.Request) {
	a.write(w, http.StatusOK, healthResponse{
		Status:        "ok",
		UptimeSeconds: int(a.now().Sub(a.started).Seconds()),
	})
}

func (a *API) write(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		a.logger.Error("write response", "err", err)
	}
}

type healthResponse struct {
	Status        string `json:"status"`
	UptimeSeconds int    `json:"uptime_seconds"`
}
