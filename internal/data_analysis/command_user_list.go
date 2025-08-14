package data_analysis

import "fmt"

func commandUserList(cfg *Config, args ...string) error {
	renderRepl(cfg)
	fmt.Printf("CurrentUser: %s\n", cfg.CurrentUser.Nickname)
	fmt.Println()
	return nil
}
