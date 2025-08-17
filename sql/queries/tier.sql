-- name: CreateTier :one
INSERT INTO
    tiers (image_url, name, mmr)
VALUES
    ($1, $2, $3)
RETURNING
    *;

-- name: DeleteTier :exec
DELETE FROM tiers
WHERE
    id = $1;

-- name: GetTier :one
SELECT
    *
FROM
    tiers
WHERE
    id = $1;

-- name: ListTiers :many
SELECT
    *
FROM
    tiers
ORDER BY
    id ASC;

-- name: PatchTier :exec
UPDATE tiers
SET
    image_url = $2,
    name = $3,
    mmr = $4
WHERE
    id = $1;
