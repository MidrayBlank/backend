-- +goose Up
CREATE TABLE IF NOT EXISTS midray.ai_report (
    code        INTEGER     PRIMARY KEY,
    report      TEXT        NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS midray.ai_report;
