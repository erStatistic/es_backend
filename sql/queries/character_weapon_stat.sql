-- name: CreateCharacterWeaponStat :one
INSERT INTO
    character_weapon_stats (cw_id, atk, def, cc, spd, sup)
VALUES
    ($1, $2, $3, $4, $5, $6)
RETURNING
    *;

-- name: DeleteCharacterWeaponStat :exec
DELETE FROM character_weapon_stats
WHERE
    cw_id = $1;

-- name: GetCharacterWeaponStat :one
SELECT
    *
FROM
    character_weapon_stats
WHERE
    cw_id = $1;

-- name: ListCharacterWeaponStats :many
SELECT
    *
FROM
    character_weapon_stats
ORDER BY
    cw_id ASC;

-- name: PatchCharacterWeaponStat :exec
UPDATE character_weapon_stats
SET
    atk = $2,
    def = $3,
    cc = $4,
    spd = $5,
    sup = $6
WHERE
    cw_id = $1;
