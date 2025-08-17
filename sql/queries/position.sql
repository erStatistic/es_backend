-- name: CreatePosition :one
INSERT INTO
    positions (name, image_url)
VALUES
    ($1, $2)
RETURNING
    *;

-- name: DeletePosition :exec
DELETE FROM positions
WHERE
    id = $1;

-- name: GetPosition :one
SELECT
    *
FROM
    positions
WHERE
    id = $1;

-- name: ListPositions :many
SELECT
    *
FROM
    positions
ORDER BY
    id ASC;

-- name: PatchPosition :exec
UPDATE positions
SET
    image_url = $2,
    name = $3
WHERE
    id = $1;
