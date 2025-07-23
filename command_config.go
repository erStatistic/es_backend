package main

import "fmt"

func commandConfig(cfg *config, args ...string) error {

	fmt.Println("=============================")
	fmt.Println("           Config            ")
	fmt.Println("=============================")

	fmt.Println()
	fmt.Println()

	fmt.Printf("CurrentUser : %s\n", cfg.currentUser.Nickname)
	fmt.Printf("Nextgame : %d\n", *cfg.nextgame)

	fmt.Println()
	fmt.Println()
	return nil
}
