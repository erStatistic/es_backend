-- name: CreateUser :one
INSERT INTO
    users (id, user_num, name)
VALUES
    ($1, $2, $3)
RETURNING
    *;

-- name: GetUser :one
SELECT
    *
FROM
    users
WHERE
    name = $1;

-- name: DeleteUser :exec
DELETE FROM users;
