-- +goose Up

CREATE TABLE IF NOT EXISTS async_requests (
    id           SERIAL PRIMARY KEY,
    hash         TEXT UNIQUE,
    code         INTEGER,
    status       INTEGER,
    request_type INTEGER,
    attempts     INTEGER,
    request_data TEXT,
    result       TEXT,
    error        TEXT,
    deadline_at  TIMESTAMP,
    created_at   TIMESTAMP,
    updated_at   TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS async_requests;
