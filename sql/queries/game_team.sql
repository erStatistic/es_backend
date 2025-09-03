-- name: CreateGameTeam :one
INSERT INTO
    game_teams (
        game_id,
        team_id,
        game_rank,
        team_kills,
        monster_credits,
        gained_mmr,
        team_avg_mmr,
        total_time
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
    *;

-- name: DeleteGameTeam :exec
DELETE FROM game_teams
WHERE
    id = $1;

-- name: GetGameTeam :one
SELECT
    *
FROM
    game_teams
WHERE
    game_id = $1
    AND team_id = $2;

-- name: GetGameTeamByGameID :many
SELECT
    *
FROM
    game_teams
WHERE
    game_id = $1;

-- name: GetListGameTeamsByGameRank :many
SELECT
    *
FROM
    game_teams
WHERE
    game_rank = $1;

-- name: ListGameTeams :many
SELECT
    *
FROM
    game_teams
ORDER BY
    id ASC;

-- name: PatchGameTeam :exec
UPDATE game_teams
SET
    game_id = $2,
    team_id = $3,
    game_rank = $4,
    team_kills = $5,
    monster_credits = $6,
    gained_mmr = $7,
    team_avg_mmr = $8,
    total_time = $9
WHERE
    id = $1;

-- name: TruncateGameTeams :exec
TRUNCATE TABLE game_teams;
