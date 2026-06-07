.PHONY: build run test lint

build:
	go build -o bin/engine cmd/engine/main.go

run:
	go run cmd/engine/main.go

test:
	go test ./... -race

lint:
	golangci-lint run ./...
