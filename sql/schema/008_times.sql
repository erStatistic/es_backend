-- +goose Up
CREATE TABLE times (
    id SERIAL PRIMARY KEY,
    day_num INT NOT NULL,
    seconds INT NOT NULL,
    is_daytime BOOL NOT NULL UNIQUE,
    start_time INT NOT NULL,
    end_time INT NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE times;
