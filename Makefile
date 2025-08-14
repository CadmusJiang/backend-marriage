SHELL := /bin/bash

.PHONY: lint test build ci

lint:
	golangci-lint run --timeout=5m

test:
	go test -race -covermode=atomic -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out | tail -n1

build:
	go build ./...

ci: lint test build


