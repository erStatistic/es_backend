package main

import "fmt"

func commandWeaponInfo(cfg *config, args ...string) error {

	weapons, err := cfg.esapiClient.GetWeaponInfo()
	if err != nil {
		return err
	}
	for i, weapon := range weapons {
		fmt.Printf("%03d | %s\n", i+1, weapon.Type)
	}
	return nil

}
