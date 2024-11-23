-- name: GetSchedule :one
SELECT * from schedules WHERE id = $1;

-- name: GetAllSchedules :many
SELECT * from schedules;

-- name: GetWorker :one
SELECT * from workers WHERE id = $1;

-- name: GetAllWorkers :many
SELECT * from workers;

-- name: GetTask :one
SELECT * from tasks WHERE id = $1;

-- name: GetAllTasks :many
SELECT * from tasks;

-- name: GetLocation :one
SELECT * from locations WHERE id = $1;

-- name: GetAllLocations :many
SELECT * from locations;

-- name: GetConstraint :one
SELECT * from constraints WHERE location_id = $1 AND task_id = $2 AND worker_id = $3 AND start_slot = $4 AND end_slot = $5 AND kind = $6;

-- name: GetAllConstraintsForTask :many
SELECT * from constraints WHERE task_id = $1;

-- name: GetAllConstraintsForLocation :many
SELECT * from constraints WHERE location_id = $1;

-- name: GetAllConstraintsForWorker :many
SELECT * from constraints WHERE worker_id = $1;

-- name: InsertSchedule :exec
INSERT INTO schedules (id, title) VALUES ($1, $2);

-- name: InsertConstraint :exec
INSERT INTO constraints (schedule_id, location_id, task_id, worker_id, start_slot, end_slot, kind) VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: InsertConstraints :copyfrom
INSERT INTO constraints (schedule_id, location_id, task_id, worker_id, start_slot, end_slot, kind) VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: InsertWorker :exec
INSERT INTO workers (id, first_name, last_name, schedule_id) VALUES ($1, $2, $3, $4);

-- name: InsertTask :exec
INSERT INTO tasks (id, title, story, schedule_id) VALUES ($1, $2, $3, $4);

-- name: InsertLocation :exec
INSERT INTO locations (id, title, story, schedule_id) VALUES ($1, $2, $3, $4);

-- name: UpsertWorker :exec
INSERT INTO workers (id, first_name, last_name, schedule_id) VALUES ($1, $2, $3, $4)
ON CONFLICT (id) DO UPDATE SET first_name = $2, last_name = $3, schedule_id = $4;

-- name: UpsertTask :exec
INSERT INTO tasks (id, title, story, schedule_id) VALUES ($1, $2, $3, $4)
ON CONFLICT (id) DO UPDATE SET title = $2, story = $3, schedule_id = $4;

-- name: UpsertLocation :exec
INSERT INTO locations (id, title, story, schedule_id) VALUES ($1, $2, $3, $4)
ON CONFLICT (id) DO UPDATE SET title = $2, story = $3, schedule_id = $4;

-- name: DeleteWorker :exec
DELETE FROM workers WHERE id = $1;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = $1;

-- name: DeleteLocation :exec
DELETE FROM locations WHERE id = $1;