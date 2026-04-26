.PHONY: build test lint vet check

build:
	go build ./...

test:
	go test ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

check: vet lint test

run:
	go run ./cmd/githook/...

migrate:
	@test -n "$(DATABASE_DSN)" || (echo "error: DATABASE_DSN is not set" && exit 1)
	psql "$(DATABASE_DSN)" -f migrations/001_create_webhook_events.sql
