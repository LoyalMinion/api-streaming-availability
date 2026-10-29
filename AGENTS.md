# AGENTS.md

Streaming Availability Service. An in-memory Go HTTP service between the transcoder and the CDN origins.

Module path: `github.com/LoyalMinion/api-streaming-availability`.

Repository layout and process layout are not decided in this file.

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
- Use a table named `testTable` with `Args`, `Fields`, and `Want`.
- Call `t.Parallel` on the test and on each `t.Run` when the cases do not share mutable state that would race.
- Exercise HTTP with `net/http/httptest`.
- Case shape is defined in `.cursor/skills/write-test`.

## Verification

```bash
gofmt -l .
go test ./...
go test -race ./...
```

When a benchmark exists, also run `go test -bench=. ./...`.
