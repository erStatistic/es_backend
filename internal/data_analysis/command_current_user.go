package data_analysis

import "fmt"

func commandCurrentUser(cfg *Config, args ...string) error {
	if cfg.CurrentUser == nil {
		return fmt.Errorf("no current user")
	}
	fmt.Println()
	fmt.Println("Current User")
	fmt.Println()
	fmt.Printf("Nickname: %s\n", cfg.CurrentUser.Nickname)
	fmt.Printf("Usernum: %d\n", cfg.CurrentUser.UserNum)
	fmt.Println()
	return nil
}
