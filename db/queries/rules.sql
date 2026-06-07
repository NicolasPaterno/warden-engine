-- name: CreateRule :one
INSERT INTO rules (id, name, room, sensor_type, operator, threshold, action_type, payload, enabled, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    RETURNING *;

-- name: GetRuleByID :one
SELECT * FROM rules WHERE id = $1;

-- name: GetAllRules :many
SELECT * FROM rules ORDER BY created_at DESC;

-- name: GetEnabledRules :many
SELECT * FROM rules WHERE enabled = true ORDER BY created_at DESC;

-- name: UpdateRule :one
UPDATE rules
SET name = $2, room = $3, sensor_type = $4, operator = $5,
    threshold = $6, action_type = $7, payload = $8, enabled = $9
WHERE id = $1
    RETURNING *;

-- name: DeleteRule :exec
DELETE FROM rules WHERE id = $1;