-- +goose Up
CREATE TABLE IF NOT EXISTS character_weapons_stats (
    id SERIAL PRIMARY KEY,
    cw_id INT NOT NULL UNIQUE REFERENCES character_weapons (id) ON DELETE CASCADE,
    atk INT NOT NULL,
    def INT NOT NULL,
    cc INT NOT NULL,
    spd INT NOT NULL,
    sup INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_cws_range CHECK (
        atk BETWEEN 0 AND 5
        AND def BETWEEN 0 AND 5
        AND cc BETWEEN 0 AND 5
        AND spd BETWEEN 0 AND 5
        AND sup BETWEEN 0 AND 5
    )
);

CREATE INDEX IF NOT EXISTS idx_cws_cw ON character_weapons_stats (cw_id);

CREATE TRIGGER trg_cws_stats_updated BEFORE
UPDATE ON character_weapons_stats FOR EACH ROW
EXECUTE FUNCTION set_updated_at ();

-- +goose Down
DROP TRIGGER IF EXISTS trg_character_weapons_stats_updated ON character_weapons_stats;

DROP INDEX IF EXISTS idx_cws_cw;

DROP TABLE IF EXISTS character_weapons_stats;
