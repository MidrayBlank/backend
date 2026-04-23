-- +goose Up
CREATE SCHEMA IF NOT EXISTS midray;

CREATE TABLE IF NOT EXISTS midray.rosstat_age (
    rosstat_id      INTEGER     NOT NULL,
    age             INTEGER     NOT NULL,
    male_amount     INTEGER         NULL,
    female_amount   INTEGER         NULL,
    PRIMARY KEY (rosstat_id, age)
);

ALTER TABLE         midray.rosstat_age 
ADD CONSTRAINT      fk_rosstat_age_rosstat 
FOREIGN KEY         (rosstat_id) 
REFERENCES          midray.rosstat(id) 
ON DELETE CASCADE;

-- +goose Down
ALTER TABLE             midray.rosstat_age      DROP CONSTRAINT IF EXISTS fk_rosstat_age_rosstat;
DROP TABLE IF EXISTS    midray.rosstat_age;
