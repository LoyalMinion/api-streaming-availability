package httpapi

import (
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func logRequests(logger *slog.Logger) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)

			logger.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"pattern", routePattern(r),
			)
		})
	}
}

func recordDuration(requests *prometheus.HistogramVec) middleware {
	return func(next http.Handler) http.Handler {
		return promhttp.InstrumentHandlerDuration(requests, next, promhttp.WithLabelFromRequest("pattern", routePattern))
	}
}

func routePattern(r *http.Request) string {
	if r.Pattern == "" {
		return "unmatched"
	}
	return r.Pattern
}

func chain(handler http.Handler, middlewares ...middleware) http.Handler {
	// Loop backwards so execution happens in the order they are passed.
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

type middleware func(http.Handler) http.Handler
