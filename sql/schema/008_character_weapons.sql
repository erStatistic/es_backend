-- +goose Up
CREATE TABLE IF NOT EXISTS character_weapons (
    id SERIAL PRIMARY KEY,
    character_id INT NOT NULL REFERENCES characters (id) ON DELETE CASCADE,
    weapon_id INT NOT NULL REFERENCES weapons (code) ON DELETE CASCADE,
    position_id INT NOT NULL REFERENCES positions (id) ON DELETE CASCADE,
    cluster_id INT NOT NULL REFERENCES clusters (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_character_weapons_combo UNIQUE (character_id, weapon_id)
);

CREATE INDEX IF NOT EXISTS idx_cw_char ON character_weapons (character_id);

CREATE INDEX IF NOT EXISTS idx_cw_weapon ON character_weapons (weapon_id);

CREATE INDEX IF NOT EXISTS idx_cw_position ON character_weapons (position_id);

CREATE INDEX IF NOT EXISTS idx_cw_cluster ON character_weapons (cluster_id);

CREATE TRIGGER trg_cw_updated BEFORE
UPDATE ON character_weapons FOR EACH ROW
EXECUTE FUNCTION set_updated_at ();

-- +goose Down
DROP TRIGGER IF EXISTS trg_cw_updated ON character_weapons;

DROP INDEX IF EXISTS idx_cw_cluster;

DROP INDEX IF EXISTS idx_cw_position;

DROP INDEX IF EXISTS idx_cw_weapon;

DROP INDEX IF EXISTS idx_cw_char;

DROP TABLE IF EXISTS character_weapons;
