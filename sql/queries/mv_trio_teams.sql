-- name: RefreshMvTrioTeams :exec
REFRESH MATERIALIZED VIEW CONCURRENTLY mv_trio_teams;

-- name: GetCwStats :many
WITH
    scope AS (
        SELECT
            *
        FROM
            mv_trio_teams
        WHERE
            (
                $1::timestamptz IS NULL
                OR started_at >= $1
            )
            AND (
                $2::timestamptz IS NULL
                OR started_at < $2
            )
    ),
    scope_with_tier AS (
        SELECT
            s.*,
            gt.gained_mmr
        FROM
            scope s
            JOIN game_teams gt ON gt.id = s.game_team_id
            LEFT JOIN tiers t ON s.team_avg_mmr <@ t.mmr_range
        WHERE
            (
                NULLIF($3, '') IS NULL
                OR t.name = $3
            )
    ),
    denom AS (
        SELECT
            COUNT(*)::float8 AS total_teams
        FROM
            scope_with_tier
    ),
    cw_hits AS (
        SELECT
            unnest(cw_ids) AS cw_id,
            COUNT(*) AS team_count,
            SUM((game_rank = 1)::int) AS wins,
            AVG(gained_mmr)::float AS avg_mmr, -- ✅ gained_mmr 평균
            AVG(total_time)::float AS avg_survival
        FROM
            scope_with_tier
        GROUP BY
            1
        HAVING
            COUNT(*) >= COALESCE(NULLIF($4::int, 0), 50)
    )
SELECT
    h.cw_id,
    cw.character_id,
    cw.weapon_id,
    cw.position_id,
    cw.cluster_id,
    h.team_count AS samples,
    h.wins,
    COALESCE(
        h.wins::float8 / NULLIF(h.team_count::float8, 0.0),
        0.0
    ) AS win_rate,
    COALESCE(
        h.team_count::float8 / NULLIF(d.total_teams, 0.0),
        0.0
    ) AS pick_rate,
    h.avg_mmr, -- ✅ gained_mmr 평균
    h.avg_survival
FROM
    cw_hits h
    JOIN character_weapons cw ON cw.id = h.cw_id
    CROSS JOIN denom d
ORDER BY
    win_rate DESC;

-- name: GetTopClusterCombos :many
WITH
    scope AS (
        SELECT
            *
        FROM
            mv_trio_teams
        WHERE
            (
                $1::timestamptz IS NULL
                OR started_at >= $1
            )
            AND (
                $2::timestamptz IS NULL
                OR started_at < $2
            )
    ),
    scope_with_tier AS (
        SELECT
            s.*,
            gt.gained_mmr
        FROM
            scope s
            JOIN game_teams gt ON gt.id = s.game_team_id
            LEFT JOIN tiers t ON s.team_avg_mmr <@ t.mmr_range
        WHERE
            (
                NULLIF($3, '') IS NULL
                OR t.name = $3
            )
    ),
    denom AS (
        SELECT
            COUNT(*)::float8 AS total_teams
        FROM
            scope_with_tier
    ),
    combo AS (
        SELECT
            cluster_ids,
            COUNT(*) AS team_count,
            SUM((game_rank = 1)::int) AS wins,
            AVG(gained_mmr)::float AS avg_mmr, -- ✅ gained_mmr 평균
            AVG(total_time)::float AS avg_survival
        FROM
            scope_with_tier
        GROUP BY
            cluster_ids
        HAVING
            COUNT(*) >= COALESCE(NULLIF($6::int, 0), 50)
    )
SELECT
    c.cluster_ids,
    ARRAY_TO_STRING(
        ARRAY(
            SELECT
                cu.name
            FROM
                UNNEST(c.cluster_ids) cid
                JOIN clusters cu ON cu.id = cid
            ORDER BY
                cu.name
        ),
        ' · '
    ) AS cluster_label,
    c.team_count AS samples,
    c.wins,
    COALESCE(
        c.wins::float8 / NULLIF(c.team_count::float8, 0.0),
        0.0
    ) AS win_rate,
    COALESCE(
        c.team_count::float8 / NULLIF(d.total_teams, 0.0),
        0.0
    ) AS pick_rate,
    c.avg_mmr, -- ✅ gained_mmr 평균
    c.avg_survival
