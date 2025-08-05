package main

import (
	"fmt"
	"strconv"
)

func commandUserRankInfo(cfg *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: userrankinfo <userID>")
	}

	arg1 := args[0]
	userID, err := strconv.Atoi(arg1)
	if err != nil {
		return err
	}

	userRank, err := cfg.esapiClient.RankUserInfo(userID, 3, 33)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("UserRank")
	fmt.Printf("UserNum: %d\n", userRank.UserNum)
	fmt.Printf("Nickname: %s\n", userRank.Nickname)
	fmt.Printf("Mmr: %d\n", userRank.Mmr)
	fmt.Printf("ServerRank: %d\n", userRank.ServerRank)
	fmt.Printf("Rank: %d\n", userRank.Rank)
	fmt.Println()

	return nil
}
