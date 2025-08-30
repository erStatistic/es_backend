-- +goose Up
CREATE TABLE IF NOT EXISTS skills (
    id SERIAL PRIMARY KEY,
    skill_id INT NOT NULL,
    character_id INT NOT NULL REFERENCES characters (id) ON DELETE CASCADE,
    name_kr VARCHAR(255) NOT NULL,
    image_url VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_skills_updated BEFORE
UPDATE ON skills FOR EACH ROW
EXECUTE FUNCTION set_updated_at ();

-- +goose Down
DROP TRIGGER IF EXISTS trg_skills_updated ON skills;

DROP TABLE IF EXISTS skills;
