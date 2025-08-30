-- +goose Up
CREATE TABLE IF NOT EXISTS game_teams (
    id SERIAL PRIMARY KEY,
    game_id BIGINT NOT NULL REFERENCES games (game_code) ON DELETE CASCADE,
    team_id INT NOT NULL,
    game_rank INT NOT NULL,
    team_kills INT NOT NULL,
    monster_credits INT NOT NULL,
    gained_mmr INT NOT NULL,
    team_avg_mmr INT NOT NULL,
    total_time INT NOT NULL,
    times_id INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_game_teams UNIQUE (game_id, team_id),
    CONSTRAINT fk_game_teams_times FOREIGN KEY (times_id) REFERENCES times (id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_game_teams_game ON game_teams (game_id);

CREATE INDEX IF NOT EXISTS idx_game_teams_rank ON game_teams (game_id, game_rank);

CREATE INDEX IF NOT EXISTS idx_team_rank ON game_teams (game_rank);

CREATE INDEX IF NOT EXISTS idx_game_teams_times_id ON game_teams (times_id);

CREATE TRIGGER trg_game_teams_updated BEFORE
UPDATE ON game_teams FOR EACH ROW
EXECUTE FUNCTION set_updated_at ();

-- +goose Down
DROP TRIGGER IF EXISTS trg_game_teams_updated ON game_teams;

DROP INDEX IF EXISTS idx_game_teams_times_id;

DROP INDEX IF EXISTS idx_team_rank;

DROP INDEX IF EXISTS idx_game_teams_rank;

DROP INDEX IF EXISTS idx_game_teams_game;

DROP TABLE IF EXISTS game_teams;
