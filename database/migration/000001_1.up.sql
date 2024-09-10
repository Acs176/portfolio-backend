CREATE SCHEMA IF NOT EXISTS testdb_schema;

CREATE TABLE IF NOT EXISTS testdb_schema.experience (
    id text PRIMARY KEY,
    company_name text NOT NULL,
    position text NOT NULL,
    start text NOT NULL,
    end text,
    description text NOT NULL
);