.PHONY: dev run generate swagger fmt test test-e2e tidy

dev:
	air

run:
	go run ./cmd/api

generate: swagger

swagger:
	swag init -g main.go -d cmd/api,internal/http -o docs --parseInternal

fmt:
	go fmt ./...

test:
	go test ./...

test-e2e:
	go test -tags=e2e ./internal/e2e

tidy:
	go mod tidy
