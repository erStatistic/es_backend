package main

import (
	"fmt"
)

func commandUserGame(cfg *config, args ...string) error {

	userNum := 0
	nextFlag := false

	if len(args) == 0 {
		userNum = cfg.currentUser.UserNum
	} else if len(args) > 2 {
		return fmt.Errorf("usage : userGame [nickname] [--next / --n]")
	} else {
		nickname := args[0]
		user, err := cfg.esapiClient.UserByNickname(nickname)
		if err != nil {
			return err
		}
		if user == nil {
			return fmt.Errorf("user not found")
		}

	}

	for _, arg := range args {
		if arg == "--next" || arg == "--n" {
			nextFlag = true
		}
	}

	if nextFlag {
		usergames, next, err := cfg.esapiClient.GameByUserNum(userNum, cfg.nextgame)
		if err != nil {
			return err
		}
		for _, game := range usergames {
			fmt.Printf("GameID: %d\n", game.GameID)
		}
		cfg.nextgame = next
	} else {
		usergames, next, err := cfg.esapiClient.GameByUserNum(userNum, nil)
		if err != nil {
			return err
		}
		for _, game := range usergames {
			fmt.Printf("GameID: %d\n", game.GameID)
		}
		cfg.nextgame = next
	}

	fmt.Printf("nextGameID: %d\n", *cfg.nextgame)

	return nil
}
