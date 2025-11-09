-- name: CreateBookmark :one
INSERT INTO bookmarks (guild_id, channel_id, author, preview, user_id, created_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetBookmarksForUser :many
SELECT
    bookmarks.id,
    guilds.name AS guild_name,
    channels.name AS channel_name,
    bookmarks.author,
    bookmarks.preview,
    bookmarks.created_at
FROM bookmarks
JOIN guilds ON bookmarks.guild_id = guilds.id
JOIN channels ON bookmarks.channel_id = channels.id
WHERE bookmarks.user_id = $1;

-- name: GetBookmarksForUserByGuild :many
SELECT
    bookmarks.id,
    guilds.name AS guild_name,
    channels.name AS channel_name,
    bookmarks.author,
    bookmarks.preview,
    bookmarks.created_at
FROM bookmarks
JOIN guilds ON bookmarks.guild_id = guilds.id
JOIN channels ON bookmarks.channel_id = channels.id
WHERE bookmarks.user_id = $1 AND bookmarks.guild_id = $2;

-- name: GetBookmarksForUserByAuthor :many
SELECT
    bookmarks.id,
    guilds.name AS guild_name,
    channels.name AS channel_name,
    bookmarks.author,
    bookmarks.preview,
    bookmarks.created_at
FROM bookmarks
JOIN guilds ON bookmarks.guild_id = guilds.id
JOIN channels ON bookmarks.channel_id = channels.id
WHERE bookmarks.user_id = $1 AND bookmarks.author = $2;

-- name: DeleteBookmarkForUserByID :exec
DELETE FROM bookmarks WHERE user_id = $1 AND id = $2;

-- name: ResetBookmarksForUser :exec
DELETE FROM bookmarks WHERE user_id = $1;

-- name: ResetBookmarksForUserByGuild :exec
DELETE FROM bookmarks WHERE user_id = $1 AND guild_id = $2;

-- name: ResetBookmarksForUserByAuthor :exec
DELETE FROM bookmarks WHERE user_id = $1 AND author = $2;
