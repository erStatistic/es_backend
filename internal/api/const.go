package rumiapi

type ctxKey string

const (
	weaponKey              ctxKey = "weapon"
	positionKey            ctxKey = "position"
	tierKey                ctxKey = "tier"
	clusterKey             ctxKey = "cluster"
	timeKey                ctxKey = "time"
	characterWeaponKey     ctxKey = "characterWeapon"
	CharacterWeaponStatKey ctxKey = "characterWeaponStat"
	GameKey                ctxKey = "game"
	GameTeamKey            ctxKey = "gameteam"
	GameTeamCWKey          ctxKey = "gameteamcw"
	characterKey           ctxKey = "character"
	pageKey                ctxKey = "page"
	limitKey               ctxKey = "limit"
)
