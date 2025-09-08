package data_analysis

import "fmt"

func commandWeaponInfo(cfg *Config, args ...string) error {

	weapons, err := cfg.EsapiClient.GetWeaponInfo()
	if err != nil {
		return err
	}
	for i, weapon := range weapons {
		fmt.Printf("%03d | %s\n", i+1, weapon.Type)
	}
	return nil

}
