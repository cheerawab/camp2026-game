# Camp 2026 Game API

Minimal Go backend template for fast iteration.

## Commands

```sh
make generate
make migrate
make test
make dev
```

`make dev` runs Air hot reload. Use `make run` to run `go run ./cmd/api`
without file watching.

`make migrate` applies pending migrations from `db/migrations` using
`DATABASE_URL`. Use `make migrate-status` to inspect applied migrations,
`make migrate-down` to revert the latest migration, and
`make migrate-new name=add_players` to create the next migration file.

## Design Docs

```text
docs/game-data-api-design.md
```

This document maps the root game design into proposed database tables and API
routes.

## Endpoints

```text
GET /api/healthz
POST /api/auth/login
POST /api/auth/logout
GET /api/users/state
GET /api/bingo/boards
POST /api/bingo/missions/{missionID}/complete
POST /api/bingo/line-rewards/{lineRewardID}/claim
GET /api/matches
POST /api/match-pairings
POST /api/matches
GET /api/matches/{matchID}
GET /api/matches/{matchID}/ws
POST /api/matches/{matchID}/answers
GET /api/qrcode/me
POST /api/qrcode/scans
GET /api/world-bosses
GET /api/world-bosses/{bossID}
POST /api/world-bosses/{bossID}/matches
GET /api/storage
GET /api/storage/sitones
GET /api/storage/items
GET /api/storage/recipes
POST /api/storage/crafting
GET /api/catalog/sitones
GET /api/catalog/items
GET /api/catalog/recipes
POST /api/staff/rewards
POST /api/staff/activity-verifications
GET /api/swagger.json
GET /api/docs
GET /api/docs/index.html
```

The API docs page is rendered with Scalar and reads the local `swagger.json`.
Swagger is generated with `swaggo/swag` from Go handler annotations:

```sh
make swagger
```

The game feature endpoints are currently contract stubs for API review. GET
routes return example JSON. POST action routes decode and validate the request,
then return `501 Not Implemented` until the business logic is built.

## API DTOs

Shared API contract DTOs live in:

```text
internal/http/apimodel/
```

These structs describe request and response shapes for Swagger and handler
binding. Do not use generated sqlc row types as public responses unless that
shape is intentionally part of the API.

## Local Config

Copy `.env.example` to `.env` and adjust `DATABASE_URL` when needed.
Runtime settings and secrets stay in env. Game content definitions such as
sitones, items, crafting recipes, bingo boards, bingo missions, and world boss
definitions should be loaded from JSON files, not env.

The app does not run migrations automatically. Add migration files under
`db/migrations`, apply them with `make migrate`, add sqlc queries under
`db/query`, then run `make generate`.

## Handler Structure

HTTP handlers are grouped by route group:

```text
internal/http/handler/<group>/
  handler.go
  <endpoint>.go
```

Each endpoint should live in its own file. For example, the system group uses
`health.go`.
