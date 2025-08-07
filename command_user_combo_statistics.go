package main

import (
	"fmt"
)

func commandUserComboStatistics(cfg *config, args ...string) error {
	users := args
	count := len(users)

	userNums := []int{}
	for _, user := range users {
		userNum, err := cfg.esapiClient.UserByNickname(user)
		if err != nil {
			return err
		}
		userNums = append(userNums, userNum.UserNum)
	}

	chars := [][]int{
		{81, 57, 16},
	}
	combos, err := cfg.esapiClient.AnalysisResult(chars)
	if err != nil {
		return err
	}

	fmt.Println()
	for i, combo := range combos {
		fmt.Printf("Combo %d\n", i+1)
		fmt.Printf("Combo: %v\n", combo.Clusters)
		fmt.Printf("Top3Ratio: %v\n", combo.Top3Ratio)
		fmt.Printf("Scores: %v\n", combo.Scores)
		fmt.Println()
	}

	switch count {

	// count = 0
	// show top 5 combo in every game
	case 0:
		fmt.Println("Top 5 combo in every game Count 0")
		fmt.Println()
		combos := cfg.esapiClient.SortClusterDistByNormalizedScore()
		// until 5 index iteam
		for i := range 5 {
			fmt.Printf("Combo %d\n", i+1)
			fmt.Printf("Combo: %v\n", combos[i].ClusterComboKey)
			fmt.Printf("Total: %v\n", combos[i].Total)
			fmt.Printf("Top3Ratio: %v\n", combos[i].Top3Ratio)
			fmt.Printf("Scores: %v\n", combos[i].NormalizedScore)
			fmt.Println()
		}

	// count = 1
	// Show the top 5 combinations of characters, each made up of the top 3 characters from this user.
	case 1:
		err := getCombos(cfg, userNums)
		if err != nil {
			return err
		}

	// count = 2
	// Show the top 5 combinations of characters, each made up of the top 3 characters from those two users.
	case 2:
		err := getCombos(cfg, userNums)
		if err != nil {
			return err
		}

	// count = 3
	// Show the top 5 combinations of characters, each made up of the top 3 characters from those three users.
	case 3:
		err := getCombos(cfg, userNums)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("usage : statistics [...user] until three users")
	}

	return nil
}

func getCombos(cfg *config, userNums []int) error {
	_, err := cfg.esapiClient.UserStatByUserIdList(userNums)
	if err != nil {
		return err
	}

	return nil
}
