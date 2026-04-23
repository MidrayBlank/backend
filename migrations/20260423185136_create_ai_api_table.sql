-- +goose Up
CREATE SCHEMA IF NOT EXISTS midray;

CREATE TABLE IF NOT EXISTS midray.ai_api (
    hash        TEXT        PRIMARY KEY,
    requests    INTEGER     NOT NULL,   
);

-- +goose Down
DROP TABLE IF EXISTS midray.ai_api;
