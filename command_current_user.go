package main

import "fmt"

func commandCurrentUser(cfg *config, args ...string) error {
	if cfg.currentUser == nil {
		return fmt.Errorf("no current user")
	}
	fmt.Println()
	fmt.Println("Current User")
	fmt.Println()
	fmt.Printf("Nickname: %s\n", cfg.currentUser.Nickname)
	fmt.Printf("Usernum: %d\n", cfg.currentUser.UserNum)
	fmt.Println()
	return nil
}
