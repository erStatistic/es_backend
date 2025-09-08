-- name: CreateTier :one
INSERT INTO
    tiers (image_url, name, mmr_range, rank)
VALUES
    ($1, $2, $3, $4)
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
    mmr_range = $4,
    rank = $5
WHERE
    id = $1;
