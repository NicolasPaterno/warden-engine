-- name: CreateAlert :one
INSERT INTO alerts (id, rule_id, room, message, severity, value, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING *;

-- name: GetAllAlerts :many
SELECT * FROM alerts ORDER BY created_at DESC;

-- name: GetAlertsByRoom :many
SELECT * FROM alerts WHERE room = $1 ORDER BY created_at DESC;

-- name: GetAlertsByRule :many
SELECT * FROM alerts WHERE rule_id = $1 ORDER BY created_at DESC;