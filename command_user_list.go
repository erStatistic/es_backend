package main

import "fmt"

func commandUserList(cfg *config, args ...string) error {
	renderRepl(cfg)
	fmt.Printf("CurrentUser: %s\n", cfg.currentUser.Nickname)
	fmt.Println()
	return nil
}
