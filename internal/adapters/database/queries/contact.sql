-- name: CreateContact :one
INSERT INTO contacts (
    contact_type, first_name, last_name, company_name, company_id, email, phone,
    address, city, state, country, zipcode, position, social_media_profiles, notes, taxation_detail_id
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
RETURNING *;

-- name: GetContact :one
SELECT * FROM contacts WHERE id = $1;

-- name: ListContacts :many
SELECT * FROM contacts
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateContact :one
UPDATE contacts
SET first_name=$2, last_name=$3, phone=$4, address=$5, city=$6, state=$7, country=$8, zipcode=$9,
    position=$10, social_media_profiles=$11, notes=$12, updated_at=CURRENT_TIMESTAMP
WHERE id=$1
RETURNING *;

-- name: DeleteContact :exec
DELETE FROM contacts WHERE id = $1;
