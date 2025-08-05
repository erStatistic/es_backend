package erapi

type userResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	User    User   `json:"user"`
}

type GameResponse struct {
	Code      int        `json:"code"`
	Message   string     `json:"message"`
	UserGames []UserGame `json:"userGames"`
	Next      int        `json:"next"`
}

type TopRankResponse struct {
	Code     int        `json:"code"`
	Message  string     `json:"message"`
	TopRanks []UserRank `json:"topRanks"`
}

type RankResponse struct {
	Code     int      `json:"code"`
	Message  string   `json:"message"`
	UserRank UserRank `json:"userRank"`
}
