DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/titanbay?sslmode=disable

.PHONY: build run test migrate-up seed swagger

build:
	go build ./cmd/api

run:
	DATABASE_URL=$(DATABASE_URL) go run ./cmd/api

test:
	go test ./...

migrate-up:
	DATABASE_URL=$(DATABASE_URL) go run ./cmd/migrate

seed:
	DATABASE_URL=$(DATABASE_URL) go run ./cmd/seed

swagger:
	@echo "OpenAPI spec lives at docs/swagger/openapi.yaml"
