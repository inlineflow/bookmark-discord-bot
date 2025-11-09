-- +goose Up
CREATE TABLE guilds (
    id INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    external_id INT UNIQUE NOT NULL
);

CREATE TABLE channels (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    external_id INT UNIQUE NOT NULL
);

CREATE TABLE bookmarks (
    id INTEGER PRIMARY KEY,
    guild_id INTEGER NOT NULL REFERENCES guilds(id),
    channel_id INTEGER NOT NULL REFERENCES channels(id),
    author TEXT NOT NULL,
    preview TEXT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE guilds;
DROP TABLE channels;
DROP TABLE bookmarks;
