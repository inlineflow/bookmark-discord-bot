-- name: CreateBookmark :one
INSERT INTO bookmarks (guild_id, channel_id, author, preview, user_id, created_at)
VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?
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
WHERE bookmarks.user_id = ?;

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
WHERE bookmarks.user_id = ? AND bookmarks.guild_id = ?;

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
WHERE bookmarks.user_id = ? AND bookmarks.author = ?;

-- name: DeleteBookmarkForUserByID :exec
DELETE FROM bookmarks WHERE user_id = ? AND id = ?;

-- name: ResetBookmarksForUser :exec
DELETE FROM bookmarks WHERE user_id = ?;

-- name: ResetBookmarksForUserByGuild :exec
DELETE FROM bookmarks WHERE user_id = ? AND guild_id = ?;

-- name: ResetBookmarksForUserByAuthor :exec
DELETE FROM bookmarks WHERE user_id = ? AND author = ?;
