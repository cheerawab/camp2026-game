# Camp 2026 Game API

Minimal Go backend template for fast iteration.

## Commands

```sh
make generate
make test
make dev
```

## Endpoints

```text
GET /api/v1/
GET /api/v1/healthz
GET /api/v1/readyz
GET /api/v1/ping
POST /api/v1/examples/validation
GET /api/v1/swagger.json
GET /api/v1/docs/index.html
```

The API docs page is rendered with Scalar and reads the local `swagger.json`.

The example validation endpoint demonstrates `go-playground/validator` with
array-level error locations such as `body.players[0].displayName`.

## Local Config

Copy `.env.example` to `.env` and adjust `DATABASE_URL` when needed.

The app does not run migrations automatically. Add migration files under
`db/migrations` and sqlc queries under `db/query`, then run `make generate`.

## Handler Structure

HTTP handlers are grouped by route group:

```text
internal/http/handler/<group>/
  handler.go
  <endpoint>.go
```

Each endpoint should live in its own file. For example, the system group uses
`root.go`, `health.go`, `ready.go`, and `ping.go`.
