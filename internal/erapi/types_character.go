package erapi

type Characters []struct {
	AdaptiveForce                    int     `json:"adaptiveForce"`
	AttackPower                      int     `json:"attackPower"`
	AttackSpeed                      float64 `json:"attackSpeed"`
	AttackSpeedLimit                 float64 `json:"attackSpeedLimit"`
	AttackSpeedMin                   float64 `json:"attackSpeedMin"`
	AttackSpeedRatio                 int     `json:"attackSpeedRatio"`
	CharArcheType1                   string  `json:"charArcheType1"`
	CharArcheType2                   string  `json:"charArcheType2"`
	Code                             int     `json:"code"`
	ComboTime                        int     `json:"comboTime"`
	CriticalStrikeChance             int     `json:"criticalStrikeChance"`
	Defense                          int     `json:"defense"`
	ForcedComboAttack                bool    `json:"forcedComboAttack"`
	HpRegen                          float64 `json:"hpRegen"`
	IncreaseBasicAttackDamageRatio   int     `json:"increaseBasicAttackDamageRatio"`
	InitExtraPoint                   int     `json:"initExtraPoint"`
	InitStateDisplayIndex            int     `json:"initStateDisplayIndex"`
	IsUIPlayerStatusBarStaminaFilled bool    `json:"isUIPlayerStatusBarStaminaFilled"`
	LobbySubObject                   string  `json:"lobbySubObject"`
	LocalScaleInCutscene             int     `json:"localScaleInCutscene"`
	LocalScaleInVictoryScene         string  `json:"localScaleInVictoryScene"`
	MaxExtraPoint                    int     `json:"maxExtraPoint"`
	MaxHp                            int     `json:"maxHp"`
	MaxSp                            int     `json:"maxSp"`
	MoveSpeed                        float64 `json:"moveSpeed"`
	Name                             string  `json:"name"`
	PathingRadius                    float64 `json:"pathingRadius"`
	PreventBasicAttackDamagedRatio   int     `json:"preventBasicAttackDamagedRatio"`
	PreventSkillDamagedRatio         int     `json:"preventSkillDamagedRatio"`
	Radius                           float64 `json:"radius"`
	Resource                         string  `json:"resource"`
	SightRange                       float64 `json:"sightRange"`
	SkillAmp                         int     `json:"skillAmp"`
	SkillAmpRatio                    int     `json:"skillAmpRatio"`
	SpRegen                          float64 `json:"spRegen"`
	StrLearnStartSkill               string  `json:"strLearnStartSkill"`
	StrUsePointLearnStartSkill       string  `json:"strUsePointLearnStartSkill"`
	UIHeight                         float64 `json:"uiHeight"`
	WeaponRangeType                  string  `json:"weaponRangeType"`
}
