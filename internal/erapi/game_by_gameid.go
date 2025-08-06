package erapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GameByGameID(gameID int) ([]UserGame, error) {

	url := fmt.Sprintf("%s/games/%d", baseURLv1, gameID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-api-key", c.apiKey)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP failed to get game by game ID: %d", res.StatusCode)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	result := GameResponse{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	for i, game := range result.UserGames {
		result.UserGames[i] = convertKeys(game)
	}

	if result.Code != 200 {
		fmt.Printf("GameByGameID StatusCode : %d\n", result.Code)
		return nil, err
	}
	return result.UserGames, nil
}
