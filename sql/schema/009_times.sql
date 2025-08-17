-- +goose Up
CREATE TABLE IF NOT EXISTS times (
    id SERIAL PRIMARY KEY,
    no INT NOT NULL,
    name VARCHAR(50) NOT NULL,
    seconds INT NOT NULL,
    start_time INT NOT NULL,
    end_time INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_times_span UNIQUE (no, start_time, end_time)
);

CREATE TRIGGER trg_times_updated BEFORE
UPDATE ON times FOR EACH ROW
EXECUTE FUNCTION set_updated_at ();

CREATE INDEX IF NOT EXISTS idx_times_span ON times (no, start_time, end_time);

-- +goose Down
DROP INDEX IF EXISTS idx_times_span;

DROP TABLE IF EXISTS times;
