-- name: CreateTask :one
INSERT INTO tasks (title, description, status, priority, due_date, activity_id)
VALUES ($1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: GetTask :one
SELECT * FROM tasks WHERE id = $1;

-- name: ListTasks :many

SELECT *
FROM tasks
WHERE activity_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- name: UpdateTask :one
UPDATE tasks
SET description=$2, status=$3, priority=$4, due_date=$5, updated_at=CURRENT_TIMESTAMP
WHERE id=$1
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = $1;
