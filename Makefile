GO ?= go

.PHONY: check test

check:
	test -z "$$($(GO)fmt -l .)"
	$(GO) vet ./...

test:
	$(GO) test ./...
	$(GO) test -race ./...
