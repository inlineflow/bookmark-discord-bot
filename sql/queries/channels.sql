-- name: CreateChannel :one
INSERT INTO channels (name, external_id)
VALUES (
    $1,
    $2
)
RETURNING *;

-- name: GetChannel :one
SELECT * FROM channels WHERE external_id = $1;
