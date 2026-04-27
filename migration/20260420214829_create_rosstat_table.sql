-- +goose Up
CREATE SCHEMA IF NOT EXISTS midray;

CREATE TABLE IF NOT EXISTS midray.rosstat (
    id                      SERIAL          PRIMARY KEY,
    code                    INTEGER         NOT NULL,
    year                    INTEGER         NOT NULL,

    population_amout        INTEGER         NULL,
    birth_amount            INTEGER         NULL,
    death_amount            INTEGER         NULL,
    arrival_amount          INTEGER         NULL,
    departure_amount        INTEGER         NULL,
    male_amount             INTEGER         NULL,
    female_amount           INTEGER         NULL,

    land_area               INTEGER         NULL,  
    avg_salary              DECIMAL(10,2)   NULL,     
    medical_facilities      INTEGER         NULL,  
    schools_count           INTEGER         NULL,       
    housing_commissioned    INTEGER         NULL,
    UNIQUE(code, year)
);

ALTER TABLE         midray.rosstat 
ADD CONSTRAINT      fk_rosstat_geo 
FOREIGN KEY         (code) 
REFERENCES          midray.geo(code) 
ON DELETE CASCADE;

-- +goose Down
ALTER TABLE             midray.rosstat  DROP CONSTRAINT IF EXISTS fk_rosstat_geo;
DROP TABLE IF EXISTS    midray.rosstat;
