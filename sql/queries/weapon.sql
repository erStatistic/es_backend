-- name: CreateWeapon :one
INSERT INTO
    weapons (code, name_kr, image_url)
VALUES
    ($1, $2, $3)
RETURNING
    *;

-- name: ListWeapons :many
SELECT
    *
FROM
    weapons
ORDER BY
    code ASC;

-- name: GetWeapon :one
SELECT
    *
FROM
    weapons
WHERE
    code = $1;

-- name: DeleteWeapon :exec
DELETE FROM weapons
WHERE
    code = $1;

-- name: PatchWeapon :exec
UPDATE weapons
SET
    image_url = $2,
    name_kr = $3
WHERE
    code = $1;
