package erapi

type Route struct {
	RecommendWeaponRoute struct {
		ID                      int     `json:"id"`
		Title                   string  `json:"title"`
		UserNum                 int     `json:"userNum"`
		UserNickname            string  `json:"userNickname"`
		CharacterCode           int     `json:"characterCode"`
		SlotID                  int     `json:"slotId"`
		WeaponType              int     `json:"weaponType"`
		WeaponCodes             string  `json:"weaponCodes"`
		TraitCodes              string  `json:"traitCodes"`
		LateGameItemCodes       string  `json:"lateGameItemCodes"`
		RemoteTransferItemCodes string  `json:"remoteTransferItemCodes"`
		TacticalSkillGroupCode  int     `json:"tacticalSkillGroupCode"`
		Paths                   string  `json:"paths"`
		Count                   int     `json:"count"`
		Version                 string  `json:"version"`
		TeamMode                int     `json:"teamMode"`
		LanguageCode            string  `json:"languageCode"`
		RouteVersion            int     `json:"routeVersion"`
		Share                   bool    `json:"share"`
		UpdateDtm               int64   `json:"updateDtm"`
		V2Like                  int     `json:"v2Like"`
		V2WinRate               float64 `json:"v2WinRate"`
		V2SeasonID              int     `json:"v2SeasonId"`
		V2AccumulateLike        int     `json:"v2AccumulateLike"`
		V2AccumulateWinRate     float64 `json:"v2AccumulateWinRate"`
		V2AccumulateSeasonID    int     `json:"v2AccumulateSeasonId"`
	} `json:"recommendWeaponRoute"`
	RecommendWeaponRouteDesc struct {
		RecommendWeaponRouteID int    `json:"recommendWeaponRouteId"`
		SkillPath              string `json:"skillPath"`
		Desc                   string `json:"desc"`
	} `json:"recommendWeaponRouteDesc"`
}
