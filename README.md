# api-streaming-availability

Stan's transcoder pipeline emits events whenever a video segment becomes available for playback. Downstream CDN origin servers poll a manifest service to discover which segments are ready to serve to players.

This is the Go HTTP service that sits between the transcoder and the CDN.

## Architecture

The store design is in [docs/architecture.md](docs/architecture.md).

## Run

`docker compose up --build` builds the image and serves `GET /health` on port 8080. `PORT` defaults to `8080` and `DRAIN_TIMEOUT` defaults to `15s`, the same as a process started on the host. That image logs revision `dev`. A published build passes `--build-arg REVISION=<commit>` to stamp the commit into the binary.

## Pipeline

Pull requests run `make check` and `make test`. A push to `main` does the same, then calls the Render deploy hook. Render builds the Go service from that commit. The service definition is in [render.yaml](render.yaml), and auto-deploy is off so a push does not go live before the checks pass.

The deploy hook is the `RENDER_DEPLOY_HOOK_URL` repository secret. Connect the blueprint once from the Render dashboard. A free instance sleeps, and the in-memory store is empty after a restart.
