package erapi

type UserStats struct {
	UserNum        int                 `json:"userNum"`
	Rank           int                 `json:"rank"`
	Mmr            int                 `json:"mmr"`
	Nickname       string              `json:"nickname"`
	TotalGames     int64               `json:"totalGames"`
	TotalWins      int64               `json:"totalWins"`
	AverageRank    float64             `json:"averageRank"`
	Top1Rage       float64             `json:"top1Rage"`
	Top2Rage       float64             `json:"top2Rage"`
	Top3Rage       float64             `json:"top3Rage"`
	CharacterStats []UserCharacterStat `json:"characterStats"`
}

type UserCharacterStat struct {
	CharacterCode int     `json:"characterCode"`
	Usage         int64   `json:"usage"`
	MaxKillings   int     `json:"maxKillings"`
	Top3          int     `json:"top3"`
	Wins          int     `json:"wins"`
	Top3Rage      float64 `json:"top3Rage"`
	AverageRank   float64 `json:"averageRank"`
}
