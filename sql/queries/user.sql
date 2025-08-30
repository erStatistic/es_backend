-- name: GetUser :one
SELECT
    *
FROM
    users
WHERE
    id = $1;

-- name: GetUserByNickname :one
SELECT
    *
FROM
    users
WHERE
    nickname = $1;

-- name: GetUserByUserNum :one
SELECT
    *
FROM
    users
WHERE
    user_num = $1;

-- name: CreateUser :one
INSERT INTO
    users (nickname, user_num)
VALUES
    ($1, $2)
RETURNING
    *;

-- name: PatchUser :exec
UPDATE users
SET
    nickname = $1,
    user_num = $2
WHERE
    id = $3;

-- name: DeleteUser :exec
DELETE FROM users
WHERE
    id = $1;

-- name: ListUsers :many
SELECT
    *
FROM
    users;
