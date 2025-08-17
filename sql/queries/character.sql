-- name: CreateCharacter :one
INSERT INTO
    characters (code, image_url, name_kr)
VALUES
    ($1, $2, $3)
RETURNING
    *;

-- name: ListCharacters :many
SELECT
    *
FROM
    characters
ORDER BY
    code ASC;

-- name: GetCharacter :one
SELECT
    *
FROM
    characters
WHERE
    code = $1;

-- name: DeleteCharacter :exec
DELETE FROM characters
WHERE
    code = $1;

-- name: PatchCharacter :exec
UPDATE characters
SET
    image_url = $2,
    name_kr = $3
WHERE
    code = $1;
