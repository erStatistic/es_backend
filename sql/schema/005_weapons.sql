-- +goose Up
CREATE TABLE weapons (
    id SERIAL PRIMARY KEY,
    name_KR VARCHAR(255) NOT NULL UNIQUE,
    name_EN VARCHAR(255) NOT NULL UNIQUE,
    attack_speed NUMERIC(2, 2) NOT NULL,
    skill_amplifier NUMERIC(2, 2) NOT NULL,
    attack_amplifier NUMERIC(2, 2) NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE characters;
