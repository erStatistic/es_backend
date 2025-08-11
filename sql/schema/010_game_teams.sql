-- +goose Up
CREATE TABLE game_teams (
    id SERIAL PRIMARY KEY,
    game_id INT REFERENCES games (id) ON DELETE CASCADE,
    team_id INT NOT NULL,
    game_rank INT NOT NULL,
    team_kills INT NOT NULL,
    moster_credits INT NOT NULL,
    gained_mmr INT NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_game_team UNIQUE (game_id, team_id)
);

-- +goose Down
DROP TABLE game_teams;
