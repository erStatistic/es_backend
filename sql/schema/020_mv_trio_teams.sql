-- +goose Up
CREATE MATERIALIZED VIEW IF NOT EXISTS mv_trio_teams AS
WITH
    trio AS (
        SELECT
            gt.id AS game_team_id
        FROM
            game_team_cws gc
            JOIN game_teams gt ON gt.id = gc.game_team_id
        GROUP BY
            gt.id
        HAVING
            COUNT(*) = 3
    )
SELECT
    gt.id AS game_team_id,
    g.game_code,
    g.started_at,
    gt.game_rank,
    gt.total_time,
    gt.team_avg_mmr,
    ARRAY_AGG(
        gc.cw_id
        ORDER BY
            gc.cw_id
    ) AS cw_ids,
    ARRAY_AGG(
        cw.cluster_id
        ORDER BY
            cw.cluster_id
    ) AS cluster_ids
FROM
    trio t
    JOIN game_teams gt ON gt.id = t.game_team_id
    JOIN games g ON g.game_code = gt.game_id
    JOIN game_team_cws gc ON gc.game_team_id = gt.id
    JOIN character_weapons cw ON cw.id = gc.cw_id
GROUP BY
    gt.id,
    g.game_code,
    g.started_at,
    gt.game_rank,
    gt.total_time,
    gt.team_avg_mmr;

CREATE UNIQUE INDEX IF NOT EXISTS uq_mv_trio_teams_id ON mv_trio_teams (game_team_id);

CREATE INDEX IF NOT EXISTS idx_mv_trio_started_at ON mv_trio_teams (started_at);

CREATE INDEX IF NOT EXISTS idx_mv_trio_rank ON mv_trio_teams (game_rank);

CREATE INDEX IF NOT EXISTS idx_mv_trio_mmr ON mv_trio_teams (team_avg_mmr);

CREATE INDEX IF NOT EXISTS idx_mv_trio_cw_ids ON mv_trio_teams USING gin (cw_ids);

CREATE INDEX IF NOT EXISTS idx_mv_trio_cluster_ids ON mv_trio_teams USING gin (cluster_ids);

-- +goose Down
DROP INDEX IF EXISTS idx_mv_trio_cluster_ids;

DROP INDEX IF EXISTS idx_mv_trio_cw_ids;

DROP INDEX IF EXISTS idx_mv_trio_mmr;

DROP INDEX IF EXISTS idx_mv_trio_rank;

DROP INDEX IF EXISTS idx_mv_trio_started_at;

DROP INDEX IF EXISTS uq_mv_trio_teams_id;

DROP MATERIALIZED VIEW IF EXISTS mv_trio_teams;
