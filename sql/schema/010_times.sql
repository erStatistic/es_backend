-- +goose Up
CREATE TABLE IF NOT EXISTS times (
    id SERIAL PRIMARY KEY,
    no INT NOT NULL,
    name VARCHAR(50) NOT NULL,
    seconds INT NOT NULL,
    time_range INT4RANGE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_times_span UNIQUE (no, time_range)
);

CREATE TRIGGER trg_times_updated BEFORE
UPDATE ON times FOR EACH ROW
EXECUTE FUNCTION set_updated_at ();

CREATE INDEX IF NOT EXISTS idx_times_span ON times (no, time_range);

-- +goose Down
DROP INDEX IF EXISTS idx_times_span;

DROP TABLE IF EXISTS times;
