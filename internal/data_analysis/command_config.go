package data_analysis

import "fmt"

func commandConfig(cfg *Config, args ...string) error {

	fmt.Println("=============================")
	fmt.Println("           Config            ")
	fmt.Println("=============================")

	fmt.Println()
	fmt.Println()

	fmt.Printf("CurrentUser : %s\n", cfg.CurrentUser.Nickname)
	fmt.Printf("Nextgame : %d\n", *cfg.Nextgame)

	fmt.Println()
	fmt.Println()
	return nil
}
