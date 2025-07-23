package main

import (
	"fmt"
)

func commandUserGame(cfg *config, args ...string) error {

	var userNum int
	var nextFlag bool

	for _, arg := range args {
		if arg == "--next" || arg == "--n" {
			nextFlag = true
			continue
		}
		if arg != "" {
			user, err := cfg.esapiClient.UserByNickname(arg)
			if err != nil {
				return err
			}
			userNum = user.UserNum
			break
		}
	}
	if userNum == 0 {
		userNum = cfg.currentUser.UserNum
	}

	// Validate argument count
	if len(args) > 2 {
		return fmt.Errorf("usage: userGame [nickname] [--next / --n]")
	}

	var nextgame *int
	if nextFlag {
		nextgame = cfg.nextgame
	}

	usergames, nxt, err := cfg.esapiClient.GameByUserNum(userNum, nextgame)
	if err != nil {
		return err
	}
	cfg.nextgame = nxt

	for _, game := range usergames {
		fmt.Printf("GameID: %d\n", game.GameID)
	}

	fmt.Printf("nextGameID: %d\n", *cfg.nextgame)

	return nil
}
