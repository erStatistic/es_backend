-- name: CreateCharacter :one
INSERT INTO
    characters (image_url_mini, image_url_full, name_kr)
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
    id ASC;

-- name: GetCharacter :one
SELECT
    *
FROM
    characters
WHERE
    id = $1;

-- name: DeleteCharacter :exec
DELETE FROM characters
WHERE
    id = $1;

-- name: PatchCharacter :exec
UPDATE characters
SET
    image_url_mini = $2,
    image_url_full = $3,
    name_kr = $4
WHERE
    id = $1;
