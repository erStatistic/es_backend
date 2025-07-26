-- +goose Up
CREATE TABLE characters (
    id UUID PRIMARY KEY,
    characater_id INT NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL UNIQUE,
    thumbnail VARCHAR(255) NOT NULL,
    skill_group INT NOT NULL
);

-- +goose Down
DROP TABLE characters;
