-- +goose Up
CREATE TABLE IF NOT EXISTS user_routes (
    id SERIAL PRIMARY KEY,
    route_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    weapon_id INT NOT NULL REFERENCES weapons (code) ON DELETE CASCADE,
    character_id INT NOT NULL REFERENCES characters (id) ON DELETE CASCADE,
    count INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_user_routes_updated BEFORE
UPDATE ON user_routes FOR EACH ROW
EXECUTE FUNCTION set_updated_at ();

-- +goose Down
DROP TRIGGER IF EXISTS trg_user_routes_updated ON user_routes;

DROP TABLE IF EXISTS user_routes;
