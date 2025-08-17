-- name: CreateCharacterWeapon :one
INSERT INTO
    character_weapons (character_id, weapon_id, position_id, cluster_id)
VALUES
    ($1, $2, $3, $4)
RETURNING
    *;

-- name: DeleteCharacterWeapon :exec
DELETE FROM character_weapons
WHERE
    character_id = $1
    AND weapon_id = $2;

-- name: GetCharacterWeapon :one
SELECT
    *
FROM
    character_weapons
WHERE
    character_id = $1
    AND weapon_id = $2;

-- name: ListCharacterWeapons :many
SELECT
    *
FROM
    character_weapons
ORDER BY
    id ASC;

-- name: PatchCharacterWeapon :exec
UPDATE character_weapons
SET
    position_id = $3,
    cluster_id = $4
WHERE
    character_id = $1
    AND weapon_id = $2;
