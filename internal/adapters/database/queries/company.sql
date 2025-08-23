-- name: CreateCompany :one
INSERT INTO companies (
    name, industry, website, phone, email, address, city, state, country, zipcode, created_by, organization_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetCompany :one
SELECT * FROM companies WHERE id = $1;

-- name: ListCompanies :many
SELECT * FROM companies
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateCompany :one
UPDATE companies
SET name = $2, industry = $3, website = $4, phone = $5, email = $6, address = $7, city = $8, state = $9, country = $10,
    zipcode = $11, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteCompany :exec
DELETE FROM companies WHERE id = $1;
