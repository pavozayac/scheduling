-- name: GetSchedule :one
SELECT * from schedules WHERE id = $1;

-- name: GetWorker :one
SELECT * from workers WHERE id = $1;