package data_analysis

import "fmt"

func commandHelp(cfg *Config, args ...string) error {
	fmt.Println()
	fmt.Println("Welcome ER command line interface")
	fmt.Println()
	for _, cmd := range getCommands() {
		fmt.Printf(" * %s: %s\n", cmd.name, cmd.descrpition)
	}
	fmt.Println()
	return nil
}
