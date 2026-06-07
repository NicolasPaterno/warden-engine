-- +goose Up
CREATE TYPE alert_severity AS ENUM ('info', 'warning', 'critical');

CREATE TABLE alerts (
                        id          TEXT            PRIMARY KEY,
                        rule_id     TEXT            NOT NULL REFERENCES rules(id),
                        room        TEXT            NOT NULL,
                        message     TEXT            NOT NULL,
                        severity    alert_severity  NOT NULL,
                        value       FLOAT8          NOT NULL,
                        created_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);