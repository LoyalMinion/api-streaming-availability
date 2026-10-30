// Package httpapi serves the service's HTTP routes and owns its wire format.
package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

// API serves the HTTP routes.
type API struct {
	logger  *slog.Logger
	started time.Time
	now     func() time.Time
}

// New returns an API that reports uptime since started, read from now.
func New(logger *slog.Logger, started time.Time, now func() time.Time) *API {
	return &API{logger: logger, started: started, now: now}
}

// Handler returns the routes served by a.
func (a *API) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", a.health)

	return mux
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
