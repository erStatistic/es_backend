-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_updated_at () RETURNS trigger AS $$
BEGIN
  NEW.updated_at := NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS set_updated_at ();

-- +goose StatementEnd
