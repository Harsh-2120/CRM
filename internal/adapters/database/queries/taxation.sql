-- name: CreateTaxationDetail :one
INSERT INTO taxation_details (country, tax_type, rate, description)
VALUES ($1,$2,$3,$4)
RETURNING *;

-- name: GetTaxationDetail :one
SELECT * FROM taxation_details WHERE id = $1;

-- name: ListTaxationDetails :many
SELECT * FROM taxation_details
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateTaxationDetail :one
UPDATE taxation_details
SET tax_type=$2, rate=$3, description=$4, updated_at=CURRENT_TIMESTAMP
WHERE id=$1
RETURNING *;

-- name: DeleteTaxationDetail :exec
DELETE FROM taxation_details WHERE id = $1;