FROM
    combo c
    CROSS JOIN denom d
ORDER BY
    win_rate DESC
LIMIT
    $4
OFFSET
    $5;

-- name: GetCompMetricsBySelectedCWs :one
WITH
    want AS (
        SELECT
            ARRAY(
                SELECT
                    cw.cluster_id
                FROM
                    character_weapons cw
                WHERE
                    cw.id = ANY ($3::INT[])
                ORDER BY
                    cw.cluster_id
            ) AS cluster_ids
    ),
    scope AS (
        SELECT
            *
        FROM
            mv_trio_teams
        WHERE
            (
                $1::timestamptz IS NULL
                OR started_at >= $1
            )
            AND (
                $2::timestamptz IS NULL
                OR started_at < $2
            )
    ),
    scope_with_tier AS (
        SELECT
            s.*,
            gt.gained_mmr
        FROM
            scope s
            JOIN game_teams gt ON gt.id = s.game_team_id
            LEFT JOIN tiers t ON s.team_avg_mmr <@ t.mmr_range
        WHERE
            (
                NULLIF($4, '') IS NULL
                OR t.name = $4
            )
    ),
    denom AS (
        SELECT
            COUNT(*)::float8 AS total_teams
        FROM
            scope_with_tier
    ),
    matched AS (
        SELECT
            s.*
        FROM
            scope_with_tier s,
            want w
        WHERE
            s.cluster_ids = w.cluster_ids
    ),
    agg AS (
        SELECT
            COUNT(*) AS team_count,
            SUM((game_rank = 1)::int) AS wins,
            AVG(gained_mmr)::float AS avg_mmr, -- ✅ gained_mmr 평균
            AVG(total_time)::float AS avg_survival
        FROM
            matched
    )
SELECT
    a.team_count AS samples,
    a.wins,
    COALESCE(
        a.wins::float8 / NULLIF(a.team_count::float8, 0.0),
        0.0
    ) AS win_rate,
    COALESCE(
        a.team_count::float8 / NULLIF(d.total_teams, 0.0),
        0.0
    ) AS pick_rate,
    a.avg_mmr, -- ✅ gained_mmr 평균
    a.avg_survival
FROM
    agg a
    CROSS JOIN denom d
WHERE
    a.team_count >= COALESCE(NULLIF($5::int, 0), 50);

-- name: GetCwDailyTrend :many
WITH
    scope AS (
        SELECT
            *
        FROM
            mv_trio_teams
        WHERE
            (
                started_at >= COALESCE($1::timestamptz, NOW() - INTERVAL '14 days')
                AND started_at < COALESCE($2::timestamptz, NOW())
            )
    ),
    scope_with_tier AS (
        SELECT
            s.*
        FROM
            scope s
            LEFT JOIN tiers t ON s.team_avg_mmr <@ t.mmr_range
        WHERE
            (
                NULLIF($3, '') IS NULL
                OR t.name = $3
            )
    ),
    daily AS (
        SELECT
            DATE_TRUNC('day', started_at) AS DAY,
            COUNT(*) AS team_count,
            SUM((game_rank = 1)::int) AS wins
        FROM
            scope_with_tier
        WHERE
            $4::int = ANY (cw_ids)
        GROUP BY
            1
        HAVING
            COUNT(*) >= COALESCE(NULLIF($5::int, 0), 50)
    ),
    totals AS (
        SELECT
            DATE_TRUNC('day', started_at) AS DAY,
            COUNT(*)::float8 AS total_teams
        FROM
            scope_with_tier
        GROUP BY
            1
    )
SELECT
    d.day,
    d.team_count AS samples,
    COALESCE(
        d.wins::float8 / NULLIF(d.team_count::float8, 0.0),
        0.0
    ) AS win_rate,
    COALESCE(
        d.team_count::float8 / NULLIF(t.total_teams, 0.0),
        0.0
    ) AS pick_rate
FROM
    daily d
    JOIN totals t ON t.day = d.day
ORDER BY
    d.day;
