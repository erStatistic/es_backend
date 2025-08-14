package data_analysis

import (
	"fmt"
	"os"
)

func commandExit(cfg *Config, args ...string) error {
	fmt.Println("Bye!")
	os.Exit(0)
	return nil
}
