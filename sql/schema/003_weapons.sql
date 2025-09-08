-- +goose Up
CREATE TABLE IF NOT EXISTS weapons (
    id SERIAL PRIMARY KEY,
    code INT NOT NULL UNIQUE,
    name_kr VARCHAR(255) NOT NULL,
    image_url VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_weapons_updated BEFORE
UPDATE ON weapons FOR EACH ROW
EXECUTE FUNCTION set_updated_at ();

-- +goose Down
DROP TRIGGER IF EXISTS trg_weapons_updated ON weapons;

DROP TABLE IF EXISTS weapons;
