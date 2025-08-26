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

-- name: ListCwDirectoryByCluster :many
WITH
    cnts AS (
        SELECT
            cluster_id,
            COUNT(*) AS cws,
            COUNT(DISTINCT character_id) AS characters
        FROM
            character_weapons
        GROUP BY
            cluster_id
    ),
    bypos AS (
        SELECT
            cluster_id,
            position_id,
            COUNT(*) AS n
        FROM
            character_weapons
        GROUP BY
            cluster_id,
            position_id
    ),
    maxn AS (
        SELECT
            cluster_id,
            MAX(n) AS max_n
        FROM
            bypos
        GROUP BY
            cluster_id
    ),
    pick AS (
        SELECT
            b.cluster_id,
            MIN(b.position_id) AS position_id
        FROM
            bypos b
            JOIN maxn m ON m.cluster_id = b.cluster_id
            AND m.max_n = b.n
        GROUP BY
            b.cluster_id
    )
SELECT
    c.id AS cluster_id,
    c.name AS label,
    p.name AS role,
    cnts.cws,
    cnts.characters
FROM
    clusters c
    JOIN cnts ON cnts.cluster_id = c.id
    LEFT JOIN pick pk ON pk.cluster_id = c.id
    LEFT JOIN positions p ON p.id = pk.position_id
ORDER BY
    c.name;

-- name: ListCwByClusterID :many
SELECT
    cw.id AS cw_id,
    ch.id AS ch_id,
    ch.name_kr AS ch_name,
    ch.image_url_mini AS ch_img_mini,
    ch.image_url_full AS ch_img_full,
    w.code AS w_code,
    w.name_kr AS w_name,
    w.image_url AS w_img,
    p.id AS p_id,
    p.name AS p_name
FROM
    character_weapons cw
    JOIN characters ch ON ch.id = cw.character_id
    JOIN weapons w ON w.code = cw.weapon_id -- ⭐ 여기!
    JOIN positions p ON p.id = cw.position_id -- NOT NULL이면 LEFT 필요 없음
WHERE
    cw.cluster_id = $1
ORDER BY
    ch.name_kr,
    w.name_kr;

-- db/queries/character_weapons.sql
-- 캐릭터가 가진 CW(무기군) 목록
-- name: ListCwsByCharacterID :many
SELECT
    cw.id AS cw_id,
    w.code AS weapon_code,
    w.name_kr AS weapon_name,
    COALESCE(w.image_url, '') AS weapon_image_url
FROM
    character_weapons cw
    JOIN weapons w ON w.code = cw.weapon_id
WHERE
    cw.character_id = $1
ORDER BY
    w.name_kr;

-- overview에 표시할 신원(캐릭/무기/포지션)만 가져오는 쿼리
-- name: GetCwIdentity :one
SELECT
    cw.id AS cw_id,
    c.id AS ch_id,
    c.name_kr AS ch_name,
    COALESCE(c.image_url_mini, COALESCE(c.image_url_full, '')) AS ch_img,
    w.code AS w_code,
    w.name_kr AS w_name,
    COALESCE(w.image_url, '') AS w_img,
    p.id AS p_id,
    p.name AS p_name
FROM
    character_weapons cw
    JOIN characters c ON c.id = cw.character_id
    JOIN weapons w ON w.code = cw.weapon_id
    JOIN positions p ON p.id = cw.position_id
WHERE
    cw.id = $1;
