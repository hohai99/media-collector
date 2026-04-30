-- Media Collector database schema

CREATE TABLE IF NOT EXISTS config (
    key   TEXT PRIMARY KEY,
    value TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS collections (
    id        TEXT PRIMARY KEY,
    name      TEXT NOT NULL,
    path      TEXT NOT NULL UNIQUE,
    parent_id TEXT REFERENCES collections(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS media (
    id            TEXT PRIMARY KEY,
    name          TEXT NOT NULL,
    path          TEXT NOT NULL UNIQUE,
    type          TEXT NOT NULL,
    size          INTEGER NOT NULL DEFAULT 0,
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    collection_id TEXT REFERENCES collections(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS player_configs (
    id              TEXT PRIMARY KEY,
    name            TEXT NOT NULL,
    total_time      INTEGER NOT NULL DEFAULT 0,
    transition_time INTEGER NOT NULL DEFAULT 0,
    time_per_picture INTEGER NOT NULL DEFAULT 5,
    sound_source    TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS player_media (
    config_id TEXT NOT NULL REFERENCES player_configs(id) ON DELETE CASCADE,
    media_id  TEXT NOT NULL REFERENCES media(id) ON DELETE CASCADE,
    position  INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (config_id, media_id)
);

CREATE TABLE IF NOT EXISTS player_collections (
    config_id     TEXT NOT NULL REFERENCES player_configs(id) ON DELETE CASCADE,
    collection_id TEXT NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
    PRIMARY KEY (config_id, collection_id)
);
