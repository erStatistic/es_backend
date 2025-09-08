package erapi

type UserRank struct {
	UserNum    int    `json:"userNum"`
	Nickname   string `json:"nickname"`
	Mmr        int    `json:"mmr"`
	ServerRank int    `json:"serverRank"`
	Rank       int    `json:"rank"`
}
