-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION assign_times_id () RETURNS trigger AS $$
DECLARE
  v_id INT;
BEGIN
  IF NEW.times_id IS NOT NULL THEN
    RETURN NEW;
  END IF;

  SELECT t.id
    INTO v_id
  FROM times t
  WHERE NEW.total_time BETWEEN t.start_time AND t.end_time
  ORDER BY t.start_time ASC       -- 겹칠 때 선택 기준 명시
  LIMIT 1;

  NEW.times_id := v_id;           -- 매칭 없으면 NULL로 설정

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
DROP TRIGGER IF EXISTS trg_assign_times_id ON game_teams;

CREATE TRIGGER trg_assign_times_id BEFORE INSERT
OR
UPDATE OF total_time ON game_teams FOR EACH ROW
EXECUTE FUNCTION assign_times_id ();

-- +goose Down
DROP TRIGGER IF EXISTS trg_assign_times_id ON game_teams;

DROP FUNCTION IF EXISTS assign_times_id ();
