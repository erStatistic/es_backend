package main

import (
	"fmt"
	"os"
	"strings"
)

// const (
// 	maxConcurrency = 40
// 	maxIterations  = 3
// 	maxGames       = 10000
// )

func commandCharacterTeamInfo(cfg *config, args ...string) error {

	if len(args) > 0 {
		return fmt.Errorf("usage: characterteaminfo")
	}
	initialUserIDs := []int{}
	for _, user := range cfg.rankers {
		initialUserIDs = append(initialUserIDs, user.UserNum)
	}
	fmt.Println()
	fmt.Println("Start character team info")
	userIDs, gameIDs, err := cfg.esapiClient.GetCharacterTeamInfo(initialUserIDs)
	if err != nil {
		return err
	}
	fmt.Println("End character team info")

	if err := saveIDsToFile("user_ids.txt", userIDs); err != nil {
		return fmt.Errorf("failed to save user IDs: %w", err)
	}

	if err := saveIDsToFile("game_ids.txt", gameIDs); err != nil {
		return fmt.Errorf("failed to save game IDs: %w", err)
	}

	fmt.Printf("Total users: %d\n", len(userIDs))
	fmt.Printf("Total games: %d\n", len(gameIDs))
	return nil
}

func saveIDsToFile(filename string, ids []int) error {
	var sb strings.Builder
	for _, id := range ids {
		sb.WriteString(fmt.Sprintf("%d\n", id))
	}
	return os.WriteFile(filename, []byte(sb.String()), 0644)
}
