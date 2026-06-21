.PHONY: dev build tidy test migrate

dev:
	go run ./cmd/server

build:
	CGO_ENABLED=0 go build -o bin/equidrug-api ./cmd/server

tidy:
	go mod tidy

test:
	go test ./...

migrate:
	@echo "Run migrations via deploy/podman-compose (db init) or:"
	@echo "  psql \$$DATABASE_URL -f migrations/001_initial.sql -f migrations/002_seed.sql"
