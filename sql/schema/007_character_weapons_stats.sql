-- +goose Up
CREATE TABLE character_weapons_stats (
    id SERIAL PRIMARY KEY,
    cw_id INT NOT NULL REFERENCES character_weapons (id) ON DELETE CASCADE,
    atk INT NOT NULL,
    def INT NOT NULL,
    cc INT NOT NULL,
    spd INT NOT NULL,
    sup INT NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE character_weapons_stats;
