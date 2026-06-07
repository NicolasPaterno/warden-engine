-- +goose Down
DROP TABLE IF EXISTS rules;
DROP TYPE IF EXISTS action_type;
DROP TYPE IF EXISTS operator;
DROP TYPE IF EXISTS sensor_type;