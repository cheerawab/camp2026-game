.PHONY: dev run seed import-staff import-staff-teams grant-initial-sitones generate swagger fmt test test-e2e tidy

dev:
	air

run:
	go run ./cmd/api

seed:
	go run ./cmd/seed

import-staff:
ifndef CSV
	$(error CSV is required, e.g. make import-staff CSV="/path/to/staff.csv")
endif
	go run ./cmd/import-staff -csv "$(CSV)" $(ARGS)

import-staff-teams:
	go run ./cmd/import-staff-teams -json "$(or $(JSON),./content/staff_team_assignments.json)" $(ARGS)

grant-initial-sitones:
	go run ./cmd/grant-initial-sitones $(ARGS)

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
