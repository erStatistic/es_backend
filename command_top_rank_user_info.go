package main

import "fmt"

func commandTopRankUserInfo(cfg *config, args ...string) error {

	teamMode := 3
	seasonId := 33
	serverCode := 10

	ranks, err := cfg.esapiClient.TopRankUserInfo(teamMode, seasonId, serverCode)
	if err != nil {
		return err
	}

	cfg.rankers = ranks

	for i, rank := range ranks {
		if i >= 900 {
			break
		}
		fmt.Printf("%03d |%d %s\n", i+1, rank.UserNum, rank.Nickname)
	}
	return nil
}
