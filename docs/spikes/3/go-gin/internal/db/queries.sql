-- name: InsertHook :one
INSERT INTO hooks (payload)
VALUES ($1)
RETURNING id, received_at, payload;

-- name: CountHooks :one
SELECT count(*) FROM hooks;
