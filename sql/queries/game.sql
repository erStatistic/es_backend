-- name: CreateGame :one
INSERT INTO
    games (game_code, average_mmr, started_at)
VALUES
    ($1, $2, $3)
RETURNING
    *;

-- name: GetGame :one
SELECT
    *
FROM
    games
WHERE
    game_code = $1;

-- name: ListGames :many
SELECT
    *
FROM
    games
ORDER BY
    started_at DESC;

-- name: PatchGame :exec
UPDATE games
SET
    average_mmr = $3,
    started_at = $2
WHERE
    game_code = $1;

-- name: DeleteGame :exec
DELETE FROM games
WHERE
    game_code = $1;

-- name: TruncateGames :exec
TRUNCATE TABLE games;
