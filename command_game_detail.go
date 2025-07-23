package main

import (
	"errors"
	"fmt"
)

func commandGameDetail(cfg *config, args ...string) error {
	if len(args) > 1 || len(args) == 0 {
		return errors.New("usage : gamedetail [gameid]")
	}

	gameID := args[0]
	games, err := cfg.esapiClient.GameByGameID(gameID)
	if err != nil {
		return err
	}

	for i, game := range games {
		fmt.Printf("%d | Nickname: %s, UserNum: %d\n", i+1, game.Nickname, game.UserNum)
	}
	return nil
}
