package erapi

type Weapons []struct {
	AttackRange float64 `json:"attackRange"`
	AttackSpeed float64 `json:"attackSpeed"`
	ShopFilter  int     `json:"shopFilter"`
	Type        string  `json:"type"`
}
