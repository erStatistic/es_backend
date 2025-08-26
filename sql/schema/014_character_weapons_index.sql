-- +goose Up
CREATE INDEX IF NOT EXISTS idx_cw_cluster_pos ON character_weapons (cluster_id, position_id);

CREATE INDEX IF NOT EXISTS idx_cw_cluster_char ON character_weapons (cluster_id, character_id);

-- +goose Down
DROP INDEX IF EXISTS idx_cw_cluster_pos;

DROP INDEX IF EXISTS idx_cw_cluster_char;
