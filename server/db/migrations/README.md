# Migrations

Database migrations live here. Keep files ordered with a numeric prefix:

```text
000001_init.sql
000002_add_players.sql
```

The template starts without game tables. Add schema only when a feature needs it.

Use the backend migration tool from `server/`:

```sh
make migrate-new name=add_players
make migrate-status
make migrate
make migrate-down
```

Migration files may use Goose-style `-- +goose Up` and `-- +goose Down`
sections. The local tool applies the `Up` section and records applied versions
in the `schema_migrations` table.
