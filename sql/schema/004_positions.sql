-- +goose Up
CREATE TABLE IF NOT EXISTS positions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    image_url VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_positions_updated BEFORE
UPDATE ON positions FOR EACH ROW
EXECUTE FUNCTION set_updated_at ();

-- +goose Down
DROP TRIGGER IF EXISTS trg_positions_updated ON positions;

DROP TABLE IF EXISTS positions;
