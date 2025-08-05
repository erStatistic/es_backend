package erapi

import (
	"fmt"
	"sync"
	"time"
)

type GameInfo struct {
	UserNums []int
	Teams    [][]int
}

const (
	maxConcurrency = 20
	maxIterations  = 3
	maxGames       = 100000
	maxUsers       = 10000
	rangeDay       = 7
)

func (c *Client) GetCharacterTeamInfo(userID []int) ([]int, []int, error) {

	userIDList, gameIDList := c.processUserIDs(userID)

	return userIDList, gameIDList, nil
}

func filterRecentGames(games []UserGame) []int {

	today := time.Now().In(time.FixedZone("KST", 9*60*60))
	cutoff := today.AddDate(0, 0, -rangeDay)

	gameIDs := []int{}
	for _, game := range games {
		fmt.Printf("GameID: %d\n", game.GameID)

		startDtm, err := time.Parse("2006-01-02T15:04:05.000-0700", game.StartDtm)
		if err != nil {
			fmt.Printf("Failed to parse date %s: %v\n", game.StartDtm, err)
		}
		if startDtm.After(cutoff) && game.SeasonID == 33 {
			gameIDs = append(gameIDs, game.GameID)
		}
	}
	return gameIDs
}

func (c *Client) processUserIDs(initialUserIDs []int) ([]int, []int) {
	var (
		userIDList     []int
		gameIDList     []int
		processedUsers = make(map[int]bool, maxUsers)
		processedGames = make(map[int]bool, maxGames)
		mutex          sync.Mutex
		semaphore      = make(chan struct{}, maxConcurrency)
	)

	queue := make([]int, 0, maxUsers)
	queue = append(queue, initialUserIDs...)

	for len(queue) > 0 {
		nextQueue := make([]int, 0, maxUsers)
		var wg sync.WaitGroup

		for _, uid := range queue {
			mutex.Lock()
			if processedUsers[uid] || len(processedUsers) >= maxUsers {
				mutex.Unlock()
				continue
			}
			processedUsers[uid] = true
			userIDList = append(userIDList, uid)
			mutex.Unlock()

			semaphore <- struct{}{}
			wg.Add(1)

			go func(userID int) {
				defer wg.Done()
				defer func() { <-semaphore }()
				games, err := c.ManyGames(userID, 3)
				if err != nil {
					fmt.Printf("GameByUserNum Error: %v\n", err)
					return
				}
				recentGameIDs := filterRecentGames(games)

				for _, gameID := range recentGameIDs {
					gameData, err := c.GameByGameID(gameID)
					if err != nil {
						fmt.Printf("GameByGameID Error: %v\n", err)
						continue
					}
					for _, g := range gameData {
						mutex.Lock()
						if processedGames[g.GameID] || len(processedGames) >= maxGames {
							mutex.Unlock()
							continue
						}
						processedGames[g.GameID] = true
						gameIDList = append(gameIDList, g.GameID)

						if !processedUsers[g.UserNum] && len(processedUsers) < maxUsers {
							nextQueue = append(nextQueue, g.UserNum)
						}
						mutex.Unlock()
					}
				}
			}(uid)
		}
		wg.Wait()
		queue = nextQueue
	}
	return userIDList, gameIDList
}

func (c *Client) ManyGames(userID, count int) ([]UserGame, error) {
	totalGames := []UserGame{}

	games, next, err := c.GameByUserNum(userID, nil)
	if err != nil {
		return nil, err
	}
	totalGames = append(totalGames, games...)

	count--

	for range count {
		games, next, err = c.GameByUserNum(userID, next)
		if err != nil {
			return nil, err
		}
		totalGames = append(totalGames, games...)
	}

	return totalGames, nil
}
