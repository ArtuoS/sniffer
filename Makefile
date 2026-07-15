.PHONY: build run-api run-ingest lint

build:
	go build ./...

run-api:
	go run ./cmd/api

run-ingest:
	go run ./cmd/ingest

lint:
	golangci-lint run ./...
