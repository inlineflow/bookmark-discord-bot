-- name: CreateGuild :one
INSERT INTO guilds (name, external_id)
VALUES (
    ?,
    ?
)
RETURNING *;

-- name: GetGuild :one
SELECT * FROM guilds WHERE external_id = ?;
