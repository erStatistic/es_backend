-- name: CreateCharacter :one
INSERT INTO
    characters (
        id,
        imageUrl,
        name_kr,
        name_en,
        updated_at,
        created_at
    )
VALUES
    ($1, $2, $3, $4, NOW(), NOW())
RETURNING
    *;

-- name: GetCharacter :one
SELECT
    *
FROM
    characters
WHERE
    name_KR = $1;

-- name: DeleteCharacterById :exec
DELETE FROM characters
WHERE
    id = $1;

-- name: DeleteCharacterByName :exec
DELETE FROM characters
WHERE
    name_KR = $1;

-- name: PatchCharacter :exec
UPDATE characters
SET
    imageurl = $2,
    name_kr = $3,
    name_en = $4,
    updated_at = NOW()
WHERE
    id = $1;
