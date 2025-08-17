-- +goose Up
CREATE TABLE IF NOT EXISTS tiers (
    id SERIAL PRIMARY KEY,
    imageUrl varchar(255) NOT NULL DEFAULT '',
    name VARCHAR(50) NOT NULL UNIQUE,
    mmr INT NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TRIGGER trg_tiers_updated BEFORE
UPDATE ON tiers FOR EACH ROW
EXECUTE FUNCTION set_updated_at ();

-- +goose Down
DROP TRIGGER IF EXISTS trg_tiers_updated ON tiers;

DROP TABLE IF EXISTS tiers;
