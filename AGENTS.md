# AGENTS.md

Streaming Availability Service. An in-memory Go HTTP service between the transcoder and the CDN origins.

Module path: `github.com/LoyalMinion/api-streaming-availability`.

## Layout

| Path | Purpose |
|---|---|
| `cmd/api-streaming-availability/` | Entry point. Builds the logger and config, wires the packages, and runs the server. |
| `internal/config/` | Reads settings from the environment into `Config`. |
| `internal/server/` | Owns `http.Server`: listening, and draining in-flight requests on shutdown. |
| `internal/httpapi/` | Routes, handlers, and request and response types. Owns the wire format. |
| `internal/store/` | The in-memory store in [`docs/architecture.md`](docs/architecture.md). Added with ingest. |
| `docs/` | Design notes. |
| `.tool-versions`, `go.mod` | The Go version. Keep them in step. |
| `Makefile` | `check` and `test`, run locally and in CI. |
| `Dockerfile`, `docker-compose.yml` | Local image, Prometheus, and Grafana. Host ports are assigned by Docker. |
| `prometheus/` | Local scrape config and the Grafana datasource. The series are stored in a Compose volume. Dashboards are configured in Grafana, not in this repo. |
| `.github/workflows/ci.yml` | Pull requests run `make check` and `make test`. `main` then calls the Render deploy hook. |
| `render.yaml` | The Render web service. Auto-deploy is off. |

- `main` only wires. It reads config, builds each dependency, and passes it to a constructor.
- Only `internal/config` reads environment variables.
- Imports point inward: `main` imports `config`, `server`, and `httpapi`. `httpapi` imports `store`. `server` does not import `httpapi`, and `store` does not import `httpapi` or `net/http`.
- Store types carry no JSON tags. `httpapi` converts them to response types.
- Add a package for a responsibility, not for one function.
- Metrics are `GET /metrics` on the same listener. The host exposes one port.

## Constraints

- Standard library only. A third-party package needs a written justification in the README.
- No single `sync.Mutex` wrapping the entire in-memory store.
- State is in memory and resets on restart.
- Segments may arrive in any order. The manifest returns them sorted ascending by `segment_index`.
- On `SIGTERM`, finish in-flight requests before exiting. The drain timeout is configurable.
- Routes are `POST /ingest`, `GET /manifest/{asset_id}`, and `GET /health`. An `asset_id` may contain slashes, so the manifest route captures the rest of the path.
- `go test ./...` and `go test -race ./...` must pass. Include at least one benchmark.

## Coding rules

- Each change is one concern, on a short branch with a meaningful name.
- Delete a feature branch after its pull request merges into main. When another open pull request uses that branch as its base, merge without deleting it, retarget the other pull request to main, and then delete the branch.
- Follow the patterns already in a file unless asked to refactor. When unsure, check a neighbour.
- Only include changes the task needs.
- Code is the source of truth. When docs and code disagree, trust the code and flag the drift.
- Do not add comments to code you did not change.
- Do not add a helper or abstraction for a one-off.
- Run the verification commands before calling a task complete.

## Skills

Skills add to this file. This file does not replace a skill, and a skill does not replace this file. Apply every rule here and every rule in the skill that fits the task. Silence in one is not permission to skip the other.

A skill example shows the shape of a rule. When an example names a third-party package, keep the shape and use the standard library.

Read the skill that fits the task before editing:

- [`.cursor/skills/shape-code/SKILL.md`](.cursor/skills/shape-code/SKILL.md) — writing or reviewing code: function shape, names, checks, and errors.
- [`.cursor/skills/commit-and-pull-request/SKILL.md`](.cursor/skills/commit-and-pull-request/SKILL.md) — a branch name, commit message, or pull request.
- [`.cursor/skills/write-test/SKILL.md`](.cursor/skills/write-test/SKILL.md) — a new test, a test case, or a `_test.go` file.

## Go

### Import order

- Two groups, separated by a blank line: the standard library, then this module.
- An external-package group exists only after a third-party dependency is justified in the README.
- Format with `gofmt`.

### Error handling

- Wrap an error that crosses a boundary with `fmt.Errorf("context: %w", err)`. Use `%w`.
- Never discard an error with `_`.
- Use `errors.Join` when several errors must be returned together.

### Doc comments

- Every exported package, type, func, const, and var has a doc comment that starts with the declared name and ends with a period. No blank line between the comment and the declaration.

### Dependencies

- Inject dependencies through constructors.
- A collaborator depends on the methods it calls, not on a wider concrete type.
- No package-level mutable singletons.

## Testing

- Tests live next to the package they cover, in an external `_test` package.
- Use a table named `testTable` with `Args`, `Fields`, and `Want`. Omit any of the three when the test does not use it.
- Call `t.Parallel` on the test and on each `t.Run` when the cases do not share mutable state that would race.
- Exercise HTTP with `net/http/httptest`.
- Case shape is defined in `.cursor/skills/write-test`.

## Verification

```bash
make check
make test
```

`make check` runs `gofmt` and `go vet`. `make test` runs `go test ./...` and `go test -race ./...`. Also run `go test -bench=. ./...` when a change touches code a benchmark covers.
