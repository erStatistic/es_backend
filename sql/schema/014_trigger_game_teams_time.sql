-- +goose Up
-- (권장) range 검색 가속용 GiST 인덱스도 함께 준비
CREATE EXTENSION IF NOT EXISTS btree_gist;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION assign_times_id () RETURNS trigger AS $$
DECLARE
  v_id INT;
BEGIN
  -- 이미 times_id가 있으면 그대로 둠
  IF NEW.times_id IS NOT NULL THEN
    RETURN NEW;
  END IF;

  -- time_range @> total_time 로 매핑
  SELECT t.id
    INTO v_id
  FROM times t
  WHERE t.time_range @> NEW.total_time
  ORDER BY lower(t.time_range) ASC      -- 겹침 시 더 이른 구간 우선
  LIMIT 1;

  NEW.times_id := v_id;                  -- 매칭 없으면 NULL 그대로

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- 기존 트리거 교체
DROP TRIGGER IF EXISTS trg_assign_times_id ON game_teams;

CREATE TRIGGER trg_assign_times_id BEFORE INSERT
OR
UPDATE OF total_time ON game_teams FOR EACH ROW
EXECUTE FUNCTION assign_times_id ();

-- (권장) btree 인덱스 대신 GiST 인덱스 사용
-- 이미 btree 인덱스를 만들어두셨다면 교체하세요.
DROP INDEX IF EXISTS idx_times_span;

CREATE INDEX IF NOT EXISTS idx_times_span_gist ON times USING gist (time_range);

-- +goose Down
-- 트리거/함수 제거
DROP TRIGGER IF EXISTS trg_assign_times_id ON game_teams;

DROP FUNCTION IF EXISTS assign_times_id ();

-- 인덱스 원복(원하시면 btree로 복구)
DROP INDEX IF EXISTS idx_times_span_gist;

CREATE INDEX IF NOT EXISTS idx_times_span ON times (time_range);
