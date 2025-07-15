-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    user_num INT NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE users;
