CREATE TABLE IF NOT EXISTS users
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT, 
    email     TEXT NOT NULL UNIQUE,
    pass_hash BLOB NOT NULL                      
);

CREATE INDEX IF NOT EXISTS idx_email ON users (email);

ALTER TABLE users ADD COLUMN is_admin BOOLEAN NOT NULL DEFAULT FALSE;

CREATE TABLE IF NOT EXISTS apps
(
    id     INTEGER PRIMARY KEY,
    name   TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL UNIQUE
);

INSERT OR IGNORE INTO apps (id, name, secret) 
VALUES (1, 'test', 'test-secret');