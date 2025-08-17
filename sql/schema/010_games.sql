-- +goose Up
CREATE TABLE IF NOT EXISTS games (
    id SERIAL PRIMARY KEY,
    game_code BIGINT NOT NULL UNIQUE, -- team_info.GameID
    started_at TIMESTAMPTZ DEFAULT NULL,
    average_mmr INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_games_started_at_desc ON games (started_at DESC);

CREATE TRIGGER trg_games_updated BEFORE
UPDATE ON games FOR EACH ROW
EXECUTE FUNCTION set_updated_at ();

-- +goose Down
DROP TRIGGER IF EXISTS trg_games_updated ON games;

DROP INDEX IF EXISTS idx_games_started_at_desc;

DROP TABLE IF EXISTS games;
