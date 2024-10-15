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