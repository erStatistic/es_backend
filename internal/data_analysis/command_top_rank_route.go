package data_analysis

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kaeba0616/es_backend/internal/erapi"
)

type CharacterRoute struct {
	CharacterCode int
	WeaponCode    int
	Title         string
	Count         int
}

func commandTopRankRoute(cfg *Config, args ...string) error {
	if len(args) > 0 {
		return nil
	}

	characaterTops, err := characterTops(cfg)
	if err != nil {
		return err
	}

	Routes := make(map[int]CharacterRoute)

	for characterCode, userNums := range characaterTops {

		for _, userNum := range userNums {
			recentGames := filterRecentGames(cfg, userNum, characterCode)
			for _, game := range recentGames {
				routeId := game.RouteIDOfStart
				route, ok := Routes[routeId]
				if !ok {
					gameRoute, err := cfg.EsapiClient.GameRoute(routeId)
					if err != nil {
						fmt.Printf("GameRoute Error (routeID %d): %v\n", routeId, err)
						continue
					}

					weaponNum := game.BestWeapon
					if characterCode == 44 {
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

					route = CharacterRoute{
						CharacterCode: characterCode,
						WeaponCode:    weaponNum,
						Count:         1,
						Title:         gameRoute.RecommendWeaponRoute.Title,
					}
					Routes[routeId] = route
				} else {
					route.Count++
					Routes[routeId] = route
				}
			}

		}
	}
	if err := saveRoutesToCSV("./internal/output/games/route.csv", Routes); err != nil {
		return fmt.Errorf("failed to save route to CSV: %w", err)
	}

	return nil
}

func saveRoutesToCSV(filename string, routes map[int]CharacterRoute) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	header := []string{"route_id", "title", "character_code", "weapon_code", "count"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for i, route := range routes {
		record := []string{
			fmt.Sprintf("%d", i),
			fmt.Sprint(route.Title),
			fmt.Sprintf("%d", route.CharacterCode),
			fmt.Sprintf("%d", route.WeaponCode),
			fmt.Sprintf("%d", route.Count),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func characterTops(cfg *Config) (map[int][]int, error) {

	characaterTops := make(map[int][]int)
	ranks := cfg.Rankers
	for _, rank := range ranks {
		userStats, err := cfg.EsapiClient.UserStatByUserId(rank.UserNum)
		if err != nil {
			return nil, err
		}
		for _, stat := range userStats {
			characaterTops[stat.CharacterCode] = append(characaterTops[stat.CharacterCode], rank.UserNum)
		}

	}

	for characterCode, userNums := range characaterTops {
		fmt.Printf("CharacterCode: %d\n", characterCode)
		for _, userNum := range userNums {
			fmt.Printf("UserNum: %d\n", userNum)
		}
	}

	return characaterTops, nil
}

func filterRecentGames(cfg *Config, userNum int, characterCode int) []erapi.UserGame {
	var (
		page  *int
		games []erapi.UserGame
	)

	const (
		rangeDay = 14
	)

	today := time.Now().In(time.FixedZone("KST", 9*60*60))
	cutoff := today.AddDate(0, 0, -rangeDay)
	count := 0
	isNext := true
	for range 4 {
		usergames, next, err := cfg.EsapiClient.GameByUserNum(userNum, page)
		if err != nil {
			fmt.Printf("GameByUserNum Error (userNum %d): %v\n", userNum, err)

		}
		if next != nil {
			fmt.Println(*next)
		}
		if !isNext {
			break
		}

		for _, game := range usergames {
			fmt.Printf("Count: %d\n", count)
			if count > 10 {
				isNext = false
				break
			}

			startDtm, err := time.Parse("2006-01-02T15:04:05.000-0700", game.StartDtm)
			if err != nil {
				fmt.Printf("Failed to parse date %s: %v\n", game.StartDtm, err)
			}
			if game.CharacterNum != characterCode || game.SeasonID != 33 {
				continue
			}

			if startDtm.After(cutoff) {
				fmt.Printf("GameID: %d\n", game.GameID)
				games = append(games, game)
				count++
			}

		}
		page = next

	}

	return games
}

func routeData(cfg *Config, routeID int) (erapi.Route, error) {

	route, err := cfg.EsapiClient.GameRoute(routeID)
	if err != nil {
		fmt.Printf("GameRoute Error (routeID %d): %v\n", routeID, err)
		return erapi.Route{}, err
	}
	return route, nil
}
