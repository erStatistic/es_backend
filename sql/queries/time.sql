-- name: CreateTime :one
INSERT INTO
    times (no, name, seconds, time_range)
VALUES
    ($1, $2, $3, $4)
RETURNING
    *;

-- name: DeleteTime :exec
DELETE FROM times
WHERE
    id = $1;

-- name: GetTime :one
SELECT
    *
FROM
    times
WHERE
    id = $1;

-- name: ListTimes :many
SELECT
    *
FROM
    times
ORDER BY
    id ASC;

-- name: PatchTime :exec
UPDATE times
SET
    name = $2,
    seconds = $3,
    time_range = $4
WHERE
    id = $1;
