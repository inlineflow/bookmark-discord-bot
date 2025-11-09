-- name: CreateGuild :one
INSERT INTO guilds (name, external_id)
VALUES (
    $1,
    $2
)
RETURNING *;

-- name: GetGuild :one
SELECT * FROM guilds WHERE external_id = $1;
