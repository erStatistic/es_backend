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
    id = $1;

-- name: GetCharacterWeapon :one
SELECT
    *
FROM
    character_weapons
WHERE
    id = $1;

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
    character_id = $2,
    weapon_id = $3,
    position_id = $4,
    cluster_id = $5
WHERE
    id = $1;
