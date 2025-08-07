package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/kaeba0616/es_backend/internal/erapi"
)

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

	var (
		allofteams = map[int][]TeamInfo{}
		mu         sync.Mutex
		wg         sync.WaitGroup
		sem        = make(chan struct{}, 5)
	)

	for _, gameID := range gameIDs {
		wg.Add(1)
		sem <- struct{}{}
		go func(gid int) {
			defer wg.Done()
			defer func() { <-sem }()
			teams := getTeamInfo(cfg, gid)
			mu.Lock()
			for gameRank := range teams {
				allofteams[gameRank] = append(allofteams[gameRank], teams[gameRank])
			}
			mu.Unlock()

		}(gameID)
	}
	wg.Wait()

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
	weaponNums     []int
	TeamKills      int
	MonsterCredits int
	TotalTime      int
	mmrGainInGame  int
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

		weaponName := getItemInfo(cfg, game.Equipment["0"])

		// for Echion : charNum = 44
		if charNum == 44 && weaponName != nil {
			weaponMap := map[string]int{
				"데스애더":   25,
				"블랙맘바":   26,
				"사이드와인더": 27,
			}
			for key, val := range weaponMap {
				if strings.Contains(*weaponName, key) {
					team.weaponNums = append(team.weaponNums, val)
				}
			}

		} else {
			team.weaponNums = append(team.weaponNums, game.BestWeapon)
		}

		team.TeamKills = kills
		team.MonsterCredits += credits
		team.TotalTime = game.TotalTime
		team.mmrGainInGame = game.MmrGainInGame
		teams[rank] = team
	}

	return teams
}

func getItemInfo(cfg *config, WeaponItemID int) *string {
	weapon, err := cfg.esapiClient.ItemInfo(WeaponItemID)
	if err != nil {
		fmt.Printf("Failed to get item info: %v\n", err)
		return nil
	}
	return &weapon.Name

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

	header := []string{"GameRank", "GameID", "CharacterNums", "WeaponNums", "TeamKills", "MonsterCredits", "TotalTime", "mmrGainInGame"}
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
			weaponNumsStr := fmt.Sprintf("%v", team.weaponNums)
			record := []string{
				fmt.Sprintf("%d", rank),
				fmt.Sprintf("%d", team.GameID),
				charNumsStr,
				weaponNumsStr,
				fmt.Sprintf("%d", team.TeamKills),
				fmt.Sprintf("%d", team.MonsterCredits),
				fmt.Sprintf("%d", team.TotalTime),
				fmt.Sprintf("%d", team.mmrGainInGame),
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
