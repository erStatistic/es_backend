package erapi

type RankResponse struct {
	Code     int      `json:"code"`
	Message  string   `json:"message"`
	UserRank UserRank `json:"userRank"`
}

type UserRank struct {
	UserNum    int `json:"userNum"`
	ServerCode int `json:"serverCode"`
	Mmr        int `json:"mmr"`
	ServerRank int `json:"serverRank"`
	Rank       int `json:"rank"`
}
