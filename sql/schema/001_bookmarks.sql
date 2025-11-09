-- +goose Up
CREATE TABLE guilds (
    id INT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE channels (
    id INT PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE bookmarks (
    id INT PRIMARY KEY,
    guild_id INT NOT NULL REFERENCES guilds(id),
    channel_id INT NOT NULL REFERENCES channels(id),
    author TEXT NOT NULL,
    preview TEXT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE guilds;
DROP TABLE channels;
DROP TABLE bookmarks;
