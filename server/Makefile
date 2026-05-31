.PHONY: dev run generate sqlc swagger fmt migrate migrate-up migrate-down migrate-status migrate-new test tidy

dev:
	air

run:
	go run ./cmd/api

generate: sqlc swagger

sqlc:
	sqlc generate

swagger:
	swag init -g main.go -d cmd/api,internal/http -o docs --parseInternal

fmt:
	go fmt ./...

migrate: migrate-up

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

migrate-status:
	go run ./cmd/migrate status

migrate-new:
	test -n "$(name)" || (echo "usage: make migrate-new name=add_players"; exit 1)
	go run ./cmd/migrate new "$(name)"

test:
	go test ./...

tidy:
	go mod tidy
