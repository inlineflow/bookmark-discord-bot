-- name: CreateChannel :one
INSERT INTO channels (name, external_id)
VALUES (
    ?,
    ?
)
RETURNING *;

-- name: GetChannel :one
SELECT * FROM channels WHERE external_id = ?;
