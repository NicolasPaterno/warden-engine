.PHONY: build run test lint compose-up compose-down migrate-up migrate-down

# Load .env if present so DATABASE_URL etc. are available to targets below.
-include .env
export

MIGRATE = go run -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path db/migrations -database "$(DATABASE_URL)"

build:
	go build -o bin/engine cmd/engine/main.go

run:
	go run cmd/engine/main.go

test:
	go test ./... -race

lint:
	golangci-lint run ./...

# Start local dependencies (Postgres + NATS) and wait until healthy.
compose-up:
	docker compose up -d --wait

compose-down:
	docker compose down

# Apply / roll back schema migrations against $(DATABASE_URL).
migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down 1
