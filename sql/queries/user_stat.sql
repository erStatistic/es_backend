-- name: GetUserStat :one
SELECT
    *
FROM
    user_stats
WHERE
    id = $1;

-- name: GetUserStatbyUserId :many
SELECT
    *
FROM
    user_stats
WHERE
    user_id = $1;

-- name: CreateUserStat :one
INSERT INTO
    user_stats (user_id, character_id)
VALUES
    ($1, $2)
RETURNING
    *;

-- name: PatchUserStat :exec
UPDATE user_stats
SET
    user_id = $1,
    character_id = $2
WHERE
    id = $1;

-- name: DeleteUserStat :exec
DELETE FROM user_stats
WHERE
    user_id = $1;

-- name: ListUserStat :many
SELECT
    *
FROM
    user_stats;
