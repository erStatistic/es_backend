package erapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func (c *Client) GameByUserNum(usernum int, page *int) ([]UserGame, *int, error) {
	url := fmt.Sprintf("%s/user/games/%d", baseURLv1, usernum)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-api-key", c.apiKey)

	if page != nil {
		q := req.URL.Query()
		q.Add("next", strconv.Itoa(*page))
		req.URL.RawQuery = q.Encode()
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	if res.StatusCode != 200 {
		return nil, nil, fmt.Errorf("HTTP failed to get Game By UserNum: %d", res.StatusCode)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	result := GameResponse{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, nil, err
	}

	for i, game := range result.UserGames {
		result.UserGames[i] = convertKeys(game)
	}

	if result.Code != 200 {
		fmt.Printf("GameByUserNum StatusCode : %d\n", result.Code)
		return nil, nil, err
	}

	return result.UserGames, &result.Next, nil
}
