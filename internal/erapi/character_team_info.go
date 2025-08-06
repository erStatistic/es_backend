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
	maxConcurrency = 5
	maxIterations  = 3
	maxGames       = 100000
	maxUsers       = 10000
	rangeDay       = 7
)

func (c *Client) GetCharacterTeamInfo(userID []int) ([]int, []int, error) {

	userIDList, gameIDList := c.processUserIDs(userID)
	filtereduserIDList := []int{}
	for _, userID := range userIDList {
		rank, err := c.RankUserInfo(userID, 3, 33)
		if err != nil {
			return nil, nil, err
		}
		if rank.Mmr >= 6400 {
			filtereduserIDList = append(filtereduserIDList, userID)
		}

	}

	return filtereduserIDList, gameIDList, nil
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

	// 초기 유저는 무조건 userIDList에 포함
	for _, uid := range initialUserIDs {
		processedUsers[uid] = true
		userIDList = append(userIDList, uid)
	}

	for len(queue) > 0 {
		nextQueue := make([]int, 0, maxUsers)
		var wg sync.WaitGroup

		for _, uid := range queue {
			semaphore <- struct{}{}
			wg.Add(1)

			go func(userID int) {
				defer wg.Done()
				defer func() { <-semaphore }()

				games, err := c.ManyGames(userID, 3)
				if err != nil {
					fmt.Printf("ManyGames Error (userID %d): %v\n", userID, err)
					return
				}

				recentGameIDs := filterRecentGames(games)

				for _, gameID := range recentGameIDs {
					gameData, err := c.GameByGameID(gameID)
					if err != nil {
						fmt.Printf("GameByGameID Error (gameID %d): %v\n", gameID, err)
						continue
					}

					ok := true
					for _, g := range gameData {
						if g.MmrAvg < 6000 {
							ok = false
							break
						}
					}
					if !ok {
						continue
					}

					mutex.Lock()
					if processedGames[gameID] || len(processedGames) >= maxGames {
						mutex.Unlock()
						continue
					}
					processedGames[gameID] = true
					gameIDList = append(gameIDList, gameID)
					mutex.Unlock()

					// Process each user in the game
					for _, g := range gameData {
						rank, err := c.RankUserInfo(g.UserNum, 3, 33)
						if err != nil {
							fmt.Printf("RankUserInfo Error (userNum %d): %v\n", g.UserNum, err)
							continue
						}
						if rank.Mmr < 6400 {
							continue
						}

						mutex.Lock()
						if !processedUsers[g.UserNum] && len(processedUsers) < maxUsers {
							processedUsers[g.UserNum] = true
							userIDList = append(userIDList, g.UserNum)
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

	// ✅ initialUserIDs는 반드시 포함되도록 보장
	finalUserSet := make(map[int]bool)
	for _, id := range userIDList {
		finalUserSet[id] = true
	}
	for _, id := range initialUserIDs {
		finalUserSet[id] = true
	}

	finalUserIDs := make([]int, 0, len(finalUserSet))
	for id := range finalUserSet {
		finalUserIDs = append(finalUserIDs, id)
	}

	return finalUserIDs, gameIDList
}

func (c *Client) ManyGames(userID, count int) ([]UserGame, error) {
	totalGames := []UserGame{}
	var next *int

	for range count {
		games, nextPage, err := c.GameByUserNum(userID, next)
		if err != nil {
			return nil, err
		}
		totalGames = append(totalGames, games...)
		if nextPage == nil {
			break
		}
		next = nextPage
	}

	return totalGames, nil
}

//	func (c *Client) processUserIDs(initialUserIDs []int) ([]int, []int) {
//		var (
//			userIDList     []int
//			gameIDList     []int
//			processedUsers = make(map[int]bool, maxUsers)
//			processedGames = make(map[int]bool, maxGames)
//			mutex          sync.Mutex
//			semaphore      = make(chan struct{}, maxConcurrency)
//		)
//
//		queue := make([]int, 0, maxUsers)
//		queue = append(queue, initialUserIDs...)
//
//		for len(queue) > 0 {
//			nextQueue := make([]int, 0, maxUsers)
//			var wg sync.WaitGroup
//
//			for _, uid := range queue {
//				mutex.Lock()
//				if processedUsers[uid] || len(processedUsers) >= maxUsers {
//					mutex.Unlock()
//					continue
//				}
//				processedUsers[uid] = true
//				userIDList = append(userIDList, uid)
//				mutex.Unlock()
//
//				semaphore <- struct{}{}
//				wg.Add(1)
//
//				go func(userID int) {
//					defer wg.Done()
//					defer func() { <-semaphore }()
//					games, err := c.ManyGames(userID, 3)
//					if err != nil {
//						fmt.Printf("GameByUserNum Error: %v\n", err)
//						return
//					}
//					recentGameIDs := filterRecentGames(games)
//
//					for _, gameID := range recentGameIDs {
//
//						gameData, err := c.GameByGameID(gameID)
//						if err != nil {
//							fmt.Printf("GameByGameID Error: %v\n", err)
//							continue
//						}
//
//						ok := true
//						for _, g := range gameData {
//							if g.MmrAvg < 6000 {
//								ok = false
//								break
//							}
//						}
//
//						if ok {
//							mutex.Lock()
//							if processedGames[gameID] || len(processedGames) >= maxGames {
//								mutex.Unlock()
//								continue
//							}
//							processedGames[gameID] = true
//							gameIDList = append(gameIDList, gameID)
//
//							for _, g := range gameData {
//								userID := g.UserNum
//								rank, err := c.RankUserInfo(userID, 3, 33)
//								if err != nil {
//									fmt.Printf("RankUserInfo Error: %v\n", err)
//									continue
//								}
//								if rank.Mmr < 6400 {
//									continue
//								}
//
//								if !processedUsers[g.UserNum] && len(processedUsers) < maxUsers {
//									nextQueue = append(nextQueue, g.UserNum)
//								}
//							}
//							mutex.Unlock()
//
//						}
//					}
//				}(uid)
//			}
//			wg.Wait()
//			queue = nextQueue
//		}
//		return userIDList, gameIDList
//	}
// func (c *Client) ManyGames(userID, count int) ([]UserGame, error) {
// 	totalGames := []UserGame{}
//
// 	games, next, err := c.GameByUserNum(userID, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	totalGames = append(totalGames, games...)
//
// 	count--
//
// 	for range count {
// 		games, next, err = c.GameByUserNum(userID, next)
// 		if err != nil {
// 			return nil, err
// 		}
// 		totalGames = append(totalGames, games...)
// 	}
//
// 	return totalGames, nil
// }
