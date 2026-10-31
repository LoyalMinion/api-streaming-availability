# api-streaming-availability

Stan's transcoder pipeline emits events whenever a video segment becomes available for playback. Downstream CDN origin servers poll a manifest service to discover which segments are ready to serve to players.

This is the Go HTTP service that sits between the transcoder and the CDN.

## Architecture

The store design is in [docs/architecture.md](docs/architecture.md).

## Run locally

- `docker compose up --build`
- `docker compose port api 8080`, `docker compose port prometheus 9090`, and `docker compose port grafana 3000` print the host ports.
- Check `GET /health` and `GET /metrics` on the API port. Grafana shows the Prometheus series.

## Pipeline

- Pull requests run `make check` and `make test`.
- A push to `main` does the same, then calls the Render deploy hook. Render builds the Go service from that commit.
- The service definition is in [render.yaml](render.yaml). Auto-deploy is off, so a push does not go live before the checks pass.
- The deploy hook is the `RENDER_DEPLOY_HOOK_URL` repository secret. Connect the blueprint once from the Render dashboard.
- A free instance sleeps, and the in-memory store is empty after a restart.

## Dependencies

- `github.com/prometheus/client_golang` is the one third-party package.
- `GET /metrics` exposes `http_request_duration_seconds`, a histogram labeled by the matched route pattern.
- The buckets start at 0.1ms, so p50, p95, and p99 stay meaningful for a fast handler.
