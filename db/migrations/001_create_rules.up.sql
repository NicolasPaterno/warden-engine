-- +goose Up
CREATE TYPE sensor_type AS ENUM ('temperature', 'humidity', 'motion', 'co2');
CREATE TYPE operator AS ENUM ('gt', 'lt', 'gte', 'lte', 'eq');
CREATE TYPE action_type AS ENUM ('alert');

CREATE TABLE rules (
    id          TEXT        PRIMARY KEY,
    name        TEXT        NOT NULL,
    room        TEXT        NOT NULL,
    sensor_type sensor_type NOT NULL,
    operator    operator    NOT NULL,
    threshold   FLOAT8      NOT NULL,
    action_type action_type NOT NULL,
    payload     TEXT        NOT NULL DEFAULT '',
    enabled     BOOLEAN     NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
