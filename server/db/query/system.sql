-- name: GetDatabaseTime :one
SELECT now()::text AS database_time;
