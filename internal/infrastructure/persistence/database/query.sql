-- name: CreateUser :one
INSERT INTO users (id, channel, channel_user_id, created_at)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: GetUserByChannel :one
SELECT * FROM users
WHERE channel = ? AND channel_user_id = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;

-- name: CreateSession :one
INSERT INTO sessions (id, user_id, created_at, updated_at)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetSessionByID :one
SELECT * FROM sessions
WHERE id = ? LIMIT 1;

-- name: GetSessionsByUserID :many
SELECT * FROM sessions
WHERE user_id = ?
ORDER BY created_at DESC;

-- name: UpdateSession :one
UPDATE sessions
SET updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE id = ?;

-- name: CreateMessage :one
INSERT INTO messages (id, session_id, role, content, created_at)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetMessageByID :one
SELECT * FROM messages
WHERE id = ? LIMIT 1;

-- name: GetMessagesBySessionID :many
SELECT * FROM messages
WHERE session_id = ?
ORDER BY created_at ASC;

-- name: DeleteMessage :exec
DELETE FROM messages WHERE id = ?;

-- name: CreateTask :one
INSERT INTO tasks (id, session_id, skill, input, status, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetTaskByID :one
SELECT * FROM tasks
WHERE id = ? LIMIT 1;

-- name: GetTasksBySessionID :many
SELECT * FROM tasks
WHERE session_id = ?
ORDER BY created_at DESC;

-- name: UpdateTask :one
UPDATE tasks
SET output = ?, status = ?, error = ?, updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = ?;

-- name: CreateSkill :one
INSERT INTO skills (id, name, version, location, permissions, metadata, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetSkillByID :one
SELECT * FROM skills
WHERE id = ? LIMIT 1;

-- name: GetSkillByName :one
SELECT * FROM skills
WHERE name = ? LIMIT 1;

-- name: ListSkills :many
SELECT * FROM skills
ORDER BY created_at DESC;

-- name: UpdateSkill :one
UPDATE skills
SET version = ?, location = ?, permissions = ?, metadata = ?
WHERE id = ?
RETURNING *;

-- name: DeleteSkill :exec
DELETE FROM skills WHERE id = ?;

-- name: CreateSchedule :one
INSERT INTO schedules (id, skill, cron_expression, input, enabled, created_at)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetScheduleByID :one
SELECT * FROM schedules
WHERE id = ? LIMIT 1;

-- name: GetSchedulesBySkill :many
SELECT * FROM schedules
WHERE skill = ?
ORDER BY created_at DESC;

-- name: ListSchedules :many
SELECT * FROM schedules
ORDER BY created_at DESC;

-- name: UpdateSchedule :one
UPDATE schedules
SET cron_expression = ?, input = ?, enabled = ?
WHERE id = ?
RETURNING *;

-- name: DeleteSchedule :exec
DELETE FROM schedules WHERE id = ?;

-- name: CreateLog :one
INSERT INTO logs (id, level, source, message, metadata, created_at)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetLogByID :one
SELECT * FROM logs
WHERE id = ? LIMIT 1;

-- name: GetLogsByLevel :many
SELECT * FROM logs
WHERE level = ?
ORDER BY created_at DESC
LIMIT ?;

-- name: GetLogsBySource :many
SELECT * FROM logs
WHERE source = ?
ORDER BY created_at DESC
LIMIT ?;

-- name: GetLogsByDateRange :many
SELECT * FROM logs
WHERE created_at >= ? AND created_at <= ?
ORDER BY created_at DESC
LIMIT ?;

-- name: DeleteLog :exec
DELETE FROM logs WHERE id = ?;

-- name: DeleteLogsOlderThan :exec
DELETE FROM logs WHERE created_at < ?;
