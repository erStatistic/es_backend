package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/kaeba0616/es_backend/internal/erapi"
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
	allofteams := map[int][]TeamInfo{}

	for _, gameID := range gameIDs {
		teams := getTeamInfo(cfg, gameID)
		for gameRank := range teams {
			allofteams[gameRank] = append(allofteams[gameRank], teams[gameRank])
		}
	}
	if err := saveTeamInfoToCSV("./internal/output/data/team_info.csv", allofteams); err != nil {
		return fmt.Errorf("failed to save team info to CSV: %w", err)
	}

	if err := saveIDsToFile("./internal/output/data/user_ids.txt", userIDs); err != nil {
		return fmt.Errorf("failed to save user IDs: %w", err)
	}

	if err := saveIDsToFile("./internal/output/data/game_ids.txt", gameIDs); err != nil {
		return fmt.Errorf("failed to save game IDs: %w", err)
	}

	fmt.Printf("Total users: %d\n", len(userIDs))
	fmt.Printf("Total games: %d\n", len(gameIDs))
	return nil
}

type TeamInfo struct {
	GameID         int
	CharacterNums  []int
	TeamKills      int
	MonsterCredits int
}

func getTeamInfo(cfg *config, gameID int) map[int]TeamInfo {
	fmt.Printf("Getting team info for game ID: %d\n", gameID)
	usergames, err := cfg.esapiClient.GameByGameID(gameID)
	if err != nil {
		fmt.Printf("Failed to get game by game ID: %v\n", err)
		return nil
	}
	teams := map[int]TeamInfo{}
	for _, game := range usergames {
		rank := game.GameRank
		charNum := game.CharacterNum
		kills := game.TeamKill
		credits := getMonsterCredits(game)

		team := teams[rank]

		team.GameID = game.GameID
		team.CharacterNums = append(team.CharacterNums, charNum)
		team.TeamKills = kills
		team.MonsterCredits += credits
		teams[rank] = team
	}

	for _, team := range teams {
		sort.Ints(team.CharacterNums)
	}

	return teams
}

func getMonsterCredits(game erapi.UserGame) int {
	creditChicken := game.KillBatGainVFCredit
	creditBoar := game.KillBoarGainVFCredit
	creditWolf := game.KillWolfGainVFCredit
	creditBear := game.KillBearGainVFCredit
	creditBat := game.KillBatGainVFCredit
	creditWildDog := game.KillWildDogGainVFCredit
	creditRaven := game.CreditSource.KillRaven + game.CreditSource.KillMutantRaven
	totalcredit := creditChicken + creditBoar + creditWolf + creditBear + creditBat + creditWildDog + int(creditRaven)
	return totalcredit
}

func saveTeamInfoToCSV(filename string, allofteams map[int][]TeamInfo) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"GameRank", "GameID", "CharacterNums", "TeamKills", "MonsterCredits"}
	if err := writer.Write(header); err != nil {
		return err
	}

	ranks := make([]int, 0, len(allofteams))
	for rank := range allofteams {
		ranks = append(ranks, rank)
	}
	sort.Ints(ranks)

	for _, rank := range ranks {
		teams := allofteams[rank]
		for _, team := range teams {
			charNumsStr := fmt.Sprintf("%v", team.CharacterNums)
			record := []string{
				fmt.Sprintf("%d", rank),
				fmt.Sprintf("%d", team.GameID),
				charNumsStr,
				fmt.Sprintf("%d", team.TeamKills),
				fmt.Sprintf("%d", team.MonsterCredits),
			}
			if err := writer.Write(record); err != nil {
				return err
			}
		}
	}
	return nil
}

func saveIDsToFile(filename string, ids []int) error {
	var sb strings.Builder
	for _, id := range ids {
		sb.WriteString(fmt.Sprintf("%d\n", id))
	}
	return os.WriteFile(filename, []byte(sb.String()), 0644)
}
