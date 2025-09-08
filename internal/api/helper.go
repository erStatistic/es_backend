package rumiapi

import (
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kaeba0616/es_backend/internal/database"
)

func normalizeTier(s string) string {
	if s == "" || s == "All" {
		return ""
	}
	return s
}
func parseRFC3339Ptr(s string) (pgtype.Timestamptz, error) {
	var t pgtype.Timestamptz
	if s == "" {
		t.Valid = false
		return t, nil
	}
	v, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return t, err
	}
	t.Time = v
	t.Valid = true
	return t, nil
}

func parseIntDefault(s string, def int) int {
	if s == "" {
		return def
	}
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return def
}

func parseInt32List(q string) ([]int32, error) {
	if q == "" {
		return nil, nil
	}
	ps := strings.Split(q, ",")
	out := make([]int32, 0, len(ps))
	for _, p := range ps {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		v, err := strconv.ParseInt(p, 10, 32)
		if err != nil || v <= 0 {
			return nil, err
		}
		out = append(out, int32(v))
	}
	return out, nil
}

// pgtype.Timestamptz -> *time.Time (nullable 컬럼용)
func tsToPtr(ts pgtype.Timestamptz) *time.Time {
	if !ts.Valid {
		return nil
	}
	t := ts.Time
	return &t
}

// pgtype.Timestamptz -> time.Time (NOT NULL 컬럼용)
func tsToTime(ts pgtype.Timestamptz) time.Time {
	if ts.Valid {
		return ts.Time
	}
	// NOT NULL이어도 방어적으로 zero-time 반환
	return time.Time{}
}

func toTierDTO(t database.Tier) TierDTO {
	return TierDTO{
		ID:        t.ID,
		ImageUrl:  t.ImageUrl,
		Name:      t.Name,
		MmrMin:    t.MmrRange.Lower.Int32,
		MmrMax:    t.MmrRange.Upper.Int32,
		CreatedAt: tsToTime(t.CreatedAt), // time.Time
		UpdatedAt: tsToTime(t.UpdatedAt), // time.Time
	}
}

func toGameDTO(g database.Game) GameDTO {
	return GameDTO{
		ID:         g.ID,
		GameCode:   g.GameCode,
		StartedAt:  tsToPtr(g.StartedAt), // *time.Time
		AverageMmr: g.AverageMmr,
		CreatedAt:  tsToTime(g.CreatedAt), // time.Time
		UpdatedAt:  tsToTime(g.UpdatedAt), // time.Time
	}
}

func toGameTeamDTO(gt database.GameTeam) GameTeamDTO {
	return GameTeamDTO{
		ID:             gt.ID,
		GameID:         gt.GameID,
		TeamID:         gt.TeamID,
		GameRank:       gt.GameRank,
		TeamKills:      gt.TeamKills,
		MonsterCredits: gt.MonsterCredits,
		GainedMmr:      gt.GainedMmr,
		TeamAvgMmr:     gt.TeamAvgMmr,
		TotalTime:      gt.TotalTime,
		TimesID:        gt.TimesID, // *int32
		CreatedAt:      tsToTime(gt.CreatedAt),
		UpdatedAt:      tsToTime(gt.UpdatedAt),
	}
}

func Map[T any, U any](in []T, f func(T) U) []U {
	out := make([]U, len(in))
	for i, v := range in {
		out[i] = f(v)
	}
	return out
}

func toTimeDTO(row database.Time) TimeDTO {
	var start, end int32
	if row.TimeRange.Lower.Valid {
		start = row.TimeRange.Lower.Int32
	}
	if row.TimeRange.Upper.Valid {
		end = row.TimeRange.Upper.Int32
	}
	return TimeDTO{
		ID:        row.ID,
		No:        row.No,
		Name:      row.Name,
		Seconds:   row.Seconds,
		StartTime: start,
		EndTime:   end,
		CreatedAt: row.CreatedAt.Time, // pgtype.Timestamptz -> time.Time
		UpdatedAt: row.UpdatedAt.Time,
	}
}
func boundsFromRange(rr pgtype.Range[pgtype.Int4]) (start *int32, end *int32) {
	if rr.Lower.Valid {
		v := rr.Lower.Int32
		start = &v
	}
	if rr.Upper.Valid {
		v := rr.Upper.Int32
		end = &v
	}
	return
}

func makeInt4Range(min *int32, max *int32) pgtype.Range[pgtype.Int4] {
	r := pgtype.Range[pgtype.Int4]{Valid: true}

	if min != nil {
		r.Lower = pgtype.Int4{Int32: *min, Valid: true}
		r.LowerType = pgtype.Inclusive // 보통 하한 포함
	} else {
		r.LowerType = pgtype.Unbounded // 하한 없음
	}

	if max != nil {
		r.Upper = pgtype.Int4{Int32: *max, Valid: true}
		r.UpperType = pgtype.Exclusive // 정수 range는 상한 미포함이 관례
	} else {
		r.UpperType = pgtype.Unbounded // ★ 상한 없음(여기 없으면 네 에러)
	}

	return r
}
func parseTS(s string) (time.Time, error) {
	// 1) RFC3339(+09:00) 시도
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t, nil
	}
	// 2) +0900 형식 지원
	layouts := []string{
		"2006-01-02T15:04:05.000-0700",
		"2006-01-02T15:04:05-0700",
		"2006-01-02 15:04:05-0700",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	// 3) 오프셋 끝에 콜론 보정해서 RFC3339 재시도 (예: +0900 -> +09:00)
	s2 := addColonToOffset(s)
	return time.Parse(time.RFC3339Nano, s2)
}

func addColonToOffset(s string) string {
	// 끝이 +HHMM 또는 -HHMM 이면 +HH:MM 으로 바꿔줌
	n := len(s)
	if n >= 5 {
		off := s[n-5:]
		if (off[0] == '+' || off[0] == '-') && off[3] != ':' {
			return s[:n-2] + ":" + s[n-2:]
		}
	}
	return s
}
