CREATE SCHEMA IF NOT EXISTS testdb_schema;

CREATE TABLE IF NOT EXISTS testdb_schema.experience (
    id uuid PRIMARY KEY,
    company_name text NOT NULL,
    position text NOT NULL,
    period_start text NOT NULL,
    period_end text,
    role_description text NOT NULL
);