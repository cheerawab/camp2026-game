.PHONY: dev generate sqlc swagger test tidy

dev:
	go run ./cmd/api

generate: sqlc swagger

sqlc:
	sqlc generate

swagger:
	swag init -g main.go -d cmd/api,internal/http -o docs --parseInternal

test:
	go test ./...

tidy:
	go mod tidy
