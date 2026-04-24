-- +goose Up
CREATE SCHEMA IF NOT EXISTS midray;

CREATE TABLE IF NOT EXISTS midray.geo (
    code        INTEGER     PRIMARY KEY,
    parent_code INTEGER     NULL,
    name        TEXT        NOT NULL,
    level       INTEGER     NOT NULL
);

ALTER TABLE         midray.geo           
ADD CONSTRAINT      fk_geo_parent   
FOREIGN KEY         (parent_code)   
REFERENCES          midray.geo(code)      
ON DELETE SET NULL;         

-- +goose Down
ALTER TABLE             midray.geo     DROP CONSTRAINT IF EXISTS fk_geo_parent;
DROP TABLE IF EXISTS    midray.geo;
