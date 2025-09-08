-- +goose Up
CREATE TABLE IF NOT EXISTS characters (
    id SERIAL PRIMARY KEY,
    name_kr VARCHAR(255) NOT NULL,
    image_url_mini VARCHAR(255) NOT NULL DEFAULT '',
    image_url_full VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_characters_updated BEFORE
UPDATE ON characters FOR EACH ROW
EXECUTE FUNCTION set_updated_at ();

-- +goose Down
DROP TRIGGER IF EXISTS trg_characters_updated ON characters;

DROP TABLE IF EXISTS characters;
