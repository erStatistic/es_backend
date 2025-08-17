package rumiapi

import (
	"database/sql"
	"time"

	"github.com/kaeba0616/es_backend/internal/database"
)

func tptrTime(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	t := nt.Time
	return &t
}

func tptrInt32(nt sql.NullInt32) *int32 {
	if !nt.Valid {
		return nil
	}
	return &nt.Int32
}

func toGameDTO(g database.Game) GameDTO {
	return GameDTO{
		ID: g.ID, GameCode: g.GameCode,
		StartedAt:  tptrTime(g.StartedAt),
		AverageMmr: g.AverageMmr,
		CreatedAt:  g.CreatedAt,
		UpdatedAt:  g.UpdatedAt,
	}
}

func toGameTeamDTO(g database.GameTeam) GameTeamDTO {
	return GameTeamDTO{
		ID:             g.ID,
		GameID:         g.GameID,
		TeamID:         g.TeamID,
		GameRank:       g.GameRank,
		TeamKills:      g.TeamKills,
		MonsterCredits: g.MonsterCredits,
		GainedMmr:      g.GainedMmr,
		TeamAvgMmr:     g.TeamAvgMmr,
		TotalTime:      g.TotalTime,
		TimesID:        tptrInt32(g.TimesID),
		CreatedAt:      g.CreatedAt,
		UpdatedAt:      g.UpdatedAt,
	}
}
func Map[T any, U any](in []T, f func(T) U) []U {
	out := make([]U, len(in))
	for i, v := range in {
		out[i] = f(v)
	}
	return out
}
