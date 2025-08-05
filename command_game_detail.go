package main

import (
	"errors"
	"fmt"
	"strconv"
)

func commandGameDetail(cfg *config, args ...string) error {
	if len(args) > 1 || len(args) == 0 {
		return errors.New("usage : gamedetail [gameid]")
	}

	arg1 := args[0]

	gameID, err := strconv.Atoi(arg1)
	games, err := cfg.esapiClient.GameByGameID(gameID)
	if err != nil {
		return err
	}

	for i, game := range games {
		fmt.Printf("%02d | Nickname: %s, UserNum: %d\n", i+1, game.Nickname, game.UserNum)
		fmt.Printf("   | Totaltime: %d\n", game.TotalTime)
	}
	return nil
}
