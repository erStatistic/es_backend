package main

import (
	"fmt"
)

func commandUser(cfg *config, args ...string) error {
	if len(args) > 1 {
		return fmt.Errorf("usage : user [nickname]")
	}

	nickname := args[0]

	user, err := cfg.esapiClient.UserByNickname(nickname)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	exists := false
	for _, u := range cfg.users {
		if u.UserNum == user.UserNum {
			exists = true
			break
		}
	}

	if !exists {
		cfg.users = append(cfg.users, *user)
		cfg.currentUser = user
	}

	fmt.Printf("Nickname: %s\n", user.Nickname)
	fmt.Printf("Usernum: %d\n", user.UserNum)
	return nil
}
