-- name: CreateActivity :one
INSERT INTO activities (title, description, type, status, due_date, contact_id)
VALUES ($1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: GetActivity :one
SELECT * FROM activities WHERE id = $1;

-- name: ListActivities :many
SELECT * FROM activities
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateActivity :one
UPDATE activities
SET description=$2, status=$3, due_date=$4, updated_at=CURRENT_TIMESTAMP
WHERE id=$1
RETURNING *;

-- name: DeleteActivity :exec
DELETE FROM activities WHERE id = $1;
