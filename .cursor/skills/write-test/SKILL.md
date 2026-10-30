---
name: write-test
description: >-
  Write or reshape a Go test to match this repo's table-test convention.
  Use when adding a test, a test case, a _test.go file, or when a new test
  does not match neighbouring tests. Covers httptest cases for ingest,
  manifest, and health, and race-safe tests.
---

# Write test

- Before writing a test, open the newest test in the same package and copy its shape. If an older test in the same file disagrees with this skill, follow this skill.
- Do a boundary analysis of the feature under test and keep only the cases that cover a distinct path. Drop cases that repeat a path already covered, and drop cases that cannot change coverage.
- Use the standard library. Do not add testify, sqlmock, autogold, or another test helper module.

## Shape

- Package is the external `_test` package.
- Call `t.Parallel()` on the test and on each `t.Run` when the test does not share mutable state that would race. Do not decide by whether a neighbouring test calls it. A case that shares one store with another case does not call `t.Parallel()`.
- One table named `testTable`. Each case is `map[string]struct { Args args; Fields fields; Want want }`.
- Declare `args`, `fields`, and `want` inside the test function, directly below the first `t.Parallel()`. Do not declare them at package level, so each test in a package owns its own.
- Args is the input. Fields holds dependencies, such as a store or a fixed clock. Want holds the expected value and the error string. The error field is `Error`, not `Err`.
- Case names are "success - ..." or "error - ...".
- One act per case. Do not loop the function under test inside a case. Split that into two cases.
- On failure, compare `err.Error()` with `Want.Error`. On success, `Want.Error` is empty and the act returns a nil error. Fail with `t.Fatalf` when the test cannot continue, and `t.Errorf` otherwise.
- Give each case its own dependencies through Fields. Do not share a mutable store across parallel cases.

## Concurrency and time

- Block a handler on a channel when the test needs an in-flight request. Release that channel when the case says the request should finish. Do not sleep to stand in for the request.

## HTTP

- Build the request with `httptest.NewRequest` and record the response with `httptest.NewRecorder`.
- One request per case. Pass the method, path, and body through Args.
- Want carries the status code and the decoded response. Decode JSON into a struct. Do not assert with a substring match on the raw body.
- Cover ingest, manifest, and health. Include the malformed ingest bodies and the missing-asset manifest from the challenge.

Do not invent a flatter table, or a second call to the function under test inside one case.
