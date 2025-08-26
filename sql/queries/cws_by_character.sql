-- name: ListCwByCharacterID :many
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
    JOIN weapons w ON w.code = cw.weapon_id
    LEFT JOIN positions p ON p.id = cw.position_id
WHERE
    cw.character_id = $1
ORDER BY
    w.name_kr,
    cw.id;

-- name: ListWeaponsByCharacterID :many
SELECT DISTINCT
    w.code AS code,
    w.name_kr AS name,
    w.image_url AS image_url
FROM
    character_weapons cw
    JOIN weapons w ON w.code = cw.weapon_id
WHERE
    cw.character_id = $1
ORDER BY
    w.name_kr;
