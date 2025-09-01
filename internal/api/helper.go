package rumiapi

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kaeba0616/es_backend/internal/database"
)

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
