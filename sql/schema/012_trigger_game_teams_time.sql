-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION assign_times_id () RETURNS trigger AS $$
BEGIN
    SELECT id INTO NEW.times_id
    FROM times
    WHERE NEW.total_time BETWEEN start_time AND end_time
    LIMIT 1;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
CREATE TRIGGER trg_assign_times_id BEFORE INSERT
OR
UPDATE OF total_time ON game_teams FOR EACH ROW
EXECUTE FUNCTION assign_times_id ();

-- +goose Down
DROP TRIGGER IF EXISTS trg_assign_times_id ON game_teams;
