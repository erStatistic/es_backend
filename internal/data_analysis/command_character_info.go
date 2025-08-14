package data_analysis

import "fmt"

func commandCharacterInfo(cfg *Config, args ...string) error {

	characters, err := cfg.EsapiClient.GetCharacterInfo()
	if err != nil {
		return err
	}
	for i, character := range characters {
		fmt.Printf("%03d | %s\n", i+1, character.Name)
	}
	return nil

}
