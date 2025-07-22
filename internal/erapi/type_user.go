package erapi

type userResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	User    User   `json:"user"`
}

type User struct {
	Nickname string `json:"nickname"`
	UserNum  int    `json:"userNum"`
}
