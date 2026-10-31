# Handlers

Implement `POST /ingest` and `GET /manifest/{asset_id}` on the store in [architecture.md](architecture.md). Health, metrics, the request log, and shutdown are already done. Do not change them.

Scenario 7 is a discussion, not a task. Do not add a panic middleware, an `ETag`, eviction, or a `500`.

## How to pick up a task

Each task is one concern, one branch, and one pull request. The branch name is the task name. Change only the files the task lists.

`segment-order` and `asset-store` can be written at the same time, and so can `ingest` and `manifest`. `ingest` and `manifest` both touch `httpapi.go`: each adds only its own route line and keeps the other line on rebase. Both use the `New` signature in the contract. Merge in the numbered order. After the task ahead of yours merges, rebase onto `main`. Do not stack the pull request on another task branch.

A later task calls the signatures in this file. Use those names. Do not rename them in a task.

Before editing, read `AGENTS.md`, `.cursor/skills/shape-code/SKILL.md`, and `.cursor/skills/write-test/SKILL.md`.

- Tests are an external `_test` package. One `testTable`. Use `Args`, `Fields`, and `Want`, and omit any of the three the test does not use.
- `t.Parallel` on the test and on each case, unless the cases share a store.
- HTTP goes through `httptest`. One request in a case. Decode JSON into a struct.
- `make check` and `make test` pass. A task that adds the benchmark also runs `go test -bench=. ./internal/store`.
- No new dependency. Store types have no JSON tags. `httpapi` owns the wire format.

## Contract

```go
type Segment struct {
	Index      int
	CDNURL     string
	ReceivedAt time.Time
}

type Manifest struct {
	AssetID  string
	Segments []Segment
	Updated  time.Time
}

func insert(segments []Segment, segment Segment) (next []Segment, changed bool)

func New(now func() time.Time) *Store
func (s *Store) Accept(assetID string, segment Segment)
func (s *Store) Manifest(assetID string) (Manifest, bool)
```

`insert` never writes into `segments[:len(segments)]`. An index past the end may be appended. An insert or a replacement allocates the next slice.

`New` stores `now`. `Accept` calls `now` inside the asset's writer lock, and only when the write changes the snapshot.

`httpapi.New` gains the store as its last argument. The private interface is the two methods above. The field is `assets`. Health tests and `main` pass `store.New(time.Now)`.

```go
mux.HandleFunc("POST /ingest", a.ingest)
mux.HandleFunc("GET /manifest/{asset_id...}", a.manifest)
```

```go
type errorResponse struct {
	Error string `json:"error"`
}
```

Error and success bodies go through the existing `write` helper. The ingest body is `asset_id`, `segment_index`, `cdn_url`, and `received_at`. Decode into pointers so a missing `segment_index` is not `0`.

| Condition | Status | `error` |
|---|---|---|
| Empty body, malformed JSON, or a wrong type | 400 | `Could not decode request: JSON parsing failed` |
| `received_at` is not RFC3339 | 400 | `Could not decode request: JSON parsing failed` |
| `asset_id` is absent or empty | 400 | `missing required field: asset_id` |
| `segment_index` is absent | 400 | `missing required field: segment_index` |
| `cdn_url` is absent or empty | 400 | `missing required field: cdn_url` |
| `received_at` is absent | 400 | `missing required field: received_at` |

Check the missing fields in that order. `segment_index` `0` is valid. A negative index is unspecified: do not add an error string for it.

A valid body is stored, then the response is `202` and `{"status":"accepted"}`.

The manifest route reads the rest of the path with `r.PathValue("asset_id")`. A missing asset is `404` and `{"error":"asset not found"}`. A hit is `200`:

```json
{
  "asset_id": "show/thunderbirds/season/1/ep3",
  "segments": [
    { "index": 0, "cdn_url": "https://cdn.example.com/thunderbirds/s1e3/seg0.ts" }
  ],
  "last_updated": "2024-05-01T10:00:05Z"
}
```

`index` and `cdn_url` come from the store. `received_at` is not on the wire. `last_updated` is `Updated` in UTC, truncated to a second, formatted with `time.RFC3339`.

The same `asset_id` and `segment_index` stay one segment. A strictly newer `received_at` replaces it, including when `cdn_url` is unchanged, and moves `last_updated`. An older or equal `received_at` changes nothing. The first write wins an equal timestamp. The manifest is sorted ascending by index. An asset is not visible until its first segment is stored. A read takes the registry read lock, then loads the snapshot, and encodes with no lock held.

## Tasks

### 1. `segment-order`

Files: `internal/store/order.go`, `internal/store/order_test.go`.

`insert` covers a new index at the end, a new index at `0`, a new index between two others, a newer replacement, a newer replacement with the same `cdn_url`, an older no-op, and an equal-timestamp no-op. Assert the input slice's elements are unchanged.

Runs at the same time as nothing else. `asset-store` can be written against the signature while this is open, and merges after it.

### 2. `asset-store`

Files: `internal/store/store.go`, `internal/store/store_test.go`.

The registry is a `map[string]*asset` under a `sync.RWMutex`, held only for lookup or insert. Each asset has a writer `sync.Mutex` and an `atomic.Pointer` to an immutable snapshot. Follow the ingest and manifest steps in the architecture doc, including the second check before insert.

Tests, each on its own store:

- Scenarios 1, 2, and 4: order, and no segments crossing assets.
- Scenario 3: one segment per index. The newer duplicate replaces and moves `Updated`. An older duplicate does not.
- A missing asset is `(Manifest{}, false)`.
- `BenchmarkStore_Manifest` reads an asset that already holds segments.

One test may start goroutines. Do not split that test into parallel cases. Cover two assets at once, and many writers on one asset. `make test` is the race run.

### 3. `ingest`

Files: `internal/httpapi/ingest.go`, `internal/httpapi/ingest_test.go`, `internal/httpapi/httpapi.go`, `internal/httpapi/httpapi_test.go`, `cmd/api-streaming-availability/main.go`.

Apply the `New`, route, and `errorResponse` changes from the contract. Add the `POST` line. If the `GET` line is already present, leave it. Pass `store.New(time.Now)` from `main` and from the existing health and request-log tests.

`ingest` decodes, rejects with the table above, calls `Accept`, then writes `202`.

HTTP cases: scenario 5a, 5b, 5c, and 5d; `segment_index` `0`; a body that is stored and returns `202`. The case keeps the `*store.Store` it passed to `New` and asserts `Manifest` in `Want`. That read is the assertion, not a second request.

### 4. `manifest`

Files: `internal/httpapi/manifest.go`, `internal/httpapi/manifest_test.go`, and the `GET` line in `internal/httpapi/httpapi.go`.

Do not change `New`, `main`, or the ingest handler. Add the `GET` line. If the `POST` line is already present, leave it.

Build the response from `Manifest`. Cases: scenario 6, including slashes in the asset id; a `200` whose segments were stored in `Fields` before the request. The response has `index`, `cdn_url`, and `last_updated`, and no `received_at`.
