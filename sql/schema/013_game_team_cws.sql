-- +goose Up
CREATE TABLE IF NOT EXISTS game_team_cws (
    id SERIAL PRIMARY KEY,
    game_team_id INT NOT NULL REFERENCES game_teams (id) ON DELETE CASCADE,
    cw_id INT NOT NULL REFERENCES character_weapons (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_game_team_cws UNIQUE (game_team_id, cw_id)
);

CREATE INDEX IF NOT EXISTS idx_gtcws_game_team ON game_team_cws (game_team_id);

CREATE INDEX IF NOT EXISTS idx_gtcws_cw ON game_team_cws (cw_id);

CREATE TRIGGER trg_game_team_cws_updated BEFORE
UPDATE ON game_team_cws FOR EACH ROW
EXECUTE FUNCTION set_updated_at ();

-- +goose Down
DROP TRIGGER IF EXISTS trg_game_team_cws_updated ON game_team_cws;

DROP INDEX IF EXISTS idx_gtcws_cw;

DROP INDEX IF EXISTS idx_gtcws_game_team;

DROP TABLE IF EXISTS game_team_cws;
