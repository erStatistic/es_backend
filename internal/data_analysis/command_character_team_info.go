package data_analysis

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/kaeba0616/es_backend/internal/erapi"
)

func commandCharacterTeamInfo(cfg *Config, args ...string) error {

	if len(args) > 0 {
		return fmt.Errorf("usage: characterteaminfo")
	}

	initialUserIDs := []int{}
	for _, user := range cfg.Rankers {
		initialUserIDs = append(initialUserIDs, user.UserNum)
	}
	fmt.Println()
	fmt.Println("Start character team info")
	userIDs, gameIDs, err := cfg.EsapiClient.GetCharacterTeamInfo(initialUserIDs)
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
			for rank, team := range teams {
				if len(team.Users) == 3 {
					allofteams[rank] = append(allofteams[rank], team)
				}
			}
			mu.Unlock()

		}(gameID)
	}
	wg.Wait()

	if err := saveTeamInfoToCSV("./internal/output/games/team_info.csv", allofteams); err != nil {
		return fmt.Errorf("failed to save team info to CSV: %w", err)
	}

	if err := saveIDsToFile("./internal/output/games/user_ids.txt", userIDs); err != nil {
		return fmt.Errorf("failed to save user IDs: %w", err)
	}

	if err := saveIDsToFile("./internal/output/games/game_ids.txt", gameIDs); err != nil {
		return fmt.Errorf("failed to save game IDs: %w", err)
	}

	fmt.Printf("Total users: %d\n", len(userIDs))
	fmt.Printf("Total games: %d\n", len(gameIDs))
	return nil
}

type TeamInfo struct {
	GameID         int
	Users          []User
	TeamKills      int
	MonsterCredits int
	TotalTime      int
	mmrGainInGame  int
	AverageMmr     int
}

type User struct {
	TeamNumber   int
	CharacterNum int
	WeaponNum    int
	CurrentMmr   int
}

func getTeamInfo(cfg *Config, gameID int) map[int]TeamInfo {
	fmt.Printf("Getting team info for game ID: %d\n", gameID)
	usergames, err := cfg.EsapiClient.GameByGameID(gameID)
	if err != nil {
		fmt.Printf("Failed to get game by game ID: %v\n", err)
		return nil
	}
	teams := map[int]TeamInfo{}
	for _, game := range usergames {
		rank := game.GameRank
		charNum := game.CharacterNum
		kills := game.TeamKill
		weaponNum := game.BestWeapon
		gameID := game.GameID
		credits := getMonsterCredits(game)
		averageMMr := game.MmrAvg

		team := teams[rank]

		// for Echion : charNum = 44
		if charNum == 44 {
			weaponName := getItemInfo(cfg, game.Equipment["weapon"])
			if weaponName != nil {
				weaponMap := map[string]int{
					"데스애더":   25,
					"블랙맘바":   26,
					"사이드와인더": 27,
				}
				for key, val := range weaponMap {
					if strings.Contains(*weaponName, key) {
						weaponNum = val
						break
					}
				}
			}

		}
		user := User{}
		user.TeamNumber = game.TeamNumber
		user.CharacterNum = charNum
		user.WeaponNum = weaponNum
		user.CurrentMmr = game.MmrBefore
		team.Users = append(team.Users, user)

		team.GameID = gameID
		team.TeamKills = kills
		team.MonsterCredits += credits
		team.TotalTime = game.TotalTime
		team.mmrGainInGame = game.MmrGainInGame
		team.AverageMmr = averageMMr
		teams[rank] = team
	}
	for rank, team := range teams {

		if len(team.Users) > 3 {
			count := make(map[int]int)
			for _, user := range team.Users {
				count[user.TeamNumber]++
			}

			var filteredUsers []User
			for _, user := range team.Users {
				if count[user.TeamNumber] == 3 {
					filteredUsers = append(filteredUsers, user)
				}
			}

			team.Users = filteredUsers
			teams[rank] = team
		}
	}
	return teams
}

func getItemInfo(cfg *Config, WeaponItemID int) *string {
	weapon, err := cfg.EsapiClient.ItemInfo(WeaponItemID)
	if err != nil {
		fmt.Printf("IteamID: %d\n", WeaponItemID)
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
func ToPGIntArray(nums []int) string {
	if len(nums) == 0 {
		return "{}"
	}
	var b strings.Builder
	b.WriteByte('{')
	for i, v := range nums {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(strconv.Itoa(v))
	}
	b.WriteByte('}')
	return b.String()
}

func saveTeamInfoToCSV(filename string, allofteams map[int][]TeamInfo) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"game_rank", "game_code", "game_avg_mmr", "team_num", "characater_nums", "weapon_nums", "character_mmrs", "team_kills", "monster_credits", "total_time", "team_avg_mmr", "mmr_gain_in_game"}
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
			var characterNums, weaponNums []int
			for _, user := range team.Users {
				characterNums = append(characterNums, user.CharacterNum)
				weaponNums = append(weaponNums, user.WeaponNum)
			}
			charNumsStr := ToPGIntArray(characterNums)
			weaponNumsStr := ToPGIntArray(weaponNums)
			var characterMmrs []int
			teamAvgMmr := 0
			for _, user := range team.Users {
				characterMmrs = append(characterMmrs, user.CurrentMmr)
				teamAvgMmr += user.CurrentMmr
			}
			characterMmrsStr := ToPGIntArray(characterMmrs)
			teamAvgMmr /= len(team.Users)
			record := []string{
				fmt.Sprintf("%d", rank),
				fmt.Sprintf("%d", team.GameID),
				fmt.Sprintf("%d", team.AverageMmr),
				fmt.Sprintf("%d", team.Users[0].TeamNumber),
				charNumsStr,
				weaponNumsStr,
				characterMmrsStr,
				fmt.Sprintf("%d", team.TeamKills),
				fmt.Sprintf("%d", team.MonsterCredits),
				fmt.Sprintf("%d", team.TotalTime),
				fmt.Sprintf("%d", teamAvgMmr),
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
