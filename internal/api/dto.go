package rumiapi

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
