package erapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) RankUserInfo(teamMode, seasonId, serverCode int) ([]User, error) {
	url := fmt.Sprintf("%s/rank/top/%d/%d/%d", baseURLv1, seasonId, teamMode, serverCode)
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

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var RankResponse RankResponse
	err = json.Unmarshal(body, &RankResponse)
	if err != nil {
		return nil, err
	}

	if RankResponse.Code != 200 {
		fmt.Printf("StatusCode : %d\n", RankResponse.Code)
		return nil, err
	}

	users := []User{}
	for _, ranker := range RankResponse.TopRanks {
		user := User{
			UserNum:  ranker.UserNum,
			Nickname: ranker.Nickname,
		}
		users = append(users, user)
	}
	return users, nil
}
