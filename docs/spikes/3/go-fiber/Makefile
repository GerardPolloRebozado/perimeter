.PHONY: tidy build run generate generate-oapi generate-sqlc migrate-up compose-up compose-down clean

tidy:
	go mod tidy

build:
	go build ./...

run:
	go run ./cmd/server

generate: generate-oapi generate-sqlc

generate-oapi:
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1 \
		-config api/oapi-codegen.yaml api/openapi.yaml

generate-sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate

migrate-up:
	go run ./cmd/migrate up

compose-up:
	docker compose up -d postgres redis keycloak

compose-down:
	docker compose down -v

clean:
	rm -f internal/gen/api.gen.go
	rm -rf internal/db/sqlc
