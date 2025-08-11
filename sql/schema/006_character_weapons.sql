-- +goose Up
CREATE TABLE character_weapons (
    id SERIAL PRIMARY KEY,
    character_id INT NOT NULL REFERENCES characters (id) ON DELETE CASCADE,
    weapon_id INT NOT NULL REFERENCES weapons (id) ON DELETE CASCADE,
    position_id INT NOT NULL REFERENCES positions (id) ON DELETE CASCADE,
    cluster_id INT NOT NULL REFERENCES clusters (id) ON DELETE CASCADE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE character_weapons;
