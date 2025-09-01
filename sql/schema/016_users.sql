-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    nickname VARCHAR(255) NOT NULL UNIQUE,
    user_num INT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_users_updated BEFORE
UPDATE ON users FOR EACH ROW
EXECUTE FUNCTION set_updated_at ();

-- +goose Down
DROP TRIGGER IF EXISTS trg_users_updated ON users;

DROP TABLE IF EXISTS users;
