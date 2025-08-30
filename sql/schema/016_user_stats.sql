-- +goose Up
CREATE TABLE IF NOT EXISTS user_stats (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users (user_num),
    character_id INT NOT NULL REFERENCES characters (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_user_stats_updated BEFORE
UPDATE ON user_stats FOR EACH ROW
EXECUTE FUNCTION set_updated_at ();

-- +goose Down
DROP TRIGGER IF EXISTS trg_user_stats_updated ON user_stats;

DROP TABLE IF EXISTS user_stats;
