package main

import "fmt"

func commandCharacterInfo(cfg *config, args ...string) error {

	characters, err := cfg.esapiClient.GetCharacterInfo()
	if err != nil {
		return err
	}
	for i, character := range characters {
		fmt.Printf("%03d | %s\n", i+1, character.Name)
	}
	return nil

}
