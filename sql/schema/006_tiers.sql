-- +goose Up
CREATE TABLE IF NOT EXISTS tiers (
    id SERIAL PRIMARY KEY,
    image_url TEXT NOT NULL DEFAULT '',
    name VARCHAR(50) NOT NULL UNIQUE,
    rank INT NOT NULL DEFAULT 0,
    mmr_range INT4RANGE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_tiers_updated BEFORE
UPDATE ON tiers FOR EACH ROW
EXECUTE FUNCTION set_updated_at ();

-- +goose Down
DROP TRIGGER IF EXISTS trg_tiers_updated ON tiers;

DROP TABLE IF EXISTS tiers;
