package rumiapi

import "time"

type CharacterCreateRequest struct {
	Code     int32  `json:"code"`
	NameKr   string `json:"nameKr"`
	ImageUrl string `json:"imageUrl"`
}
type CharacterPatchRequest struct {
	NameKr   *string `json:"nameKr,omitempty"`
	ImageUrl *string `json:"imageUrl,omitempty"`
}

type APIResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

type GameDTO struct {
	ID         int32      `json:"id"`
	GameCode   int64      `json:"game_code"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	AverageMmr int32      `json:"average_mmr"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type GameTeamDTO struct {
	ID             int32     `json:"id"`
	GameID         int64     `json:"gameId"`
	TeamID         int32     `json:"teamId"`
	GameRank       int32     `json:"gameRank"`
	TeamKills      int32     `json:"teamKills"`
	MonsterCredits int32     `json:"monsterCredits"`
	GainedMmr      int32     `json:"gainedMmr"`
	TeamAvgMmr     int32     `json:"teamAvgMmr"`
	TotalTime      int32     `json:"totalTime"`
	TimesID        *int32    `json:"timesId,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type TierDTO struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	Rank      int32     `json:"rank"`
	ImageUrl  string    `json:"imageUrl"`
	MmrMin    int32     `json:"mmr_min"`
	MmrMax    int32     `json:"mmr_max"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TimeDTO struct {
	ID        int32     `json:"id"`
	No        int32     `json:"no"`
	Name      string    `json:"name"`
	Seconds   int32     `json:"seconds"`
	StartTime int32     `json:"startTime"` // time_range 하한 (포함, [start,end))
	EndTime   int32     `json:"endTime"`   // time_range 상한 (제외, [start,end))
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
