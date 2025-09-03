-- name: CreateGameTeamCW :one
INSERT INTO
    game_team_cws (game_team_id, cw_id, mmr)
VALUES
    ($1, $2, $3)
RETURNING
    *;

-- name: DeleteGameTeamCW :exec
DELETE FROM game_team_cws
WHERE
    id = $1;

-- name: GetGameTeamCW :one
SELECT
    *
FROM
    game_team_cws
WHERE
    id = $1;

-- name: ListGameTeamCWs :many
SELECT
    *
FROM
    game_team_cws
ORDER BY
    id;

-- name: ListGameSameTeamCWs :many
SELECT
    *
FROM
    game_team_cws
WHERE
    game_team_id = $1;

-- name: PatchGameTeamCW :exec
UPDATE game_team_cws
SET
    game_team_id = $1,
    cw_id = $2,
    mmr = $3
WHERE
    id = $1;

-- name: TruncateGameTeamCWs :exec
TRUNCATE TABLE game_team_cws;
