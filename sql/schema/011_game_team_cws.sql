-- +goose Up
CREATE TABLE game_team_cws (
    id SERIAL PRIMARY KEY,
    game_team_id INT NOT NULL REFERENCES game_teams (id) ON DELETE CASCADE,
    cw_id INT NOT NULL REFERENCES character_weapons (id) ON DELETE CASCADE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_game_team_cw UNIQUE (game_team_id, cw_id)
);

-- +goose Down
DROP TABLE game_team_cws;
