-- name: CreateUser :one
INSERT INTO users ("name", "email")
VALUES ($1, $2)
RETURNING *;
-- name: GetUser :one
SELECT "id",
    "name",
    "email",
    "created_at",
    "updated_at",
    "deleted_at"
FROM users
WHERE id = $1
LIMIT 1;
-- name: ListUsers :many
SELECT "id",
    "name",
    "email",
    "created_at",
    "updated_at",
    "deleted_at"
FROM users
ORDER BY id;
-- name: UpdateName :one
UPDATE users
set name = $1
WHERE id = $2
RETURNING *;
-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
