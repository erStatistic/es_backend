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
	GameID         int32     `json:"gameId"`
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
