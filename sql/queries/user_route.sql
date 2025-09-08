-- name: CreateUserRoute :one
INSERT INTO
    user_routes (route_id, title, weapon_id, character_id, count)
VALUES
    ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: DeleteUserRoute :exec
DELETE FROM user_routes
WHERE
    route_id = $1;

-- name: GetUserRoute :one
SELECT
    *
FROM
    user_routes
WHERE
    route_id = $1;

-- name: ListUserRoutes :many
SELECT
    *
FROM
    user_routes
ORDER BY
    route_id ASC;

-- name: PatchUserRoute :exec
UPDATE user_routes
SET
    weapon_id = $2,
    character_id = $3,
    title = $4,
    count = $5
WHERE
    route_id = $1;

-- name: ListCWRoutes :many
SELECT
    *
FROM
    user_routes
WHERE
    character_id = $1
    AND weapon_id = $2
ORDER BY
    count DESC
LIMIT
    3;
